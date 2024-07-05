package fieldutil

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"

	"github.com/wangweihong/gotoolbox/pkg/errors"
)

// 当apiObject存在一个字段包含tagName:tagValues、且类型和internalObject一致,则设置该字段值为internalObject
func SetWhenTagValueMatch(apiObject interface{}, internalObject interface{}, tagName, tagValue string) error {
	if tagName == "" || tagValue == "" {
		return errors.New("tagName and tagValue is empty")
	}

	rv := reflect.ValueOf(apiObject)
	rt := reflect.TypeOf(apiObject)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
		rt = rt.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return errors.New("apiObject not struct")
	}

	if !rv.IsValid() {
		return errors.New("apiObject not vaild")
	}

	for i := 0; i < rt.NumField(); i++ {
		if rt.Field(i).Tag.Get(tagName) == tagValue {
			// ensure apiObject has tagName/tagValue

			// if internalObject is nil, do nothing.
			if internalObject == nil {
				return nil
			}

			if rt.Field(i).Type != reflect.TypeOf(internalObject) {
				return errors.Errorf("apiObject field with %v:\"%v\" type (%v) not equal to internalObject （%v）",
					tagName, tagValue, rt.Field(i).Type, reflect.TypeOf(internalObject))
			}

			if rv.Field(i).CanSet() {
				rv.Field(i).Set(reflect.ValueOf(internalObject))
			}
			return nil
		}
	}
	return errors.New("missing field with tag json:data")
}

func CheckIfStructFieldMatch(apiObject interface{}, tagName, tagValuePoint string, compareValue interface{}) error {
	if tagValuePoint == "" || strings.TrimSpace(tagValuePoint) == "." {
		return errors.New("invalid tagValuePoint")
	}

	if apiObject == nil {
		return errors.New("apiObject is nil")
	}

	valuePoints := strings.Split(tagValuePoint, ".")

	rv := reflect.ValueOf(apiObject)
	rt := reflect.TypeOf(apiObject)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
		rt = rt.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return errors.New("apiObject not pointer")
	}

	if !rv.IsValid() {
		return errors.New("apiObject not vaild")
	}

	for i := 0; i < rt.NumField(); i++ {
		if rt.Field(i).Tag.Get(tagName) == valuePoints[0] {
			// 查看valuepoint是否还有下一级
			if len(valuePoints) > 1 {
				if rv.Field(i).Kind() != reflect.Struct {
					return errors.New("field kind should be struct when valuepoint has multple elem")
				}

				return CheckIfStructFieldMatch(
					rv.Field(i).Interface(),
					tagName,
					strings.Join(valuePoints[1:], "."),
					compareValue,
				)
			}

			if rt.Field(i).Type != reflect.TypeOf(compareValue) {
				return errors.Errorf(
					"value type not match fieldValue, value type (%v),field type(%v)",
					rt.Field(i).Type,
					reflect.TypeOf(compareValue),
				)
			}

			// will this work?
			// https://github.com/golang/go/issues/9504
			if rv.Field(i).Interface() != reflect.ValueOf(compareValue).Interface() {
				return errors.New("value not match")
			}

			return nil
		}
	}
	return errors.Errorf("cannot find field with tag %v:%v in object", tagName, valuePoints[0])
}

// 检测json Marshal bytes字节流中是否指定的字段如,而且值为compareValue。 如 china.shenzhen.baoan.cars[0].numberprefix是否等于"粤B"
func CheckIfBytesStructFieldMatch(apiObjectBytes []byte, tagValuePoint string, compareValue interface{}) error {
	if tagValuePoint == "" || strings.TrimSpace(tagValuePoint) == "." {
		return errors.New("invalid tagValuePoint")
	}

	if apiObjectBytes == nil {
		return errors.New("apiObject is nil")
	}

	valuePoints := strings.Split(tagValuePoint, ".")
	mapKey := valuePoints[0]
	// 如果是 data.slice[0], 则提取出slice索引，后续提取该元素进行比对
	var sliceIndex int = -1
	if strings.HasSuffix(mapKey, "]") {
		k := strings.Index(mapKey, "[")
		if k != 0 {
			var err error
			sliceIndex, err = strconv.Atoi(mapKey[k+1 : len(mapKey)-1])
			if err != nil {
				return err
			}
			//移除字段名中索引部分，才能从map中找到对应的值
			mapKey = mapKey[:k]
		}
	}

	//json unmarhsal默认会将数值转换成float64类型来存储。
	//调用UserNumber会使用json.Number类型保存
	//后续通过转换jsonNumber成Float或者Int
	// 不要直接用json.Unmarshal.会导致数字转换成浮点类型
	d := json.NewDecoder(bytes.NewReader(apiObjectBytes))
	d.UseNumber()
	obj := make(map[string]interface{})
	if err := d.Decode(&obj); err != nil {
		return err
	}

	fieldValue, ok := obj[mapKey]
	if !ok {
		return errors.Errorf("field not exist:%v", mapKey)
	}
	fvt := reflect.TypeOf(fieldValue)
	fvv := reflect.ValueOf(fieldValue)

	// 只会有map/slice/string/int/float/bool类型!
	switch fvt.Kind() {
	//如果是数组,这需要处理两种情况1.比较整个数组，2.比较数组中某个元素
	case reflect.Slice:
		//某个元素
		if sliceIndex != -1 {
			//设置为指定元素值
			if sliceIndex >= fvv.Len() {
				return errors.Errorf("invalid slice index [%v], out of slice range", sliceIndex)
			}
			indexValue := fvv.Index(sliceIndex).Interface()
			fieldValue = indexValue
			// 还需要查找下一层
			if len(valuePoints) != 1 {
				//没有遍历完，继续遍历
				b, err := json.Marshal(fieldValue)
				if err != nil {
					return err
				}

				return CheckIfBytesStructFieldMatch(b, strings.Join(valuePoints[1:], "."), compareValue)
			}

			if reflect.TypeOf(fieldValue) != reflect.TypeOf(compareValue) {
				return errors.Errorf("type not match, source(type:%v), target(type:%v)",
					reflect.TypeOf(fieldValue), reflect.TypeOf(compareValue))
			}

			if fieldValue != compareValue {
				return errors.Errorf("value not match, source(value:%v), target(value:%v)",
					fieldValue, compareValue)
			}
			return nil
		} else {
			// 数组, 如果是数组就不能再走下一层
			if len(valuePoints) != 1 {
				return errors.New("data.slice[0].key or data.slice is good. but data.slice.key is not")
			}

			cv := reflect.ValueOf(compareValue)
			if cv.Kind() != reflect.Slice {
				return errors.Errorf("type not match, source(type:%v), target(type:%v)",
					reflect.TypeOf(fieldValue), reflect.TypeOf(compareValue))
			}

			if cv.Len() != fvv.Len() {
				return errors.New("compare field slice len not match")
			}

			for i := 0; i < cv.Len(); i++ {
				if cv.Index(i).Interface() != fvv.Index(i).Interface() {
					return errors.Errorf("value not match, source(value:%v), target(value:%v)",
						cv.Index(i), fvv.Index(i))
				}
			}
			return nil
		}

	case reflect.Map:
		// 还需要查找下一层
		if len(valuePoints) != 1 {
			//没有遍历完，继续遍历
			b, err := json.Marshal(fieldValue)
			if err != nil {
				return err
			}

			return CheckIfBytesStructFieldMatch(b, strings.Join(valuePoints[1:], "."), compareValue)
		}

		if reflect.TypeOf(fieldValue) != reflect.TypeOf(compareValue) {
			return errors.Errorf("type not match, source(type:%v), target(type:%v)",
				reflect.TypeOf(fieldValue), reflect.TypeOf(compareValue))
		}

		if fieldValue != compareValue {
			return errors.Errorf("value not match, source(value:%v), target(value:%v)",
				fieldValue, compareValue)
		}

		return nil

	default:
		// 除了map和slice其他类型不会有下一层
		if len(valuePoints) != 1 {
			return errors.New("non-slice or non-map should has no child point")
		}

		if fvt == reflect.TypeOf(json.Number("")) {
			switch reflect.TypeOf(compareValue).Kind() {
			case reflect.Int, reflect.Int16, reflect.Int8, reflect.Int32, reflect.Int64:
				jn := fieldValue.(json.Number)
				jnInt64, err := jn.Int64()
				if err != nil {
					return err
				}
				if jnInt64 != reflect.ValueOf(compareValue).Int() {
					return errors.Errorf("value not match, source(value:%v), target(value:%v)",
						fieldValue, compareValue)
				}

			case reflect.Float64, reflect.Float32:
				jn := fieldValue.(json.Number)
				jnFloat, err := jn.Float64()
				if err != nil {
					return err
				}
				if jnFloat != reflect.ValueOf(compareValue).Float() {
					return errors.Errorf("value not match, source(value:%v), target(value:%v)",
						fieldValue, compareValue)
				}
			default:
				if fvt != reflect.TypeOf(compareValue) {
					return errors.Errorf("type not match when type is json number, source(type:%v), target(type:%v)",
						"json.Number", reflect.TypeOf(compareValue))
				}
			}
		} else {
			if fvt != reflect.TypeOf(compareValue) {
				return errors.Errorf("type not match, source(type:%v), target(type:%v)",
					fvt, reflect.TypeOf(compareValue))
			}

			if fieldValue != compareValue {
				return errors.Errorf("value not match, source(value:%v), target(value:%v)",
					fieldValue, compareValue)
			}
		}
	}
	return nil
}

func GetStructFieldValue(object interface{}, tagValuePoint string) (interface{}, error) {
	if object == nil {
		return nil, errors.New("object is nil")
	}

	b, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}
	return GetBytesStructField(b, tagValuePoint)
}

// 提取json Marshal bytes字节流中是否指定的字段值 如 china.shenzhen.baoan.cars[0].numberprefix
func GetBytesStructField(apiObjectBytes []byte, tagValuePoint string) (interface{}, error) {
	if tagValuePoint == "" || strings.TrimSpace(tagValuePoint) == "." {
		return nil, errors.New("invalid tagValuePoint")
	}

	if apiObjectBytes == nil {
		return nil, errors.New("apiObject is nil")
	}

	valuePoints := strings.Split(tagValuePoint, ".")
	mapKey := valuePoints[0]
	// 如果是 data.slice[0], 则提取出slice索引，后续提取该元素进行比对
	var sliceIndex int = -1
	if strings.HasSuffix(mapKey, "]") {
		k := strings.Index(mapKey, "[")
		if k != 0 {
			var err error
			sliceIndex, err = strconv.Atoi(mapKey[k+1 : len(mapKey)-1])
			if err != nil {
				return nil, err
			}
			//移除字段名中索引部分，才能从map中找到对应的值
			mapKey = mapKey[:k]
		}
	}

	//json unmarhsal默认会将数值转换成float64类型来存储。
	//调用UserNumber会使用json.Number类型保存
	//后续通过转换jsonNumber成Float或者Int
	// 不要直接用json.Unmarshal.会导致数字转换成浮点类型
	d := json.NewDecoder(bytes.NewReader(apiObjectBytes))
	d.UseNumber()
	obj := make(map[string]interface{})
	if err := d.Decode(&obj); err != nil {
		return nil, err
	}

	fieldValue, ok := obj[mapKey]
	if !ok {
		return nil, errors.Errorf("field not exist:%v", mapKey)
	}
	fvt := reflect.TypeOf(fieldValue)
	fvv := reflect.ValueOf(fieldValue)

	// 只会有map/slice/string/int/float/bool类型!
	switch fvt.Kind() {
	//如果是数组,这需要处理两种情况1.比较整个数组，2.比较数组中某个元素
	case reflect.Slice:
		//某个元素
		if sliceIndex != -1 {
			//设置为指定元素值
			if sliceIndex >= fvv.Len() {
				return nil, errors.Errorf("invalid slice index [%v], out of slice range", sliceIndex)
			}
			indexValue := fvv.Index(sliceIndex).Interface()
			fieldValue = indexValue
			// 还需要查找下一层
			if len(valuePoints) != 1 {
				//没有遍历完，继续遍历
				b, err := json.Marshal(fieldValue)
				if err != nil {
					return nil, errors.WithStack(err)
				}

				return GetBytesStructField(b, strings.Join(valuePoints[1:], "."))
			}

			return fieldValue, nil
		} else {
			// 数组, 如果是数组就不能再走下一层
			if len(valuePoints) != 1 {
				return nil, errors.Errorf("data.slice[0].key or data.slice is good. but data.slice.key is not")
			}

			return fvv.Interface(), nil
		}

	case reflect.Map:
		// 还需要查找下一层
		if len(valuePoints) != 1 {
			//没有遍历完，继续遍历
			b, err := json.Marshal(fieldValue)
			if err != nil {
				return nil, err
			}

			return GetBytesStructField(b, strings.Join(valuePoints[1:], "."))
		}

		return fieldValue, nil

	default:
		// 除了map和slice其他类型不会有下一层
		if len(valuePoints) != 1 {
			return nil, errors.New("non-slice or non-map should has no child point")
		}
	}
	return fvv.Interface(), nil
}

func SetWhenFieldValueTypeMatch(apiObject interface{}, fieldName string, fieldValues interface{}) error {
	if fieldName == "" {
		return errors.New("fieldName not set")
	}

	rv := reflect.ValueOf(apiObject)
	rt := reflect.TypeOf(apiObject)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
		rt = rt.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return errors.Errorf("apiObject not struct,%v", rv.Kind())
	}

	if !rv.IsValid() {
		return errors.New("apiObject not vaild")
	}

	sft, exist := rt.FieldByName(fieldName)
	if !exist {
		return errors.Errorf("apiObject doesn't has field named %v", fieldName)
	}
	fv := rv.FieldByName(fieldName)
	ft := sft.Type
	if !fv.CanSet() {
		return errors.Errorf("field %v cannot set", fieldName)
	}

	if ft != reflect.TypeOf(fieldValues) {
		return errors.Errorf(
			"type not match, field %v type is %v while passing value type is %v",
			fieldName,
			ft,
			reflect.TypeOf(fieldValues),
		)
	}

	fv.Set(reflect.ValueOf(fieldValues))
	return nil
}

func SetFieldZeroValueIfCondition(apiObject interface{}, condition func(string) bool) {
	setZeroValueIfCondition(apiObject, condition)
}

func setZeroValueIfCondition(v interface{}, condition func(string) bool) {
	if v == nil {
		return
	}
	vl := reflect.ValueOf(v)
	if vl.Kind() == reflect.Ptr {
		vl = reflect.ValueOf(v).Elem()
	}

	setZeroValue(vl, condition)
}

func setZeroValue(val reflect.Value, condition func(string) bool) {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := val.Type().Field(i).Name
		if condition(fieldName) {
			if field.CanSet() {
				field.Set(reflect.Zero(field.Type()))
			}
		} else {
			switch field.Kind() {
			case reflect.Ptr:
				if field.IsNil() {
					continue
				}
				if field.Elem().Kind() == reflect.Struct {
					setZeroValue(field, condition)
				}
			case reflect.Array, reflect.Slice:
				for j := 0; j < field.Len(); j++ {
					elem := field.Index(j)
					if elem.Kind() == reflect.Ptr && elem.Elem().Kind() == reflect.Struct {
						setZeroValue(elem, condition)
					} else if elem.Kind() == reflect.Struct {
						setZeroValue(elem, condition)
					}
				}
			case reflect.Map:
				for _, key := range field.MapKeys() {
					elem := field.MapIndex(key)
					if elem.Kind() == reflect.Ptr && elem.Elem().Kind() == reflect.Struct {
						setZeroValue(elem, condition)
					} else if elem.Kind() == reflect.Struct {
						//FIXME: 如何支持map[string]struct?
					}
				}
			}
		}
	}
}
