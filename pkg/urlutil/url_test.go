package urlutil_test

import (
	"fmt"
	"net/url"
	"strconv"
	"testing"

	"github.com/wangweihong/gotoolbox/pkg/urlutil"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSplitURL_BUILDURL(t *testing.T) {
	Convey("SplitURL", t, func() {
		SkipConvey("http://127.0.0.1:1999", func() {
			schema, ip, port, err := urlutil.SplitURL("http://127.0.0.1:1999")
			So(err, ShouldBeNil)
			So(schema, ShouldEqual, "http")
			So(ip, ShouldEqual, "127.0.0.1")
			So(port, ShouldEqual, "1999")
			So(urlutil.BuildURL(schema, ip, strconv.Itoa(port)), ShouldEqual, "http://127.0.0.1:1999")
		})
		SkipConvey("127.0.0.1", func() {
			schema, ip, port, err := urlutil.SplitURL("127.0.0.1")
			So(err, ShouldBeNil)
			So(schema, ShouldEqual, "http")
			So(ip, ShouldEqual, "127.0.0.1")
			So(port, ShouldEqual, "80")
			So(urlutil.BuildURL(schema, ip, strconv.Itoa(port)), ShouldEqual, "http://127.0.0.1")

		})
		SkipConvey("https://127.0.0.1", func() {
			schema, ip, port, err := urlutil.SplitURL("https://127.0.0.1")
			So(err, ShouldBeNil)
			So(schema, ShouldEqual, "https")
			So(ip, ShouldEqual, "127.0.0.1")
			So(port, ShouldEqual, "443")
			So(urlutil.BuildURL(schema, ip, strconv.Itoa(port)), ShouldEqual, "https://127.0.0.1")
		})

		Convey("https://127.0.0.1/rest/v2", func() {
			u, _ := url.Parse("https://127.0.0.1/rest/v2")
			fmt.Println("host", u.Host)
			fmt.Println("path", u.Path)
			fmt.Println("port", u.Port())

			u, _ = url.Parse("127.0.0.1:9999/rest/v2")
			fmt.Println("host", u.Host)
			fmt.Println("path", u.Path)
			fmt.Println("port", u.Port())

			u, _ = url.Parse("http://abc/b")
			fmt.Println("host", u.Host)
			fmt.Println("path", u.Path)
			fmt.Println("port", u.Port())
			//schema, ip, port, err := urlutil.SplitURL("https://127.0.0.1")
			//So(err, ShouldBeNil)
			//So(schema, ShouldEqual, "https")
			//So(ip, ShouldEqual, "127.0.0.1")
			//So(port, ShouldEqual, "443")
			//So(urlutil.BuildURL(schema, ip, strconv.Itoa(port)), ShouldEqual, "https://127.0.0.1")
		})
	})
}

//func TestReplaceHost(t *testing.T) {
//	Convey("TestReplaceHost", t, func() {
//		t1 := "https://127.0.0.1"
//
//	})
//}
