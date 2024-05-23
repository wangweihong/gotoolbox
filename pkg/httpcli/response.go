package httpcli

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/wangweihong/gotoolbox/pkg/httpcli/decode"
	"github.com/wangweihong/gotoolbox/pkg/json"
)

type HttpResponse struct {
	Response *http.Response
	Request  *HttpRequest
}

func NewHttpResponse(req *HttpRequest, response *http.Response) *HttpResponse {
	return &HttpResponse{
		Response: response,
		Request:  req,
	}
}

func (r *HttpResponse) GetStatusCode() int {
	return r.Response.StatusCode
}

func (r *HttpResponse) GetHeaders() map[string]string {
	headerParams := map[string]string{}
	for key, values := range r.Response.Header {
		if values == nil || len(values) <= 0 {
			continue
		}
		headerParams[key] = values[0]
	}
	return headerParams
}

func (r *HttpResponse) GetBody() string {
	body, err := ioutil.ReadAll(r.Response.Body)
	if err != nil {
		return ""
	}

	// 1. 将HTTP响应主体替换成buffer, 读取HTTP响应主题数据后，关闭HTTP响应主体避免泄露
	// 2. 其次允许多次读取响应主体的内容
	if err := r.Response.Body.Close(); err == nil {
		r.Response.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	}
	return string(body)
}

func (r *HttpResponse) GetHeader(key string) string {
	header := r.Response.Header
	return header.Get(key)
}

func (r *HttpResponse) Decode(resp interface{}) error {
	if resp == nil {
		return nil
	}

	data := r.GetBody()
	if data == "" {
		return fmt.Errorf("body data is empty")
	}
	byteData := []byte(data)
	mm := decode.NewMarshalMapping()
	_ = mm.Register(decode.ApplicationJson, json.Unmarshal)
	_ = mm.Register(decode.ApplicationXml, xml.Unmarshal)
	return mm.UnmarshalManifest(r.GetHeader(decode.ContentType), byteData, resp)
}
