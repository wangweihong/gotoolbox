package httpcli_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wangweihong/gotoolbox/src/maputil"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/wangweihong/gotoolbox/src/httpcli"
)

func TestClient_Interceptor(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body != nil {
			defer r.Body.Close()

			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if _, err := w.Write(b); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	Convey("拦截器", t, func() {
		Convey("拦截器添加查询参数，并修改返回头部", func() {
			inter1 := func(ctx context.Context, req *httpcli.HttpRequest, arg, reply interface{}, cc *httpcli.Client, invoker httpcli.Invoker, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
				req.Builder().AddQueryParam("inter1", "aaaa").Build()

				resp, err := invoker(ctx, req, arg, reply, cc, opts...)
				if err != nil {
					return resp, err
				}
				resp.Response.Header.Set("inter1", "bbbb")
				return resp, err
			}
			inter2 := func(ctx context.Context, req *httpcli.HttpRequest, arg, reply interface{}, cc *httpcli.Client, invoker httpcli.Invoker, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
				req.Builder().AddQueryParam("inter2", "bbbb").Build()

				ctx = context.WithValue(ctx, "inter2", "bbbb")
				resp, err := invoker(ctx, req, arg, reply, cc, opts...)
				if err != nil {
					return resp, err
				}
				resp.Response.Header.Set("inter2", "bbbb")
				return resp, err
			}

			c, err := httpcli.NewClient(nil, httpcli.WithIntercepts(inter1, inter2))
			So(err, ShouldBeNil)
			ctx := context.Background()

			req := httpcli.NewHttpRequestBuilder().
				WithEndpoint(server.URL).
				WithMethod("GET").
				WithPath("/version").
				Build()
			resp, err := c.Invoke(ctx, req, nil, nil)
			So(err, ShouldBeNil)
			So(resp.Response.StatusCode, ShouldEqual, 200)
			So(resp.Response.Header.Get("inter1"), ShouldEqual, "bbbb")
			So(resp.Response.Header.Get("inter2"), ShouldEqual, "bbbb")
			So(maputil.StringInterfaceMap(resp.Request.GetQueryParams()).HasKeyAndValue("inter2", "bbbb"), ShouldBeTrue)
			So(maputil.StringInterfaceMap(resp.Request.GetQueryParams()).HasKeyAndValue("inter1", "aaaa"), ShouldBeTrue)
		})
	})
}
