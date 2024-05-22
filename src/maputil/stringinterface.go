package maputil

import "reflect"

type StringInterfaceMap map[string]interface{}

func NewStringInterfaceMap() StringInterfaceMap {
	nm := make(map[string]interface{})
	return nm
}

func (m StringInterfaceMap) DeepCopy() StringInterfaceMap {
	o := make(map[string]interface{}, len(m))
	for k, v := range m {
		o[k] = v
	}
	return o
}

func (m StringInterfaceMap) Init() StringInterfaceMap {
	if m == nil {
		return make(map[string]interface{})
	}
	return m
}

func (m StringInterfaceMap) Delete(key string) {
	if m == nil {
		return
	}
	delete(m, key)
}

func (m StringInterfaceMap) DeleteIfKey(condition func(string) bool) {
	if m == nil {
		return
	}

	for k := range m {
		if condition(k) {
			delete(m, k)
		}
	}
}

func (m StringInterfaceMap) DeleteIfValue(condition func(interface{}) bool) {
	if m == nil {
		return
	}

	for k, v := range m {
		if condition(v) {
			delete(m, k)
		}
	}
}

func (m StringInterfaceMap) Has(key string) bool {
	if m != nil {
		if _, exist := m[key]; exist {
			return true
		}
	}
	return false
}

func (m StringInterfaceMap) HasKeyAndValue(key string, value interface{}) bool {
	if m != nil {
		v, exist := m[key]
		if !exist {
			return false
		}

		if reflect.DeepEqual(v, value) {
			return true
		}
	}
	return false
}

func (m StringInterfaceMap) IsSuperSet(m2 map[string]interface{}) bool {
	if m2 == nil || m == nil {
		return false
	}

	for k, v := range m2 {
		if !m.HasKeyAndValue(k, v) {
			return false
		}
	}
	return true
}

func (m StringInterfaceMap) Set(key string, value interface{}) StringInterfaceMap {
	if m == nil {
		o := make(map[string]interface{})
		o[key] = value
		return o
	}
	m[key] = value
	return m
}

func (m StringInterfaceMap) Get(key string) interface{} {
	if m == nil {
		return nil
	}
	v, _ := m[key]
	return v
}

func (m StringInterfaceMap) Keys() []string {
	if m == nil {
		return []string{}
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func (m StringInterfaceMap) Equal(m2 map[string]interface{}) bool {
	if len(m) != len(m2) {
		return false
	}

	for k1, v1 := range m {
		v2, ok := m2[k1]
		if !ok {
			return false
		}

		if !reflect.DeepEqual(v1, v2) {
			return false
		}
	}
	return true
}
