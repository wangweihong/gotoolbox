package sets_test

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/wangweihong/gotoolbox/src/sets"
)

func TestString_ToSetString(t *testing.T) {
	Convey("TestStringBoolMap_ToSetString", t, func() {
		d := sets.NewString()
		d.Insert("stringa", "sstringb")

		d2 := d.FindMatch(func(s string, s2 string) bool {
			if strings.Contains(s, s2) {
				return true
			} else {
				return false
			}
		}, "string")

		So(d2.Len(), ShouldEqual, d.Len())
		d3 := d.FindMatch(func(s string, s2 string) bool {
			if strings.HasPrefix(s, s2) {
				return true
			} else {
				return false
			}
		}, "string")
		So(d3.Len(), ShouldEqual, 1)
		So(d3.List(), ShouldResemble, []string{"stringa"})

	})
}
