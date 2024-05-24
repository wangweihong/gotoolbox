package fieldutil_test

import (
	"testing"

	"github.com/wangweihong/gotoolbox/pkg/fieldutil"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParseStructField(t *testing.T) {
	Convey("提取结构体中的字段信息", t, func() {
		So(fieldutil.ParseStructFields(nil), ShouldBeNil)
		So(fieldutil.ParseStructFields("str"), ShouldBeNil)
		So(fieldutil.ParseStructFields(0), ShouldBeNil)

		type Embedded struct {
			B string `json:"b"`
		}
		type Mystruct struct {
			Embedded
			A string `json:"a"`
			c string `json:"c"`
		}
		var myP *Mystruct

		pfs := fieldutil.ParseStructFields(myP)
		So(pfs, ShouldBeNil)

		myP = &Mystruct{
			Embedded: Embedded{B: "b"},
			A:        "123",
			c:        "o",
		}

		fs := fieldutil.ParseStructFields(myP)
		So(fs, ShouldNotBeNil)
		So(len(fs), ShouldEqual, 3)
	})
}

func TestParseStructFieldTags(t *testing.T) {
	Convey("提取结构体中的字段和标签信息", t, func() {
		So(fieldutil.ParseStructFieldTags(nil, "json"), ShouldBeNil)
		So(fieldutil.ParseStructFieldTags("str", "json"), ShouldBeNil)
		So(fieldutil.ParseStructFieldTags(0, "json"), ShouldBeNil)
		So(fieldutil.ParseStructFieldTags(0, ""), ShouldBeNil)

		type Embedded struct {
			B string `json:"b"`
		}
		type Mystruct struct {
			Embedded
			A string `json:"a"`
			c string `json:"c"`
		}
		var myP *Mystruct

		pfs := fieldutil.ParseStructFieldTags(myP, "json")
		So(pfs, ShouldBeNil)

		myP = &Mystruct{
			Embedded: Embedded{B: "b"},
			A:        "123",
			c:        "o",
		}

		fs := fieldutil.ParseStructFieldTags(myP, "json")
		So(fs, ShouldNotBeNil)
		So(len(fs), ShouldEqual, 2)
		So(len(fs.Tags()), ShouldEqual, 1)
		So(fs.Tags()[0], ShouldEqual, "a")
	})
}
