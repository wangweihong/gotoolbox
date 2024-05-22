package maputil_test

import (
	"strings"
	"testing"

	"github.com/wangweihong/gotoolbox/src/maputil"
	"github.com/wangweihong/gotoolbox/src/sets"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStringStringMap_Init(t *testing.T) {
	Convey("TestStringStringMap_Init", t, func() {
		var nilMap map[string]string

		Convey("not nil", func() {
			So(nilMap, ShouldBeNil)

			maputil.StringStringMap(nilMap).Init()
			So(nilMap, ShouldBeNil)

			nilMap = maputil.StringStringMap(nilMap).Init()
			So(nilMap, ShouldNotBeNil)
		})
	})
}

func TestStringStringMap_Set(t *testing.T) {
	Convey("TestStringStringMap_Set", t, func() {
		var nilMap map[string]string

		Convey("not nil", func() {
			So(nilMap, ShouldBeNil)

			maputil.StringStringMap(nilMap).Set("1", "2")
			So(nilMap, ShouldBeNil)

			nilMap = maputil.StringStringMap(nilMap).Set("a", "b").Set("c", "D")
			So(nilMap, ShouldNotBeNil)
			So(len(nilMap), ShouldEqual, 2)
		})
	})
}

func TestStringStringMap_DeepCopy(t *testing.T) {
	Convey("TestStringStringMap_Set", t, func() {
		var nilMap map[string]string

		Convey("not nil", func() {
			So(nilMap, ShouldBeNil)

			nilMap = maputil.StringStringMap(nilMap).DeepCopy()
			So(nilMap, ShouldNotBeNil)

			nilMap = maputil.StringStringMap(nilMap).Set("a", "b")
			So(nilMap, ShouldNotBeNil)
			So(len(nilMap), ShouldEqual, 1)
			So(maputil.StringStringMap(nilMap).Has("a"), ShouldBeTrue)
		})
	})
}

func TestStringStringMap_Delete(t *testing.T) {
	Convey("TestStringStringMap_Delete", t, func() {
		Convey("nil", func() {
			var nilMap map[string]string
			maputil.StringStringMap(nilMap).Delete("a")
		})
		Convey("not nil", func() {
			d := make(map[string]string)
			d["a"] = "b"

			maputil.StringStringMap(d).Delete("a")
			So(maputil.StringStringMap(d).Has("a"), ShouldBeFalse)
		})
	})
}

func TestStringStringMap_DeleteIfKey(t *testing.T) {
	Convey("TestStringStringMap_DeleteIfKey", t, func() {
		condition := func(k string) bool {
			if strings.Contains(k, "b") {
				return true
			}
			return false
		}
		Convey("nil", func() {
			var nilMap map[string]string
			maputil.StringStringMap(nilMap).DeleteIfKey(condition)
		})
		Convey("not nil", func() {
			d := make(map[string]string)
			d["ab"] = "ab"
			d["bb"] = "bb"
			d["cc"] = "cc"
			maputil.StringStringMap(d).DeleteIfKey(condition)

			So(maputil.StringStringMap(d).Has("ab"), ShouldBeFalse)
			So(maputil.StringStringMap(d).Has("bb"), ShouldBeFalse)
			So(maputil.StringStringMap(d).Has("cc"), ShouldBeTrue)
		})
	})
}

func TestStringStringMap_DeleteIfValue(t *testing.T) {
	Convey("TestStringStringMap_DeleteIfValue", t, func() {
		condition := func(k string) bool {
			if strings.Contains(k, "b") {
				return true
			}
			return false
		}
		Convey("nil", func() {
			var nilMap map[string]string
			maputil.StringStringMap(nilMap).DeleteIfValue(condition)
		})
		Convey("not nil", func() {
			d := make(map[string]string)
			d["ab"] = "ab"
			d["bb"] = "bb"
			d["cc"] = "cc"
			maputil.StringStringMap(d).DeleteIfValue(condition)
			So(maputil.StringStringMap(d).Has("ab"), ShouldBeFalse)
			So(maputil.StringStringMap(d).Has("bb"), ShouldBeFalse)
			So(maputil.StringStringMap(d).Has("cc"), ShouldBeTrue)
		})
	})
}

func TestStringStringMap_Get(t *testing.T) {
	Convey("TestStringStringMap_Get", t, func() {
		Convey("nil", func() {
			var nilMap map[string]string

			So(maputil.StringStringMap(nilMap).Get("notexist"), ShouldEqual, "")
		})
		Convey("not nil", func() {
			d := make(map[string]string)
			d["a"] = "b"

			So(maputil.StringStringMap(d).Get("a"), ShouldEqual, "b")
			So(maputil.StringStringMap(d).Get("notexist"), ShouldEqual, "")
		})
	})
}

func TestStringStringMap_Keys(t *testing.T) {
	Convey("TestStringStringMap_Keys", t, func() {
		Convey("nil", func() {
			var nilMap map[string]string
			keys := maputil.StringStringMap(nilMap).Keys()

			So(len(keys), ShouldEqual, 0)
		})
		Convey("not nil", func() {
			d := make(map[string]string)
			d["a"] = "1"
			d["b"] = "2"

			keys := maputil.StringStringMap(d).Keys()
			So(len(keys), ShouldEqual, 2)
			So(sets.NewString(keys...).Equal(sets.NewString("a", "b")), ShouldBeTrue)
		})
	})
}
