package httpcli

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/wangweihong/gotoolbox/pkg/tls/httptls"

	"github.com/wangweihong/gotoolbox/pkg/httpcli/def"
	"github.com/wangweihong/gotoolbox/pkg/typeutil"
)

type HttpRequest struct {
	endpoint string
	path     string
	method   string

	queryParams  map[string]interface{}
	pathParams   map[string]string
	headerParams http.Header
	formParams   map[string]def.FormData
	bodyData     interface{}

	autoFilledPathParams map[string]string
	timeout              time.Duration
}

// 填充路径参数.
func (r *HttpRequest) fillParamsInPath() *HttpRequest {
	for key, value := range r.pathParams {
		r.path = strings.ReplaceAll(r.path, "{"+key+"}", value)
	}
	for key, value := range r.autoFilledPathParams {
		r.path = strings.ReplaceAll(r.path, "{"+key+"}", value)
	}
	return r
}

// Builder转换成构建器,用于修改请求.
func (r *HttpRequest) Builder() *HttpRequestBuilder {
	httpRequestBuilder := HttpRequestBuilder{httpRequest: r}
	return &httpRequestBuilder
}

func (r *HttpRequest) GetMethod() string {
	return r.method
}

func (r *HttpRequest) GetEndpoint() string {
	return r.endpoint
}

func (r *HttpRequest) GetPath() string {
	return r.path
}

func (r *HttpRequest) GetQueryParams() map[string]interface{} {
	return r.queryParams
}

func (r *HttpRequest) GetHeaderParams() http.Header {
	return r.headerParams
}

func (r *HttpRequest) GetPathPrams() map[string]string {
	return r.pathParams
}

func (r *HttpRequest) GetFormPrams() map[string]def.FormData {
	return r.formParams
}

func (r *HttpRequest) GetFullRequestAddress() string {
	req, err := r.ConvertRequest()
	if err != nil {
		return ""
	}
	return req.URL.String()
}

func (r *HttpRequest) GetBodyData() interface{} {
	return r.bodyData
}

func (r *HttpRequest) GetBodyToBytes() (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}

	if r.bodyData != nil {
		if str, ok := r.bodyData.(json.RawMessage); ok {
			buf.WriteString(string(str))
		} else {
			v := reflect.ValueOf(r.bodyData)
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}

			if v.Kind() == reflect.String {
				buf.WriteString(v.Interface().(string))
			} else {
				var err error
				if r.headerParams.Get("Content-Type") == "application/xml" {
					encoder := xml.NewEncoder(buf)
					err = encoder.Encode(r.bodyData)
				} else {
					encoder := json.NewEncoder(buf)
					encoder.SetEscapeHTML(false)
					err = encoder.Encode(r.bodyData)
				}
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return buf, nil
}

func (r *HttpRequest) GetTimeout() time.Duration {
	return r.timeout
}

func (r *HttpRequest) ConvertRequestWithContext(ctx context.Context) (*http.Request, error) {
	t := reflect.TypeOf(r.bodyData)
	if t != nil && t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	var req *http.Request
	var err error

	// 1. 如果bodyData的类型为File, 则请求体读取流数据。常用于传输简单二进制流
	// 2. 如果是表单数据,则请求体转换成表单数据。常用于表单提交, 如需要上传文件,并携带一些文本字段
	// 3. 其他类型的请求体
	if r.bodyData != nil && t != nil && t.Name() == "File" {
		req, err = r.convertStreamBody(ctx)
		if err != nil {
			return nil, err
		}
	} else if len(r.GetFormPrams()) != 0 {
		req, err = r.covertFormBody(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		var buf *bytes.Buffer

		buf, err = r.GetBodyToBytes()
		if err != nil {
			return nil, err
		}

		req, err = http.NewRequestWithContext(ctx, r.GetMethod(), r.GetEndpoint(), buf)
		if err != nil {
			return nil, err
		}
	}
	r.fillPath(req)
	r.fillQueryParams(req)
	r.fillHeaderParams(req)

	return req, nil
}

// ConvertRequest convert to raw http request.
func (r *HttpRequest) ConvertRequest() (*http.Request, error) {
	return r.ConvertRequestWithContext(context.Background())
}

func (r *HttpRequest) covertFormBody(ctx context.Context) (*http.Request, error) {
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)

	sortedKeys := make([]string, 0, len(r.GetFormPrams()))
	for k, v := range r.GetFormPrams() {
		if _, ok := v.(*def.FilePart); ok {
			sortedKeys = append(sortedKeys, k)
		} else {
			sortedKeys = append([]string{k}, sortedKeys...)
		}
	}

	for _, k := range sortedKeys {
		if err := r.GetFormPrams()[k].Write(bodyWriter, k); err != nil {
			return nil, err
		}
	}

	contentType := bodyWriter.FormDataContentType()
	if err := bodyWriter.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, r.GetMethod(), r.GetEndpoint(), bodyBuffer)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-type", contentType)
	return req, nil
}

func (r *HttpRequest) convertStreamBody(ctx context.Context) (*http.Request, error) {
	f, ok := r.bodyData.(os.File)
	if !ok {
		return nil, errors.New("failed to get stream request body")
	}
	var reader io.Reader
	reader = &f
	req, err := http.NewRequestWithContext(ctx, r.GetMethod(), r.GetEndpoint(), reader)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (r *HttpRequest) fillHeaderParams(req *http.Request) {
	if len(r.GetHeaderParams()) == 0 {
		return
	}

	for key, values := range r.GetHeaderParams() {
		if strings.EqualFold(key, "Content-type") && req.Header.Get("Content-type") != "" {
			continue
		}
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
}

func (r *HttpRequest) fillQueryParams(req *http.Request) {
	if len(r.GetQueryParams()) == 0 {
		return
	}

	q := req.URL.Query()
	for key, value := range r.GetQueryParams() {
		q.Add(key, typeutil.ConvertInterfaceToString(value))
	}

	req.URL.RawQuery = strings.ReplaceAll(strings.ReplaceAll(strings.Trim(q.Encode(), "="), "=&", "&"), "+", "%20")
}

func (r *HttpRequest) fillPath(req *http.Request) {
	if r.GetPath() != "" {
		req.URL.Path = r.GetPath()
	}
}

// ConvertRequest convert to raw http request.
func (r *HttpRequest) Invoke(opts ...CallOption) (*HttpResponse, error) {
	return r.InvokeWithContext(context.Background(), opts...)
}

func (r *HttpRequest) InvokeWithContext(ctx context.Context, opts ...CallOption) (*HttpResponse, error) {
	ci := &callInfo{}
	for _, o := range opts {
		o(ci)
	}

	httpReq, err := r.ConvertRequestWithContext(ctx)
	if err != nil {
		return nil, err
	}

	tr := &http.Transport{}
	if ci.httpTransport != nil {
		tr = ci.httpTransport
	}

	c := http.Client{
		Transport: tr,
		Timeout:   ci.timeout,
	}

	if ci.TlsEnabled {
		creds, err := buildCredentials(ci)
		if err != nil {
			return nil, err
		}
		tr.TLSClientConfig = creds
	}
	tr.Proxy = ci.HttpProxy

	resp, err := c.Do(httpReq)
	if err != nil {
		return nil, err
	}
	return NewHttpResponse(r, resp), nil
}

func buildCredentials(c *callInfo) (*tls.Config, error) {
	var creds *tls.Config
	if c.TlsEnabled {
		var err error
		if c.SkipTlsVerified {
			creds = httptls.NewTlsClientSkipVerifiedCredentials()
		} else {
			if c.MutualTlsEnabled {
				// 如果开启双向认证,需要加载服务器
				creds, err = httptls.NewMutualTlsClientCredentials(
					[]byte(c.ServerCA),
					[]byte(c.ClientCertData),
					[]byte(c.ClientKeyData))
			} else {
				creds, err = httptls.NewTlsClientCredentials([]byte(c.ServerCA))
			}
		}
		if err != nil {
			return nil, err
		}
	}
	return creds, nil
}
