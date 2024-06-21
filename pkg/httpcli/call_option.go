package httpcli

import (
	"net/http"
	"net/url"
	"time"
)

type callInfo struct {
	timeout            time.Duration
	httpRequestProcess func(req *http.Request) (*http.Request, error)
	httpTransport      *http.Transport
	// 拦截器列表
	chainInterceptors []Interceptor

	TlsEnabled       bool
	SkipTlsVerified  bool
	ServerCA         string
	MutualTlsEnabled bool
	ClientKeyData    string
	ClientCertData   string

	HttpProxy func(*http.Request) (*url.URL, error)
}

type MTLS struct {
}

type CallOption func(*callInfo)

// TimeoutCallOption 设置某个连接超时操作.
func TimeoutCallOption(timeout time.Duration) CallOption {
	return func(c *callInfo) {
		if timeout < 0 {
			return
		}
		c.timeout = timeout
	}
}

// CallOptionInsecure 是否跳过服务端证书检测.
func CallOptionInsecure() CallOption {
	return func(c *callInfo) {
		c.TlsEnabled = true
		c.SkipTlsVerified = true
	}
}

// ServerCACallOption 设置服务端CA证书数据.
func ServerCACallOption(serverCAData string) CallOption {
	return func(c *callInfo) {
		c.TlsEnabled = true
		c.ServerCA = serverCAData
	}
}

// MTLSCallOption 是否开启双向认证.
func MTLSCallOption(serverCAData string, clientCertData string, clientKeyData string) CallOption {
	return func(c *callInfo) {
		c.MutualTlsEnabled = true
		c.TlsEnabled = true
		c.ClientCertData = clientCertData
		c.ClientKeyData = clientKeyData
		c.ServerCA = serverCAData
	}
}

type ProcessRequestFunc func(req *http.Request) (*http.Request, error)

// HttpRequestProcessOption 在http请求发起调用前，对http请求进行处理. 如根据url/请求头进行加密,并写入httpReq.
func HttpRequestProcessOption(fun ProcessRequestFunc) CallOption {
	return func(c *callInfo) {
		c.httpRequestProcess = fun
	}
}

// CallOptionTransport 通用请求选项.
func CallOptionTransport(tp *http.Transport) CallOption {
	return func(c *callInfo) {
		c.httpTransport = tp
	}
}

// CallOptionProxy 请求代理选项.
// WithProxy(http.ProxyFromEnvironment).
func CallOptionProxy(proxy func(*http.Request) (*url.URL, error)) CallOption {
	return func(c *callInfo) {
		c.HttpProxy = proxy
	}
}
