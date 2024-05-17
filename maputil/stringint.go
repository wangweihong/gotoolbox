package maputil

import "github.com/wangweihong/gotoolbox/sets"

type StringIntMap map[string]int

func (m StringIntMap) DeepCopy() StringIntMap {
	o := make(map[string]int, len(m))
	for k, v := range m {
		o[k] = v
	}
	return o
}

func (m StringIntMap) Init() StringIntMap {
	if m == nil {
		return make(map[string]int)
	}
	return m
}

func (m StringIntMap) Delete(key string) {
	if m == nil {
		return
	}
	delete(m, key)
}

func (m StringIntMap) DeleteIfKey(condition func(string) bool) {
	if m == nil {
		return
	}

	for k := range m {
		if condition(k) {
			delete(m, k)
		}
	}
}

func (m StringIntMap) DeleteIfValue(condition func(int) bool) {
	if m == nil {
		return
	}

	for k, v := range m {
		if condition(v) {
			delete(m, k)
		}
	}
}

func (m StringIntMap) Has(key string) bool {
	if m != nil {
		if _, exist := m[key]; exist {
			return true
		}
	}
	return false
}

func (m StringIntMap) Set(key string, value int) StringIntMap {
	if m == nil {
		o := make(map[string]int)
		o[key] = value
		return o
	}
	m[key] = value
	return m
}

func (m StringIntMap) Get(key string) int {
	if m == nil {
		return 0
	}
	v, _ := m[key]
	return v
}

func (m StringIntMap) Keys() []string {
	if m == nil {
		return []string{}
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func (m StringIntMap) ToSetString() sets.String {
	ss := sets.NewString()
	if m == nil {
		return ss
	}
	for k := range m {
		ss.Insert(k)
	}
	return ss
}

func (m StringIntMap) Equal(m2 map[string]int) bool {
	if len(m) != len(m2) {
		return false
	}

	for k1, v1 := range m {
		v2, ok := m2[k1]
		if !ok {
			return false
		}

		if v1 != v2 {
			return false
		}
	}
	return true
}
