package compareutil

import (
	"fmt"
	"reflect"
)

func StructReflectValueCompare(v1, v2 reflect.Value, fieldName string) int {
	if v1.Type().Kind() != reflect.Struct || v2.Type().Kind() != reflect.Struct {
		return 0
	}

	fi := v1.FieldByName(fieldName)
	fj := v2.FieldByName(fieldName)

	if !fi.IsValid() || !fj.IsValid() {
		return 0
	}

	if fi.Interface() == fj.Interface() {
		return 0
	}

	switch fi.Kind() {
	case reflect.String:
		if fi.String() < fj.String() {
			return -1
		}

		if fi.String() > fj.String() {
			return 1
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if fi.Int() < fj.Int() {
			return -1
		}

		if fi.Int() > fj.Int() {
			return 1
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if fi.Uint() < fj.Uint() {
			return -1
		}

		if fi.Uint() > fj.Uint() {
			return 1
		}
	case reflect.Float32, reflect.Float64:
		if fi.Float() < fj.Float() {
			return -1
		}

		if fi.Float() > fj.Float() {
			return 1
		}
	default:
		// 默认使用字符串比较
		si := fmt.Sprintf("%v", fi.Interface())
		sj := fmt.Sprintf("%v", fj.Interface())

		if si < sj {
			return -1
		}

		if si > sj {
			return 1
		}
	}
	return 0
}
