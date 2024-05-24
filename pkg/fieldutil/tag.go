package fieldutil

import (
	"reflect"
	"strings"
)

const (
	JsonTag = "json"
)

type Fields []reflect.StructField

// ParseStructFields 从结构体或者结构体指针(指针非空值)中获取其字段信息
func ParseStructFields(s interface{}) Fields {
	if s == nil {
		return nil
	}

	typ := reflect.TypeOf(s)
	vyp := reflect.ValueOf(s)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		vyp = vyp.Elem()
		if !vyp.IsValid() {
			return nil
		}
	}

	if typ.Kind() != reflect.Struct {
		return nil
	}

	if typ.NumField() != 0 {
		fts := make([]reflect.StructField, 0, typ.NumField())
		// 这里不处理匿名字段，应由字段本身去根据类型去处理
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			fts = append(fts, field)
		}
		return fts
	}

	return nil
}

func (fs Fields) ToMap() map[string]reflect.StructField {
	ftm := make(map[string]reflect.StructField, len(fs))
	for _, v := range fs {
		ftm[v.Name] = v
	}
	return ftm
}

// 字段和标签信息
// 注意字段可能是embedded struct, 通过Field.Anonymous来判定
type FieldTag struct {
	Field reflect.StructField
	Tag   string
}

type FieldTags []FieldTag

func (fts FieldTags) ToMap() map[string]FieldTag {
	ftm := make(map[string]FieldTag, len(fts))
	for _, v := range fts {
		ftm[v.Field.Name] = v
	}
	return ftm
}

func (fts FieldTags) ToTagMap() map[string]FieldTag {
	ftm := make(map[string]FieldTag, len(fts))
	for _, v := range fts {
		if v.Tag != "" {
			ftm[v.Tag] = v
		}
	}
	return ftm
}

// ParseStructFields 从结构体或者结构体指针(指针非空值)中获取其字段以及标签信息
func ParseStructFieldTags(s interface{}, tagName string) FieldTags {
	fs := ParseStructFields(s)
	if fs == nil {
		return nil
	}

	fts := make([]FieldTag, 0, len(fs))
	for _, v := range fs {
		// 忽略非导出字段
		if !v.IsExported() {
			continue
		}
		tag := v.Tag.Get(tagName)
		if tagName == JsonTag {
			tag = strings.TrimSuffix(tag, ",omitempty")
			if tag == "-" {
				continue
			}
		}

		fts = append(fts, FieldTag{
			Field: v,
			Tag:   tag,
		})
	}
	return fts
}

func (fts FieldTags) Tags() []string {
	tags := make([]string, 0, len(fts))
	for _, v := range fts {
		if v.Tag != "" {
			tags = append(tags, v.Tag)
		}
	}
	return tags
}
