package fieldutil

import (
	"reflect"
	"strings"
)

const (
	JsonTag = "json"
)

type Fields []reflect.StructField

func ParseStructFields(s interface{}) Fields {
	if s == nil {
		return nil
	}

	if reflect.TypeOf(s).Kind() != reflect.Struct {
		return nil
	}

	typ := reflect.TypeOf(s)
	fts := make([]reflect.StructField, 0, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fts = append(fts, field)
	}
	return fts
}

func (fs Fields) ToMap() map[string]reflect.StructField {
	ftm := make(map[string]reflect.StructField, len(fs))
	for _, v := range fs {
		ftm[v.Name] = v
	}
	return ftm
}

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

func ParseStructFieldTags(s interface{}, tagName string) FieldTags {
	fs := ParseStructFields(s)
	if fs == nil {
		return nil
	}

	fts := make([]FieldTag, 0, len(fs))
	for _, v := range fs {
		if !v.IsExported() {
			continue
		}

		tag := v.Tag.Get(tagName)

		if tagName == "json" {
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
