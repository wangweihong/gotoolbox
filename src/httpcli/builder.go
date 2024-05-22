package httpcli

import (
	"encoding/base64"
	"reflect"
	"strings"
	"time"

	"github.com/wangweihong/gotoolbox/src/httpcli/def"
)

type HttpRequestBuilder struct {
	httpRequest *HttpRequest
}

func NewHttpRequestBuilder() *HttpRequestBuilder {
	httpRequest := &HttpRequest{
		queryParams:          make(map[string]interface{}),
		headerParams:         make(map[string][]string),
		pathParams:           make(map[string]string),
		autoFilledPathParams: make(map[string]string),
		formParams:           make(map[string]def.FormData),
	}
	httpRequestBuilder := &HttpRequestBuilder{
		httpRequest: httpRequest,
	}
	return httpRequestBuilder
}

func (builder *HttpRequestBuilder) WithEndpoint(endpoint string) *HttpRequestBuilder {
	// 如果endpoint不带有协议,在转换成http request时, endpoint会被当作path存在url.path,而不是url.Host.
	// fillPath会替换掉endpoint.
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "http://" + endpoint
	}
	builder.httpRequest.endpoint = endpoint
	return builder
}

func (builder *HttpRequestBuilder) WithPath(path string) *HttpRequestBuilder {
	builder.httpRequest.path = path
	return builder
}

func (builder *HttpRequestBuilder) POST() *HttpRequestBuilder {
	return builder.WithMethod("POST")
}

func (builder *HttpRequestBuilder) GET() *HttpRequestBuilder {
	return builder.WithMethod("GET")
}

func (builder *HttpRequestBuilder) PUT() *HttpRequestBuilder {
	return builder.WithMethod("PUT")
}

func (builder *HttpRequestBuilder) DELETE() *HttpRequestBuilder {
	return builder.WithMethod("DELETE")
}

func (builder *HttpRequestBuilder) WithMethod(method string) *HttpRequestBuilder {
	builder.httpRequest.method = method
	return builder
}

func (builder *HttpRequestBuilder) AddQueryParam(key string, value interface{}) *HttpRequestBuilder {
	builder.httpRequest.queryParams[key] = value
	return builder
}

func (builder *HttpRequestBuilder) AddQueryParamByObjectMarshal(input interface{}) *HttpRequestBuilder {
	if input == nil {
		return builder
	}

	v := reflect.ValueOf(input)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()
	if t.Kind() != reflect.Struct {
		return builder
	}
	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)

		if !fieldType.IsExported() {
			continue
		}

		// 忽略空指针
		if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
			continue
		}
		key := fieldType.Name
		jsonTag := fieldType.Tag.Get("json")
		if jsonTag != "" {
			if fieldValue.IsZero() && strings.Contains(jsonTag, ",omitempty") {
				continue
			}
			key = strings.TrimPrefix(jsonTag, ",omitempty")
		}

		builder.httpRequest.queryParams[key] = fieldValue.Interface()
	}
	return builder
}

func (builder *HttpRequestBuilder) AddQueryParamByObject(input interface{}) *HttpRequestBuilder {
	if input == nil {
		return builder
	}

	v := reflect.ValueOf(input)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()
	if t.Kind() != reflect.Struct {
		return builder
	}
	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)

		if !fieldType.IsExported() {
			continue
		}

		// 忽略空指针
		if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
			continue
		}
		key := fieldType.Name
		jsonTag := fieldType.Tag.Get("json")
		if jsonTag != "" {
			if fieldValue.IsZero() && strings.Contains(jsonTag, ",omitempty") {
				continue
			}
			key = strings.TrimPrefix(jsonTag, ",omitempty")
		}

		builder.httpRequest.queryParams[key] = fieldValue.Interface()
	}
	return builder
}

func (builder *HttpRequestBuilder) AddPathParam(key string, value string) *HttpRequestBuilder {
	builder.httpRequest.pathParams[key] = value
	return builder
}

func (builder *HttpRequestBuilder) AddAutoFilledPathParam(key string, value string) *HttpRequestBuilder {
	builder.httpRequest.autoFilledPathParams[key] = value
	return builder
}

func (builder *HttpRequestBuilder) AddHeaderParam(key string, value string) *HttpRequestBuilder {
	builder.httpRequest.headerParams.Add(key, value)
	return builder
}

func (builder *HttpRequestBuilder) AddBasicAuthHeaderParam(user string, password string) *HttpRequestBuilder {
	auth := user + ":" + password
	authEncoded := base64.StdEncoding.EncodeToString([]byte(auth))
	builder.AddHeaderParam("Authorization", "Basic "+authEncoded)
	return builder
}

func (builder *HttpRequestBuilder) AddTokenAuthHeaderParam(token string) *HttpRequestBuilder {
	builder.AddHeaderParam("Authorization", "Bearer "+token)
	return builder
}

func (builder *HttpRequestBuilder) AddFormParam(key string, value def.FormData) *HttpRequestBuilder {
	builder.httpRequest.formParams[key] = value
	return builder
}

func (builder *HttpRequestBuilder) WithBody(kind string, body interface{}) *HttpRequestBuilder {
	// if body is multipart data, add to form
	if kind == "multipart" {
		v := reflect.ValueOf(body)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		t := reflect.TypeOf(body)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		// avoid panic
		if t.Kind() != reflect.Struct {
			return builder
		}

		fieldNum := t.NumField()
		for i := 0; i < fieldNum; i++ {
			jsonTag := t.Field(i).Tag.Get("json")
			if jsonTag != "" {
				if v.FieldByName(t.Field(i).Name).IsNil() && strings.Contains(jsonTag, "omitempty") {
					continue
				}
				builder.AddFormParam(
					strings.Split(jsonTag, ",")[0],
					v.FieldByName(t.Field(i).Name).Interface().(def.FormData),
				)
			} else {
				builder.AddFormParam(t.Field(i).Name, v.FieldByName(t.Field(i).Name).Interface().(def.FormData))
			}
		}
	} else {
		builder.httpRequest.bodyData = body
	}

	return builder
}

func (builder *HttpRequestBuilder) WithTimeout(t time.Duration) *HttpRequestBuilder {
	builder.httpRequest.timeout = t
	return builder
}

func (builder *HttpRequestBuilder) Build() *HttpRequest {
	return builder.httpRequest.fillParamsInPath()
}
