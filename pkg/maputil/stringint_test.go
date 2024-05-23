package maputil_test

import (
	"strings"
	"testing"

	"github.com/wangweihong/gotoolbox/pkg/maputil"
	"github.com/wangweihong/gotoolbox/pkg/sets"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStringIntMap_Init(t *testing.T) {
	Convey("TestStringIntMap_Init", t, func() {
		var nilMap map[string]int

		Convey("not nil", func() {
			So(nilMap, ShouldBeNil)

			maputil.StringIntMap(nilMap).Init()
			So(nilMap, ShouldBeNil)

			nilMap = maputil.StringIntMap(nilMap).Init()
			So(nilMap, ShouldNotBeNil)
		})
	})
}

func TestStringIntMap_Set(t *testing.T) {
	Convey("TestStringIntMap_Set", t, func() {
		var nilMap map[string]int

		Convey("not nil", func() {
			So(nilMap, ShouldBeNil)

			maputil.StringIntMap(nilMap).Set("1", 2)
			So(nilMap, ShouldBeNil)

			nilMap = maputil.StringIntMap(nilMap).Set("a", 3).Set("c", 3)
			So(nilMap, ShouldNotBeNil)
			So(len(nilMap), ShouldEqual, 2)
		})
	})
}

func TestStringIntMap_DeepCopy(t *testing.T) {
	Convey("TestStringIntMap_Set", t, func() {
		var nilMap map[string]int

		Convey("not nil", func() {
			So(nilMap, ShouldBeNil)

			nilMap = maputil.StringIntMap(nilMap).DeepCopy()
			So(nilMap, ShouldNotBeNil)

			nilMap = maputil.StringIntMap(nilMap).Set("a", 4)
			So(nilMap, ShouldNotBeNil)
			So(len(nilMap), ShouldEqual, 1)
			So(maputil.StringIntMap(nilMap).Has("a"), ShouldBeTrue)
		})
	})
}

func TestStringIntMap_Delete(t *testing.T) {
	Convey("TestStringIntMap_Delete", t, func() {
		Convey("nil", func() {
			var nilMap map[string]int
			maputil.StringIntMap(nilMap).Delete("a")
		})
		Convey("not nil", func() {
			d := make(map[string]int)
			d["a"] = 3

			maputil.StringIntMap(d).Delete("a")
			So(maputil.StringIntMap(d).Has("a"), ShouldBeFalse)
		})
	})
}

func TestStringIntMap_DeleteIfKey(t *testing.T) {
	Convey("TestStringIntMap_DeleteIfKey", t, func() {
		condition := func(k string) bool {
			if strings.Contains(k, "b") {
				return true
			}
			return false
		}
		Convey("nil", func() {
			var nilMap map[string]int
			maputil.StringIntMap(nilMap).DeleteIfKey(condition)
		})
		Convey("not nil", func() {
			d := make(map[string]int)
			d["ab"] = 10
			d["bb"] = 20
			d["cc"] = 31
			maputil.StringIntMap(d).DeleteIfKey(condition)

			So(maputil.StringIntMap(d).Has("ab"), ShouldBeFalse)
			So(maputil.StringIntMap(d).Has("bb"), ShouldBeFalse)
			So(maputil.StringIntMap(d).Has("cc"), ShouldBeTrue)
		})
	})
}

func TestStringIntMap_DeleteIfValue(t *testing.T) {
	Convey("TestStringIntMap_DeleteIfValue", t, func() {
		condition := func(k int) bool {
			if k%10 == 0 {
				return true
			}
			return false
		}
		Convey("nil", func() {
			var nilMap map[string]int
			maputil.StringIntMap(nilMap).DeleteIfValue(condition)
		})
		Convey("not nil", func() {
			d := make(map[string]int)
			d["ab"] = 10
			d["bb"] = 20
			d["cc"] = 31
			maputil.StringIntMap(d).DeleteIfValue(condition)
			So(maputil.StringIntMap(d).Has("ab"), ShouldBeFalse)
			So(maputil.StringIntMap(d).Has("bb"), ShouldBeFalse)
			So(maputil.StringIntMap(d).Has("cc"), ShouldBeTrue)
		})
	})
}

func TestStringIntMap_Get(t *testing.T) {
	Convey("TestStringIntMap_Get", t, func() {
		Convey("nil", func() {
			var nilMap map[string]int

			So(maputil.StringIntMap(nilMap).Get("notexist"), ShouldEqual, 0)
		})
		Convey("not nil", func() {
			d := make(map[string]int)
			d["a"] = 2

			So(maputil.StringIntMap(d).Get("a"), ShouldEqual, 2)
			So(maputil.StringIntMap(d).Get("notexist"), ShouldEqual, 0)
		})
	})
}

func TestStringIntMap_Keys(t *testing.T) {
	Convey("TestStringIntMap_Keys", t, func() {
		Convey("nil", func() {
			var nilMap map[string]int
			keys := maputil.StringIntMap(nilMap).Keys()

			So(len(keys), ShouldEqual, 0)
		})
		Convey("not nil", func() {
			d := make(map[string]int)
			d["a"] = 1
			d["b"] = 2

			keys := maputil.StringIntMap(d).Keys()
			So(len(keys), ShouldEqual, 2)
			So(sets.NewString(keys...).Equal(sets.NewString("a", "b")), ShouldBeTrue)
		})
	})
}
