package maputil

import "github.com/wangweihong/gotoolbox/sets"

type StringBoolMap map[string]bool

func (m StringBoolMap) DeepCopy() StringBoolMap {
	o := make(map[string]bool, len(m))
	for k, v := range m {
		o[k] = v
	}
	return o
}

func (m StringBoolMap) Init() StringBoolMap {
	if m == nil {
		return make(map[string]bool)
	}
	return m
}

func (m StringBoolMap) Delete(key string) {
	if m == nil {
		return
	}
	delete(m, key)
}

func (m StringBoolMap) DeleteIfKey(condition func(string) bool) {
	if m == nil {
		return
	}

	for k := range m {
		if condition(k) {
			delete(m, k)
		}
	}
}

func (m StringBoolMap) DeleteIfValue(condition func(bool) bool) {
	if m == nil {
		return
	}

	for k, v := range m {
		if condition(v) {
			delete(m, k)
		}
	}
}

func (m StringBoolMap) Has(key string) bool {
	if m != nil {
		if _, exist := m[key]; exist {
			return true
		}
	}
	return false
}

func (m StringBoolMap) HasAny(key string) bool {
	if m != nil {
		if _, exist := m[key]; exist {
			return true
		}
	}
	return false
}

func (m StringBoolMap) Set(key string, value bool) StringBoolMap {
	if m == nil {
		o := make(map[string]bool)
		o[key] = value
		return o
	}
	m[key] = value
	return m
}

func (m StringBoolMap) Get(key string) bool {
	if m == nil {
		return false
	}
	v, _ := m[key]
	return v
}

func (m StringBoolMap) Keys() []string {
	if m == nil {
		return []string{}
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func (m StringBoolMap) ToSetString() sets.String {
	ss := sets.NewString()

	if m == nil {
		return ss
	}
	for k := range m {
		ss.Insert(k)
	}
	return ss
}
