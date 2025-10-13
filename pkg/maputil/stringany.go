package maputil

import "reflect"

type StringAnyMap map[string]any

func ToStringAnyMap(d map[string]any) StringAnyMap {
	return StringAnyMap(d)
}

// TODO : lock
func NewStringAnyMap() StringAnyMap {
	nm := make(map[string]any)
	return nm
}

func (m StringAnyMap) DeepCopy() StringAnyMap {
	o := make(map[string]any, len(m))
	for k, v := range m {
		o[k] = v
	}
	return o
}

func (m StringAnyMap) Init() StringAnyMap {
	if m == nil {
		return make(map[string]any)
	}
	return m
}

func (m StringAnyMap) Delete(key string) {
	if m == nil {
		return
	}
	delete(m, key)
}

func (m StringAnyMap) DeleteIfKey(condition func(string) bool) {
	if m == nil {
		return
	}

	for k := range m {
		if condition(k) {
			delete(m, k)
		}
	}
}

func (m StringAnyMap) DeleteIfValue(condition func(any) bool) {
	if m == nil {
		return
	}

	for k, v := range m {
		if condition(v) {
			delete(m, k)
		}
	}
}

func (m StringAnyMap) Has(key string) bool {
	if m != nil {
		if _, exist := m[key]; exist {
			return true
		}
	}
	return false
}

func (m StringAnyMap) HasKeyAndValue(key string, value any) bool {
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

func (m StringAnyMap) IsSuperSet(m2 map[string]any) bool {
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

func (m StringAnyMap) Set(key string, value any) StringAnyMap {
	if m == nil {
		o := make(map[string]any)
		o[key] = value
		return o
	}
	m[key] = value
	return m
}

func (m StringAnyMap) Get(key string) any {
	if m == nil {
		return nil
	}
	v, _ := m[key]
	return v
}

func (m StringAnyMap) Keys() []string {
	if m == nil {
		return []string{}
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func (m StringAnyMap) Equal(m2 map[string]any) bool {
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

func (m StringAnyMap) GetInt(key string) int {
	if m != nil && key != "" {
		d := m.Get(key)
		if d != nil {
			dv, _ := d.(int)
			return dv
		}
	}
	return 0
}

func (m StringAnyMap) GetString(key string) string {
	if m != nil && key != "" {
		d := m.Get(key)
		if d != nil {
			dv, _ := d.(string)
			return dv
		}
	}
	return ""
}

func (m StringAnyMap) GetMap(key string) map[string]any {
	if m != nil && key != "" {
		d := m.Get(key)
		if d != nil {
			dv, _ := d.(map[string]any)
			return dv
		}
	}
	return nil
}

func (m StringAnyMap) GetMapSlice(key string) []map[string]any {
	if m != nil && key != "" {
		d := m.Get(key)
		if d != nil {
			dv, _ := d.([]any)
			if dv != nil {
				m := make([]map[string]any, 0, len(dv))
				for _, v := range dv {
					vd, _ := v.(map[string]any)
					if vd != nil {
						m = append(m, vd)
					}
				}
				return m
			}
		}
	}
	return nil
}

