package stringutil_test

import (
	. "github.com/smartystreets/goconvey/convey"

	"github.com/wangweihong/gotoolbox/pkg/stringutil"

	"testing"
)

func TestBothEmptyOrNone(t *testing.T) {
	Convey("BothEmptyOrNone", t, func() {
		So(stringutil.BothEmptyOrNone("a", ""), ShouldBeFalse)
		So(stringutil.BothEmptyOrNone("a", "b"), ShouldBeTrue)
		So(stringutil.BothEmptyOrNone("", "b"), ShouldBeFalse)
	})
}

func TestHasAnyPrefix(t *testing.T) {
	Convey("HasAnyPrefix", t, func() {
		str := "tcp://192.168.134.132"
		So(stringutil.HasAnyPrefix(str, ""), ShouldBeFalse)
		So(stringutil.HasAnyPrefix("", ""), ShouldBeFalse)
		So(stringutil.HasAnyPrefix(str, "http", "https"), ShouldBeFalse)
		So(stringutil.HasAnyPrefix(str, "http", ""), ShouldBeFalse)
		So(stringutil.HasAnyPrefix(str, "tcp", "unix"), ShouldBeTrue)
	})
}

func TestPointerToString(t *testing.T) {
	Convey("ToString", t, func() {
		s := "a"
		var sp *string
		So(stringutil.PointerToString(sp), ShouldEqual, "")
		sp = &s
		So(stringutil.PointerToString(sp), ShouldEqual, "a")
	})
}

func TestAddIf(t *testing.T) {
	Convey("AddIf", t, func() {
		a := "str"
		So(stringutil.AddPrefixIfNotHas(a, "my"), ShouldEqual, "mystr")
		So(stringutil.AddSuffixIfNotHas(a, "my"), ShouldEqual, "strmy")

		b := "prefixmysuffix"
		So(stringutil.AddSuffixIfNotHas(b, "suffix"), ShouldEqual, "prefixmysuffix")
		So(stringutil.AddPrefixIfNotHas(b, "prefix"), ShouldEqual, "prefixmysuffix")
	})
}
