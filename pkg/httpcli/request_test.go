package httpcli_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/wangweihong/gotoolbox/pkg/httpcli"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHttpRequest_GetURL(t *testing.T) {
	Convey("TestHttpRequest_GetURL", t, func() {
		Convey("不带协议", func() {
			r := httpcli.NewHttpRequestBuilder().
				WithEndpoint("www.example.com").
				WithPath("/{path}").
				AddQueryParam("q", 123).AddQueryParam("q2", true).
				AddPathParam("path", "abc").Build()
			So(r.GetFullRequestAddress(), ShouldEqual, "http://www.example.com/abc?q=123&q2=true")
		})
		Convey("http协议", func() {
			r := httpcli.NewHttpRequestBuilder().
				WithEndpoint("www.example.com").
				WithPath("/{path}").
				AddQueryParam("q", 123).AddQueryParam("q2", true).
				AddPathParam("path", "abc").Build()
			So(r.GetFullRequestAddress(), ShouldEqual, "http://www.example.com/abc?q=123&q2=true")
		})
		Convey("https协议", func() {
			r := httpcli.NewHttpRequestBuilder().
				WithEndpoint("https://www.example.com").
				WithPath("/{path}").
				AddQueryParam("q", 123).AddQueryParam("q2", true).
				AddPathParam("path", "abc").Build()
			So(r.GetFullRequestAddress(), ShouldEqual, "https://www.example.com/abc?q=123&q2=true")
		})
	})
}

func TestHttpRequest_Invoke(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		// 模拟响应
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, world!"))
	}))
	defer server.Close()

	Convey("TestHttpRequest_Invoke", t, func() {
		Convey("no timeout", func() {
			resp, err := httpcli.NewHttpRequestBuilder().
				GET().
				WithEndpoint(server.URL).
				Build().
				Invoke()
			So(err, ShouldBeNil)
			So(resp.GetBody(), ShouldEqual, "Hello, world!")
		})
		Convey("timeout", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			resp, err := httpcli.NewHttpRequestBuilder().
				GET().
				WithEndpoint(server.URL).
				Build().
				InvokeWithContext(ctx)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "context deadline exceeded")
			So(resp, ShouldBeNil)
		})
	})
}
