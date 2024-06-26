package grpctls_test

import (
	"testing"

	"github.com/wangweihong/gotoolbox/pkg/tls/grpctls"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	serverCA  = "-----BEGIN CERTIFICATE-----\nMIIFCTCCAvGgAwIBAgIUZSYstTIHWDtQI/rQs67/29i5n9IwDQYJKoZIhvcNAQEN\nBQAwFDESMBAGA1UEAwwJZWF6eWNsb3VkMB4XDTIzMDgwODA1NTEwM1oXDTMzMDgw\nNTA1NTEwM1owFDESMBAGA1UEAwwJZWF6eWNsb3VkMIICIjANBgkqhkiG9w0BAQEF\nAAOCAg8AMIICCgKCAgEA338ueymX3rZ29GN177ealPS6izUdWd6u6aKiPVY8ME8N\nxWZ6SAcAWg/iiO4vwTWADNP0r2Y4T8fodUKyIzVcl4Jy08UynwblAX+hHVZx+U5c\ncSq6gBVHoe1oNB4AiQqYog7ut710SeoD7bas1iepE3pUSr1GcyWpRVBOq3020Xxi\nV4XKRx80OfuHnqykR1UbE23DHexqd8+Q8UswivSBPdptA17u1YCz2HutzlSv4AIS\n1hW7peqTMlO08uCEwa8lFocp56tF3wi1khtkoyy3rvG+BIuATvpK7fl3km4Cozty\nyNU2AhE0L5W1+heevhw4Fh/iLNM/0EADmxIkL4Bc9JXJw/lz2vBNsedQ7dBKyBv0\n7FjcgDwU7X6zDx9OveU3NzleHHqBvg8AgGOJRQr5sU2EPvfTJmshGi5h6mR3Qrjr\nVxzMKuBte+aigI+Q2tYvpj2mARvD2zPj/cPj4Xkyjc3+W+ctXLYsMyzxkLYzbZJF\nMFvSeH/ybQbEMXXSuD1G70d9dmDnE7M4/tqXH1QJ0BiIQM2YjlIaDo/is7Yn12fg\noQoNoSontGzaX1hBgYFC7zOTwM8Nyim2ETEVQMPp/0t6Zs0rfJqK9n6QmfUTV1dW\n7QwV4DaP3kz6BBT08t6NTwEpYe+B7VJtWcPd+9zcFD3R/lrNG1QegduoSV/UC+0C\nAwEAAaNTMFEwHQYDVR0OBBYEFHyEAow4b/+3xm0Ljg1u9b5ZNYmcMB8GA1UdIwQY\nMBaAFHyEAow4b/+3xm0Ljg1u9b5ZNYmcMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZI\nhvcNAQENBQADggIBAMLfqMjdWP8fK93EdJ7kggdHnvM/yaAVNo6Mjyyxf4zwVEBi\nfQG7ejCb9lrMUkgj4Y30fMTwIG0gRKJIEogmMDJW+2ed/4ZWVE8oY3hChqV3HOYZ\nWr5Bt70G3oOp6aw11VqFr/b53zSlbGPGMzw4x0rjPotzMBPQ2M6S4KbfAx1CLJwl\n0G170MZxlsJkESB/Ha49ny3l49ElucmCtUdtaK5Pchs2kaTzN5OLEG1a5Msgv7Aa\npPQTNaA0WkXSZyVapxF0wWru9YUbL/GyumverlXYRBQMBDAU8q9j77wq/WXxaNo1\nu/qawcLadurcficgXj69Sgw7vObHDX2UBBxAr0g4nhJIPhQNV3PUHeWaV8udkhOt\nup6QV37nYwzfGiALcKmaSfs+gQsGUowTs9f7Og4qoPS0FPLrG8CfIXHu2EPAI9AZ\nO/Udv1uskXXCLPdnbidOooHfivf4/f6UGO9+Nv8doTwCLh2MDVmcm+kR9xepJW2g\nPYJHLt44VA5EuAuKFsLdCq7eoAyy36FrT5x/f+qjzErREXBF9LyE2hrr+FdQJLyo\nkASIOGN2TfLKzRNvSMt9mIXcHzu+sT7cApcYUCa/W3CDtnr57dI9jktYlxm+AO7U\nsfQQKfKSA/oYM0fzKUSZi80GeDsNpEhbQFfwvFmcYkBVEdTp9G4PrCHlBqP9\n-----END CERTIFICATE-----\n"
	serverCrt = "-----BEGIN CERTIFICATE-----\nMIIFRDCCAyygAwIBAgIUCyzw0Udh8o6pxdb+Mu9vWX18UpMwDQYJKoZIhvcNAQEN\nBQAwFDESMBAGA1UEAwwJZWF6eWNsb3VkMB4XDTIzMDgwODA3NDU0MloXDTMzMDgw\nNTA3NDU0MlowFDESMBAGA1UEAwwJZWF6eWNsb3VkMIICIjANBgkqhkiG9w0BAQEF\nAAOCAg8AMIICCgKCAgEAyf4H3XJGbt6cpfFRnQoxojwfplnnoXvxpTf5Y66byLrY\nFPAVbGUUx620xAWXvjLxmLzF7DMSOVHXPmxdRUV+nsfbGwWcjgZ606YqSWTen7NG\ngMBD/0M95ZnOSWDaOy2N9k1G2X4RooPw8ZWaJsqRjDw6GHImRZWKlIC9k4QJThOb\nMYcRyPeIHrorQOCBp0yTA0RxH1lJmkeInAFsTf3sG8cupsdeDpgeQN7MVliw3tM0\n9JXZdA/ip8hswvIkoLyNHMD8xlSMa9V+hIrOpDdD7vd4/8o8MAGcbUrzbyLUcMRY\nZjxpvZUVJqv690oheccKZhqvp94rTKpe2n3nBVwKLOJ45eCo8aPSCMEW7xN4R+mT\nB/XrD9zRozJO3EPrhv6p7oEiOIeWlvbGH44u5pCKJTmtM+fBUSpaRNsy+aSjo8QR\n14FF8Xs8Gfe0NHHA7k0dQuEddhMLaTWdTdNLkOjOHv4zH5UapnUyjxFLVaf/ic1D\nr/4miTDGpc1kkI5pzR+kQxGPKJbVrIDWP5ubbTYwQxIpEeoQit0kn0hm9I0x9AsW\nbHXTw/G82cY94w52UrwMeDkHQGvnknlYhrUN3QmDk7KN7g5TBg17JRjctGeO76e2\ns4iuShHGUcv2bbQCqHGiDQdKSvQIX5gZMdtIJzGcPhcfnfUarD+b0B8kI9FVlJUC\nAwEAAaOBjTCBijAfBgNVHSMEGDAWgBR8hAKMOG//t8ZtC44NbvW+WTWJnDAJBgNV\nHRMEAjAAMAsGA1UdDwQEAwIE8DATBgNVHSUEDDAKBggrBgEFBQcDATAbBgNVHREE\nFDAShwQAAAAAhwR/AAABhwTAqIaLMB0GA1UdDgQWBBRoAkL/11IAOtTs06ggbcH2\nZXkMoTANBgkqhkiG9w0BAQ0FAAOCAgEAsic3hIaRTDsn9rUUEfdQORxOj+utxA21\nqomT+3AIh8SYWb6QMiCukAKgZ1pckjHxCWlbgChESjEtN8jaegmYO1+zLaRcN1BA\nvHHv/kjzYqgSUvyOk4pnpC7KALajsrOj994ZK0MvO8H7zE1BRVBSjhT02SOR4QJz\n7zuNH2Jy1SnpFWDy4YDtmfwk4uMdTA2azjojZ7wgbDENbwtiQazeGYkCd6qykLt4\n3jmvYkQ+sZHkJUsb8SB+YDTC8Q/4fkmRYUHVhJ8ft5DVGTEkn8iyjiUgnr4+Cxt9\n9oSMoZ5boJFv9pMTJkT3+bF8G60tMzvJEbGCUVlGeXMt2VP95WMwGnWHAhkfaHDJ\nol/qTZb+kew2mWp9mvpsEORX3PkYsX3EAoC9FTLPQ6O9A7jcWJK4H85EGxVpFkwD\n2WXJ/Ud0h5WG4fJN958DKQByCdHi1ii2D5KZIauzE93wmy0v3mUG78nhkPA8/wJa\nngvUpC6gI5An1bwkEpGPefrsCLIHw4ZXuE0vzghXvSRcWbhqh+tVfSx3soM6Y1zF\nNwPycBXk/N3oltXm91BYxV2zsOygvRveJm/yHiDqx9jOOf6pZ0Rf5rMIt+79Om/A\nPVXuk5qBnRxyc3uo+8UQkob5WVj62HpAV1dlp3Hd2uZ/RxhTL88/9R9nc4d/0WgA\npERLDf7APaQ=\n-----END CERTIFICATE-----\n"
	serverKey = "-----BEGIN PRIVATE KEY-----\n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\nmLbgErgFK6kGGHs9jn72WmY0ysTAE9g=\n-----END PRIVATE KEY-----\n"
	clientCA  = serverCA
	clientCrt = "-----BEGIN CERTIFICATE-----\nMIIFJTCCAw2gAwIBAgIUM5xEnChPMrUIIYlmJ/5T6+BXrygwDQYJKoZIhvcNAQEN\nBQAwFDESMBAGA1UEAwwJZWF6eWNsb3VkMB4XDTIzMDgwODA3NDU0NVoXDTMzMDgw\nNTA3NDU0NVowFDESMBAGA1UEAwwJZWF6eWNsb3VkMIICIjANBgkqhkiG9w0BAQEF\nAAOCAg8AMIICCgKCAgEArkE9FdQl04Z6N0fa26jnzOnrwc+8VBdM/b88XZBrbi1M\nwO0T5pwWWRRYjowbbCLOUxKh6Mq96di8Rhd6/8Qkd3T0Uych0sgcMlhAjF3MG3Ow\nQs8BKmOkT1sVNHHBCb+rNbFS5yB/EcJDb2MPyqB8/AEv34IrE8v+lBNeicZMl79A\nPtUZoCD3gg54MvLfxotK3Qm3tYyFwDfBsgxK5KYKuFfnY+9bOg4yJ87Prwo5IF2g\nDp/1MUyUbWAFHntGL1BR8PoNqE2dF3635innOCq1MGEEWvXallOztcZncGwu10u5\nzTZ60EqOmqyb4mpJQZtwpRIW6WvBbfHdXbljJxnbBU/210yd2vnLdrV1QDNLW0OT\nptUk5sYG0s0PV/iEaUBaK2v2/gPtCB0xPs2ODAs1CcBRzhqsIadFSU92TdfGQGZF\nCokGvfdiUgQCAYXMYeGgvbsQvx0+jIwdpP+gg8SrlusPGQIN1Q5f871+3VIfZhGL\nZ8yDnVFAzDdlZJf4nCQrRhUYNnmykX8NnzQV58eRK+6Tjmm9EBfkTxmo5NNMXHO/\nYOKYuqTMpb7/5o+yLShlN3ZMvx5gUbTBe1srnltIrqw7hedWqd02zvBhZfWBAY0z\n/x+Lzl0It3pUhds4B96kiSX5KaI5x6zaUJJwSXZmwERoaZwq8Z5hJ8PZNWSojGcC\nAwEAAaNvMG0wHwYDVR0jBBgwFoAUfIQCjDhv/7fGbQuODW71vlk1iZwwCQYDVR0T\nBAIwADALBgNVHQ8EBAMCBPAwEwYDVR0lBAwwCgYIKwYBBQUHAwIwHQYDVR0OBBYE\nFOhQTNQW/8qlZ97FBETzAmn9j5miMA0GCSqGSIb3DQEBDQUAA4ICAQBoaHa/Y23i\nsqpLrkMnz6wglgTSxfeT9GlhLonZLepMcr1pJyidcee1iFkcj+qhgN/Bdrk8aFp6\nlU+3mGBr5U64F2M/LeZqaEO/miAib0Ags5L3LGdMRpOBGIqEx9nMqs8vrfApBka0\n3zHYMe+5hypuDluZVHnQhgY6Orc9T2s4ue0tIRCBABSG8id7t7sD+cV2/DnHGd+0\nM4TOgKUi9Z+wvE6FLHPyxuD2aTbcTmjBTxnlti+l4gnYITuhh26zqg4DP9OCDIpU\n6d/ApfE6Y3Z54uHHOf8ke1jugv/G+e9khR4qqkYYlB0QBxaM7jRWIpPvu1vxlpUV\nh1C53QHNm9c2B48WB1f9OVHICG33UQ2Jnc8ohgcGtks+4jaXcxUym9boGlHDs9ks\nOJpVv0rBQDyRkKyvn9ROh8NklRw+avkU6w0XhJY7Nr8VFJUbt73Wtmlt1wYQBDpC\nK1OLIADOP1ksLL0dqiNETcO/QmGvW8SeeYPftZkdXAjI9PeAEyBLQJIj/a9LB6O1\nkWAd5NivGQiD/ub9545U/MDWMJm2Y5icfFVejQPYaeJUfv8jD/6wdXdxfdPDIbhv\ngphx5MPJl2BO5D7tg3XaxzY450HNvXU0SGPcxUjgh83RsZW9Naf6SL2JMG20EODz\n9laevl99rio9zMa7slKZXmAY/w5dXLf/bw==\n-----END CERTIFICATE-----\n"
	clientKey = "-----BEGIN PRIVATE KEY-----\n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\nhahcLECnUocOJcbv/bmZ2dRL7oJlDsk=\n-----END PRIVATE KEY-----\n"
	wrongCA   = "-----BEGIN CERTIFICATE-----\nMIIFCTCCAvGgAwIBAgIUKXAB9kIZyN0r8PhR+wbKo5aKraQwDQYJKoZIhvcNAQEN\nBQAwFDESMBAGA1UEAwwJZWF6eWNsb3VkMB4XDTIzMDgwOTAzMTYyMFoXDTMzMDgw\nNjAzMTYyMFowFDESMBAGA1UEAwwJZWF6eWNsb3VkMIICIjANBgkqhkiG9w0BAQEF\nAAOCAg8AMIICCgKCAgEA5PcNkVaOiGmbA30P0NvL2/9+wNrRJD4NdVeHucLQwZ+3\n8ErOA8oiexTExlUFlIzSmLKciLIrwanyMCE/r4/dUCs8pQl+E3jPjT6sKRF4BtbN\nQFBZhzPWa0Ia8FHZ0D7wM6D9duA63UlvyeqK25ChjZC1FX7vOIyyskZWZfP9I6jU\nSrE+B6tXkXUorgaQSz6bmqMBiUM+v8R52XFC/ucwWFFzmI52oG/utR104/a97t5s\n7OpobMlMz2Ll0tsxg0tKiU9nlwKIgMaHHP3R7MMusLDkcoLlzdQNGHXM+c3usShH\n9bipKVm7KbKiDqCuBm91aCTcd5sMDiS+oKyorgrGGDJNHUC6pNoFHNb/k5FOmiz6\n1lfmgSM9R7FrmnHr6bkahYOUClPsnd9IHzyzJmekncbJsx2mlNTZGXYe0SNzMUau\nKxZqNWZfEbVOo7i2IB/688XZMo7srkHeHn6Y67h9PPNB6oR0UPsk68ZImip4CrQG\n9sZkNX73Ujkxeq57xLNJaOqWvq8xfIBbNchGgEfPIalZ07ei5hjzX6pEOlrDWV0c\n2VOCPlTEypFr9rsqXw++zsTqVYiRLLG4cmRKEBHnaNeELU8IYpQuFsnqGWpCpwz3\noHmemW/EY25k0luk7KmaLI7cIH9XgGca9VfnsneSzi4XcnW+aYIWJoGcPVro2qkC\nAwEAAaNTMFEwHQYDVR0OBBYEFAzV6u2lJUEQ02P++kz9BYhl7vaiMB8GA1UdIwQY\nMBaAFAzV6u2lJUEQ02P++kz9BYhl7vaiMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZI\nhvcNAQENBQADggIBAC2XHwAjuU8sxqABduSfrhRwEfFDwqtMlOw8mtmLqkOI8w50\nEelIV2EfdVkpbU6wGEoLDFTJm17BrcR3SiJp194ZH6h3Qof1t/dlSFWssTGWdbFA\ncJaN9TWc5NzxXvCddV0clIzW8jZc5rRFgkAU+/+yvd17iStf7j20ON8qZ4JriRI3\neiXT/XOfz3sWf0qtqqLjJrbp4XSX60axxKPRiGq3G1UBI1WcRvdpmKYPjno5YS1Q\n4ND6WqjlNBGg0ANwthm7V7RQNgvGg45Jt/Lw1cjLPo4qvzC9c2b6Oo08AUleRV6q\nPFPgaC/lGyDRQcquQrZxuJxagO5EWyv3phTqKJLnpNEAcxdX4J70GP1Qu/WkCncT\nuaSL/j35dX11HXCiDeUTOx6VGCKGtQ7FBu+sm4TEK2BAgfskm8DAJYeIp+vQ7Hmp\n7y60zNxT+pg7eydZx5FSmeyyMD6g67sDQ7zb9XDjDpPOsN1uOhOcF0PsiSZuvnGn\nxKvzfkn9tTA5W46RMIjj3PFqkIMKbY9KxzzM2aw0CTsvAGbB6Sj1y1dYHurEfea4\nkg8javdyIuNZZklTQfjoviDrpum0zxh7NcHAkkWgRCyVkiwBsiBm5BlMxxdnmg3l\nTG13I5JG3WS+qTvaKwOd5fYB+JyJqyyAYvociyKykj/dAQ5w7dlWCwAJqdnW\n-----END CERTIFICATE-----\n"
)

func Test_NewTlsServerCredentials(t *testing.T) {
	Convey("NewTlsServerCredentials", t, func() {
		Convey("正确证书密钥", func() {
			_, err := grpctls.NewTlsServerCredentials([]byte(serverCrt), []byte(serverKey))
			So(err, ShouldBeNil)
		})

		Convey("错误证书密钥", func() {
			var err error
			// 证书密钥相反
			_, err = grpctls.NewTlsServerCredentials([]byte(serverKey), []byte(serverCA))
			So(err, ShouldNotBeNil)

			// 证书为空
			_, err = grpctls.NewTlsServerCredentials([]byte(nil), []byte(serverKey))
			So(err, ShouldNotBeNil)

			// 密钥为空
			_, err = grpctls.NewTlsServerCredentials([]byte(serverCrt), []byte(nil))
			So(err, ShouldNotBeNil)

			// 证书和密钥不匹配
			_, err = grpctls.NewTlsServerCredentials([]byte(serverCrt), []byte(clientKey))
			So(err, ShouldNotBeNil)

			// 错误证书格式
			_, err = grpctls.NewTlsServerCredentials([]byte("xxxx"), []byte(serverKey))
			So(err, ShouldNotBeNil)

			_, err = grpctls.NewTlsServerCredentials([]byte(serverCrt), []byte("serverKey"))
			So(err, ShouldNotBeNil)
		})
	})
}

func Test_NewMutualTlsServerCredentials(t *testing.T) {
	Convey("NewMutualTlsServerCredentials", t, func() {
		Convey("正确证书密钥", func() {
			_, err := grpctls.NewMutualTlsServerCredentials([]byte(clientCA), []byte(serverCrt), []byte(serverKey))
			So(err, ShouldBeNil)
		})

		Convey("错误证书密钥", func() {
			var err error
			// 客户端CA证书为nil
			_, err = grpctls.NewMutualTlsServerCredentials([]byte(nil), []byte(serverCrt), []byte(serverKey))
			So(err, ShouldNotBeNil)

			// 服务端证书密钥为空
			_, err = grpctls.NewMutualTlsServerCredentials([]byte(clientCA), []byte(nil), []byte(nil))
			So(err, ShouldNotBeNil)

			// 错误ca格式
			_, err = grpctls.NewMutualTlsServerCredentials([]byte("clientCA"), []byte(serverCrt), []byte(serverKey))
			So(err, ShouldNotBeNil)

			// 错误证书格式
			_, err = grpctls.NewMutualTlsServerCredentials([]byte(clientCA), []byte("xxxx"), []byte(serverKey))
			So(err, ShouldNotBeNil)

			_, err = grpctls.NewMutualTlsServerCredentials([]byte(clientCA), []byte(serverCrt), []byte("serverKey"))
			So(err, ShouldNotBeNil)

		})
	})
}
