package httpcli

import (
	"net/http"
	"time"
)

type callInfo struct {
	timeout            time.Duration
	httpRequestProcess func(req *http.Request) (*http.Request, error)
	urlSetter          func() (string, error)
	httpTransport      *http.Transport
	// 拦截器列表
	chainInterceptors []Interceptor
}

type CallOption func(*callInfo)

type CallOptions []CallOption

func Combine(o1 []CallOption, o2 []CallOption) []CallOption {
	if len(o1) == 0 {
		return o2
	} else if len(o2) == 0 {
		return o1
	}
	ret := make([]CallOption, len(o1)+len(o2))
	copy(ret, o1)
	copy(ret[len(o1):], o2)
	return ret
}

func (cs CallOptions) Duplicate() []CallOption {
	if cs == nil {
		return nil
	}

	n := make([]CallOption, 0, len(cs))
	for _, v := range cs {
		n = append(n, v)
	}
	return n
}

// TimeoutCallOption 设置某个连接超时操作.
func TimeoutCallOption(timeout time.Duration) CallOption {
	return func(c *callInfo) {
		if timeout < 0 {
			return
		}
		c.timeout = timeout
	}
}

type ProcessRequestFunc func(req *http.Request) (*http.Request, error)

// HttpRequestProcessOption 在http请求发起调用前，对http请求进行处理. 如根据url/请求头进行加密,并写入httpReq.
func HttpRequestProcessOption(fun ProcessRequestFunc) CallOption {
	return func(c *callInfo) {
		c.httpRequestProcess = fun
	}
}

type URLSetter func() (string, error)

// 有可能需要根据资源/rawURL动态更改请求URL
func URLCallOption(epf URLSetter) CallOption {
	return func(c *callInfo) {
		c.urlSetter = epf
	}
}

// 更改访问的拦截器列表
func InterceptorsCallOption(chainInterceptors []Interceptor) CallOption {
	return func(c *callInfo) {
		c.chainInterceptors = chainInterceptors
	}
}

// CallOptionTransport 通用请求选项.
func CallOptionTransport(tp *http.Transport) CallOption {
	return func(c *callInfo) {
		c.httpTransport = tp
	}
}
