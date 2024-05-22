package httpcli_test

import (
	"testing"

	"github.com/wangweihong/gotoolbox/src/httpcli"
	"github.com/wangweihong/gotoolbox/src/maputil"

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
