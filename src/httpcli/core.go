package httpcli

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/wangweihong/gotoolbox/src/callerutil"

	"github.com/wangweihong/gotoolbox/src/httpcli/httpconfig"
	"github.com/wangweihong/gotoolbox/src/log"
)

type Client struct {
	conn              *http.Client
	config            httpconfig.HttpConfig
	chainInterceptors []Interceptor
}

func NewClient(cfg *httpconfig.HttpConfig, options ...Option) (*Client, error) {
	config := httpconfig.DefaultHttpConfig()
	if cfg != nil {
		config = cfg
	}

	c := &Client{
		config: *config,
	}
	for _, o := range options {
		o(c)
	}

	tr := &http.Transport{}
	if c.config.HttpTransport != nil {
		tr = c.config.HttpTransport
	} else {
		creds, err := c.config.BuildCredentials()
		if err != nil {
			return nil, err
		}
		tr.Proxy = c.config.HttpProxy
		tr.TLSClientConfig = creds
	}

	c.conn = &http.Client{
		Transport: tr,
		Timeout:   c.config.Timeout,
	}

	return c, nil
}

type Interceptor func(ctx context.Context, req *HttpRequest, arg, reply interface{}, cc *Client, invoker Invoker, opts ...CallOption) (*HttpResponse, error)

func (c *Client) Invoke(
	/*
			注意1. ctx的生命周期,如果Client在一个服务请求的异步动作,不能直接使用服务请求的ctx。否则当服务器结束后,context撤销
		会导致所有Client的请求都会立即`context cancel`失败返回
			注意2. log fieldCtx的信息覆盖问题. 如果ctx来自一个服务请求,服务的中间件可能在该请求中设置了请求相关信息。如果client
		也使用这个ctx, 在记录信息，注意字段的覆盖和冗余问题.
	*/
	ctx context.Context,
	req *HttpRequest,
	arg, reply interface{},
	opts ...CallOption,
) (*HttpResponse, error) {
	file, line, fn := callerutil.CallerDepth(2)
	callerMsg := fmt.Sprintf("%s:%s:%d", file, fn, line)

	//  允许特定请求单独设置拦截器
	ci := &callInfo{}
	for _, o := range opts {
		o(ci)
	}
	//if ci.header != nil {
	//	for k := range ci.header {
	//		req.Builder().AddHeaderParam(k, ci.header.Get(k))
	//	}
	//}

	chainInterceptors := c.chainInterceptors
	if ci.chainInterceptors != nil {
		chainInterceptors = ci.chainInterceptors
	}

	if chainInterceptors != nil {
		rawResp, err := chainInterceptors[0](
			ctx,
			req,
			arg,
			reply,
			c,
			getChainUnaryInvoker(chainInterceptors, 0, invoke),
			opts...)

		log.F(ctx).L(ctx).
			Debug("Interceptor Invoked called.", log.String("caller", callerMsg), log.Err(err), log.Every("arg", arg), log.Every("reply", reply))
		return rawResp, err
	}

	rawResp, err := invoke(ctx, req, arg, reply, c, opts...)
	log.F(ctx).L(ctx).
		Debug("Invoked called.", log.String("caller", callerMsg), log.Err(err), log.Every("arg", arg), log.Every("reply", reply))
	return rawResp, err
}

func getChainUnaryInvoker(interceptors []Interceptor, curr int, finalInvoker Invoker) Invoker {
	if curr == len(interceptors)-1 {
		return finalInvoker
	}
	return func(ctx context.Context, req *HttpRequest, arg, reply interface{}, cc *Client, opts ...CallOption) (*HttpResponse, error) {
		return interceptors[curr+1](
			ctx,
			req,
			arg,
			reply,
			cc,
			getChainUnaryInvoker(interceptors, curr+1, finalInvoker),
			opts...)
	}
}

type Invoker func(ctx context.Context, req *HttpRequest, arg, reply interface{}, cc *Client, opt ...CallOption) (*HttpResponse, error)

func logEnabled() bool {
	debugEnv := os.Getenv("HTTPCLI_DEBUG")
	return debugEnv != "" && debugEnv != "0"
}

// nolint: funlen,gocognit
func invoke(
	ctx context.Context,
	req *HttpRequest,
	arg interface{},
	reply interface{},
	c *Client,
	opt ...CallOption,
) (*HttpResponse, error) {
	ci := &callInfo{}
	for _, o := range opt {
		o(ci)
	}

	// refer to https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	var cancel context.CancelFunc
	// 如果某个请求指定的超时时间, 则采用该超时时间
	if ci.timeout != 0 {
		ctx, cancel = context.WithTimeout(ctx, ci.timeout)
		defer cancel()
	}

	// 转换成原生http请求
	httpReq, err := req.ConvertRequestWithContext(ctx)
	if err != nil {
		return nil, err
	}

	// 请求前处理. 如一些场景需要对请求参数进行签名
	if err := c.preRequestProcess(httpReq, ci); err != nil {
		return nil, err
	}

	resp, err := c.conn.Do(httpReq)
	if err != nil {
		return nil, err
	}

	// 请求后处理
	if err := c.postRequestProcess(resp, ci); err != nil {
		return nil, err
	}

	return NewHttpResponse(req, resp), nil
}

func (c *Client) preRequestProcess(req *http.Request, info *callInfo) error {
	if info != nil && info.httpRequestProcess != nil {
		if _, err := info.httpRequestProcess(req); err != nil {
			return err
		}
	}

	if c.config.HttpHandler == nil || c.config.HttpHandler.RequestHandlers == nil || req == nil {
		return nil
	}
	return c.config.HttpHandler.RequestHandlers(req)
}

func (c *Client) postRequestProcess(resp *http.Response, info *callInfo) error {
	if c.config.HttpHandler == nil || c.config.HttpHandler.ResponseHandlers == nil || resp == nil {
		return nil
	}

	return c.config.HttpHandler.ResponseHandlers(resp)
}
