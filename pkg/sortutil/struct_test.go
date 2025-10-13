package sortutil_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/wangweihong/gotoolbox/pkg/sortutil"
)

func TestFieldTagSorter_Sort(t *testing.T) {
	Convey("TestFieldTagSorter_Sort", t, func() {
		type Entry struct {
			Key  int    `json:"key"`
			Name string `json:"name"`
		}
		var entrys []Entry
		entrys = append(entrys, Entry{
			Key:  3,
			Name: "c",
		})
		entrys = append(entrys, Entry{
			Key:  4,
			Name: "b",
		})
		entrys = append(entrys, Entry{
			Key:  5,
			Name: "a",
		})
		entrys = append(entrys, Entry{
			Key:  6,
			Name: "a",
		})
		err := sortutil.FieldTagSort(entrys, "name", true, func(i, j int) bool {
			return entrys[i].Key < entrys[j].Key
		}, nil)
		So(err, ShouldBeNil)
		So(entrys, ShouldResemble, []Entry{{5, "a"}, {6, "a"}, {4, "b"}, {3, "c"}})
	})
}
