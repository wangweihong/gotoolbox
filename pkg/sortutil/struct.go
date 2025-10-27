package sortutil

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/wangweihong/gotoolbox/pkg/compareutil"

	"github.com/wangweihong/gotoolbox/pkg/structutil"

	"github.com/wangweihong/gotoolbox/pkg/fieldutil"

	"github.com/wangweihong/gotoolbox/pkg/sliceutil"
)

func FieldTagSort(si any, tag string, asc bool, defaultComparer, condition func(i, j int) bool) error {
	if !sliceutil.IsSliceOfStructs(si) {
		return fmt.Errorf("no slice or elem no struct")
	}

	sl := sliceutil.ToInterfaceSlice(si)
	if len(sl) == 0 {
		return nil
	}

	tagMap := fieldutil.GetFieldTagMapping(structutil.InitializeStruct(sl[0]))
	field, ok := tagMap[tag]
	if ok {
		fmt.Println("find tag:", tag)
		sort.SliceStable(sl, func(i, j int) bool {
			ret := 0

			vi, vj := reflect.ValueOf(sl[i]), reflect.ValueOf(sl[j])
			fi, fj, canCompare := compareutil.IsStructReflectValueCanCompare(vi, vj, field)
			if canCompare {
				ret = compareutil.ReflectCompare(fi, fj)
			}
			fmt.Printf("canCompare:%v, ret:%v\n", canCompare, ret)

			var result bool
			switch ret {
			case 1:
				result = true
			case -1:
				result = false
			default:
				result = defaultComparer(i, j)
			}

			if asc {
				return !result
			} else {
				return result
			}
		})
	} else {
		sort.SliceStable(sl, defaultComparer)
	}
	return nil
}
