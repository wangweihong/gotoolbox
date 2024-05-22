package httpconfig

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/wangweihong/gotoolbox/src/tls/httptls"

	"github.com/wangweihong/gotoolbox/src/httpcli/httphandler"
)

const (
	DefaultTimeout = 120 * time.Second
)

type HttpConfig struct {
	// Timeout 设置整个客户端的超时, 优先级低于每个请求单独的请求
	Timeout          time.Duration
	HttpProxy        func(*http.Request) (*url.URL, error)
	HttpHandler      *httphandler.HttpHandler
	HttpTransport    *http.Transport
	TlsEnabled       bool
	SkipTlsVerified  bool
	ServerCA         string
	MutualTlsEnabled bool
	ClientKeyData    string
	ClientCertData   string
}

func DefaultHttpConfig() *HttpConfig {
	return &HttpConfig{
		Timeout: DefaultTimeout,
	}
}

func (c *HttpConfig) Validate() error {
	if c.TlsEnabled {
		if !c.SkipTlsVerified {
			if c.ServerCA == "" {
				return fmt.Errorf("must set serverCA when tlsEnabled and not skipTlsVerified")
			}
		}
	}

	if c.MutualTlsEnabled {
		if c.ClientKeyData == "" || c.ClientCertData == "" {
			return fmt.Errorf("must provide clientKeyPEMData and clientCertPEMData when enable mTls")
		}

		if c.ServerCA == "" {
			return fmt.Errorf("must set serverCA when mtlsEnabled enable")
		}
	}
	return nil
}

func (c *HttpConfig) BuildCredentials() (*tls.Config, error) {
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
