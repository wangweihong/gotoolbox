package maputil_test

import (
	"strings"
	"testing"

	"github.com/wangweihong/gotoolbox/pkg/maputil"
	"github.com/wangweihong/gotoolbox/pkg/sets"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStringBoolMap_Init(t *testing.T) {
	Convey("TestStringBoolMap_Init", t, func() {
		var nilMap map[string]bool

		Convey("not nil", func() {
			So(nilMap, ShouldBeNil)

			maputil.StringBoolMap(nilMap).Init()
			So(nilMap, ShouldBeNil)

			nilMap = maputil.StringBoolMap(nilMap).Init()
			So(nilMap, ShouldNotBeNil)
		})
	})
}

func TestStringBoolMap_Set(t *testing.T) {
	Convey("TestStringBoolMap_Set", t, func() {
		var nilMap map[string]bool

		Convey("not nil", func() {
			So(nilMap, ShouldBeNil)

			maputil.StringBoolMap(nilMap).Set("1", true)
			So(nilMap, ShouldBeNil)

			nilMap = maputil.StringBoolMap(nilMap).Set("a", true).Set("c", true)
			So(nilMap, ShouldNotBeNil)
			So(len(nilMap), ShouldEqual, 2)
		})
	})
}

func TestStringBoolMap_DeepCopy(t *testing.T) {
	Convey("TestStringBoolMap_Set", t, func() {
		var nilMap map[string]bool

		Convey("not nil", func() {
			So(nilMap, ShouldBeNil)

			nilMap = maputil.StringBoolMap(nilMap).DeepCopy()
			So(nilMap, ShouldNotBeNil)

			nilMap = maputil.StringBoolMap(nilMap).Set("a", true)
			So(nilMap, ShouldNotBeNil)
			So(len(nilMap), ShouldEqual, 1)
			So(maputil.StringBoolMap(nilMap).Has("a"), ShouldBeTrue)
		})
	})
}

func TestStringBoolMap_Delete(t *testing.T) {
	Convey("TestStringBoolMap_Delete", t, func() {
		Convey("nil", func() {
			var nilMap map[string]bool
			maputil.StringBoolMap(nilMap).Delete("a")
		})
		Convey("not nil", func() {
			d := make(map[string]bool)
			d["a"] = true

			maputil.StringBoolMap(d).Delete("a")
			So(maputil.StringBoolMap(d).Has("a"), ShouldBeFalse)
		})
	})
}

func TestStringBoolMap_DeleteIfKey(t *testing.T) {
	Convey("TestStringBoolMap_DeleteIfKey", t, func() {
		condition := func(k string) bool {
			if strings.Contains(k, "b") {
				return true
			}
			return false
		}
		Convey("nil", func() {
			var nilMap map[string]bool
			maputil.StringBoolMap(nilMap).DeleteIfKey(condition)
		})
		Convey("not nil", func() {
			d := make(map[string]bool)
			d["ab"] = true
			d["bb"] = true
			d["cc"] = true
			maputil.StringBoolMap(d).DeleteIfKey(condition)

			So(maputil.StringBoolMap(d).Has("ab"), ShouldBeFalse)
			So(maputil.StringBoolMap(d).Has("bb"), ShouldBeFalse)
			So(maputil.StringBoolMap(d).Has("cc"), ShouldBeTrue)
		})
	})
}

func TestStringBoolMap_DeleteIfValue(t *testing.T) {
	Convey("TestStringBoolMap_DeleteIfValue", t, func() {
		condition := func(k bool) bool {
			return k
		}
		Convey("nil", func() {
			var nilMap map[string]bool
			maputil.StringBoolMap(nilMap).DeleteIfValue(condition)
		})
		Convey("not nil", func() {
			d := make(map[string]bool)
			d["ab"] = true
			d["bb"] = false
			d["cc"] = true
			maputil.StringBoolMap(d).DeleteIfValue(condition)
			So(maputil.StringBoolMap(d).Has("ab"), ShouldBeFalse)
			So(maputil.StringBoolMap(d).Has("bb"), ShouldBeTrue)
			So(maputil.StringBoolMap(d).Has("cc"), ShouldBeFalse)
		})
	})
}

func TestStringBoolMap_Get(t *testing.T) {
	Convey("TestStringBoolMap_Get", t, func() {
		Convey("nil", func() {
			var nilMap map[string]bool
			So(maputil.StringBoolMap(nilMap).Get("notexist"), ShouldBeFalse)
		})
		Convey("not nil", func() {
			d := make(map[string]bool)
			d["a"] = true

			So(maputil.StringBoolMap(d).Get("a"), ShouldBeTrue)
			So(maputil.StringBoolMap(d).Get("noexist"), ShouldBeFalse)
		})
	})
}

func TestStringBoolMap_Keys(t *testing.T) {
	Convey("TestStringBoolMap_Keys", t, func() {
		Convey("nil", func() {
			var nilMap map[string]bool
			keys := maputil.StringBoolMap(nilMap).Keys()

			So(len(keys), ShouldEqual, 0)
		})
		Convey("not nil", func() {
			d := make(map[string]bool)
			d["a"] = true
			d["b"] = true

			keys := maputil.StringBoolMap(d).Keys()
			So(len(keys), ShouldEqual, 2)
			So(sets.NewString(keys...).Equal(sets.NewString("a", "b")), ShouldBeTrue)
		})
	})
}

func TestStringBoolMap_ToSetString(t *testing.T) {
	Convey("TestStringBoolMap_ToSetString", t, func() {
		So(maputil.StringBoolMap(nil).ToSetString(), ShouldNotBeNil)

		d := make(map[string]bool)
		d["a"] = true
		d["b"] = true

		ss := maputil.StringBoolMap(d).ToSetString()
		So(ss.Len(), ShouldEqual, 2)
	})
}
