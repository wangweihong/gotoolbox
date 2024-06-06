package interceptorcli

import (
	"context"
	"fmt"
	"net/http"

	"github.com/wangweihong/gotoolbox/pkg/httpcli"
	"github.com/wangweihong/gotoolbox/pkg/log"
	"github.com/wangweihong/gotoolbox/pkg/skipper"
)

func StatusCodeInterceptor(name string, skipperFunc ...skipper.SkipperFunc) httpcli.Interceptor {
	return httpcli.NewInterceptor(name, func(ctx context.Context, req *httpcli.HttpRequest, arg, reply interface{}, cc *httpcli.Client,
		invoker httpcli.Invoker, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
		if skipper.Skip(req.GetPath(), skipperFunc...) {
			log.F(ctx).Debugf("skip interceptor %s for rawrurl %s", name, req.GetPath())

			return invoker(ctx, req, arg, reply, cc, opts...)
		}
		rawResp, err := invoker(ctx, req, arg, reply, cc, opts...)
		if err != nil {
			return rawResp, err
		}

		if rawResp.GetStatusCode() != http.StatusOK {
			return rawResp, fmt.Errorf("status code not 200, is %v", rawResp.GetStatusCode())
		}

		return rawResp, nil
	})
}
