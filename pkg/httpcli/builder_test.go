package httpcli_test

import (
	"testing"

	"github.com/wangweihong/gotoolbox/pkg/httpcli"
	"github.com/wangweihong/gotoolbox/pkg/maputil"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHttpRequestBuilder_AddQueryParamByObject(t *testing.T) {
	Convey("AddQueryParamByObject", t, func() {
		type objType struct {
			String  string
			Num     int
			Flag    bool
			Pointer *string
			WithTag string `json:"with_tag"`
			OmitTag string `json:"omit_tag,omitempty"`
			exx     string
		}
		ot := objType{
			String:  "name",
			Num:     123,
			Flag:    false,
			Pointer: nil,
			WithTag: "tag1",
			OmitTag: "",
		}
		//expect := make(map[string]interface{})
		expect := make(map[string]interface{})
		expect = map[string]interface{}{
			"String":   "name",
			"Num":      123,
			"Flag":     false,
			"with_tag": "tag1",
		}
		params := httpcli.NewHttpRequestBuilder().AddQueryParamByObject(ot).Build().GetQueryParams()

		So(maputil.StringInterfaceMap(params).Equal(expect), ShouldBeTrue)
	})
}

func TestHttpRequestBuilder_AddQueryParamByObjectEmbedded(t *testing.T) {
	Convey("AddQueryParamByObject", t, func() {
		type PagingParam struct {
			PageNum  int `form:"page_num" json:"page_num"`
			PageSize int `form:"page_size" json:"page_size"`
		}

		type objType struct {
			PagingParam
			String string
		}

		Convey("匿名结构字段不设置", func() {
			ot := objType{
				String: "name",
			}
			//expect := make(map[string]interface{})
			expect := make(map[string]interface{})
			expect = map[string]interface{}{
				"String":    "name",
				"page_num":  0,
				"page_size": 0,
			}
			So(maputil.StringInterfaceMap(httpcli.NewHttpRequestBuilder().AddQueryParamByObject(ot).Build().GetQueryParams()).Equal(expect), ShouldBeTrue)
		})

		Convey("匿名结构字段设置", func() {
			ot := objType{
				String: "name",
				PagingParam: PagingParam{
					PageNum:  1,
					PageSize: 3,
				},
			}
			//expect := make(map[string]interface{})
			expect := make(map[string]interface{})
			expect = map[string]interface{}{
				"String":    "name",
				"page_num":  1,
				"page_size": 3,
			}
			So(maputil.StringInterfaceMap(httpcli.NewHttpRequestBuilder().AddQueryParamByObject(ot).Build().GetQueryParams()).Equal(expect), ShouldBeTrue)
		})
	})
}
