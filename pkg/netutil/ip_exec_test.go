package netutil_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/wangweihong/gotoolbox/pkg/netutil"
)

func TestPing(t *testing.T) {
	Convey("ping", t, func() {
		So(netutil.Ping("127.0.0.1"), ShouldBeTrue)
	})
}
