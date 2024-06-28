package stringutil

import (
	"fmt"
	"strings"
)

func BothEmptyOrNone(str1, str2 string) bool {
	return (str1 == "" && str2 == "") || (str1 != "" && str2 != "")
}

func HasAnyPrefix(str string, prefixes ...string) bool {
	if str != "" {
		for _, p := range prefixes {
			if p != "" {
				if strings.HasPrefix(str, p) {
					return true
				}
			}
		}
	}
	return false
}

func HasAnySuffix(str string, suffixes ...string) bool {
	if str != "" {
		for _, p := range suffixes {
			if p != "" {
				if strings.HasSuffix(str, p) {
					return true
				}
			}
		}
	}
	return false
}

// use typetuil instead
func PointerToString(p *string) string {
	if p != nil {
		return *p
	}
	return ""
}

// 打印字符时不转义
// "\n{\"msgtype\": " -- > "\n{\"msgtype\":
func PrintUnescape(p string) {
	fmt.Println(fmt.Sprintf("%#v", p))
}

// TrimAnySuffix 字符串移除指定前缀(如果有的的话)
func TrimAnyPrefix(str string, prefixes ...string) string {
	if str != "" {
		for _, p := range prefixes {
			str = strings.TrimPrefix(str, p)
		}
	}
	return str
}

// TrimAnySuffix 字符串移除指定后缀(如果有的的话)
func TrimAnySuffix(str string, suffixes ...string) string {
	if str != "" {
		for _, p := range suffixes {
			str = strings.TrimPrefix(str, p)
		}
	}
	return str
}

// AddSuffixIfNotHas 如果字符串没有指定后缀则添加
func AddSuffixIfNotHas(str, suffix string) string {
	if !strings.HasSuffix(str, suffix) {
		str += suffix
	}
	return str
}

// AddPrefixIfNotHas 如果字符串没有指定前缀则添加
func AddPrefixIfNotHas(str, prefix string) string {
	if !strings.HasPrefix(str, prefix) {
		str = prefix + str
	}
	return str
}

// RemoveSubBefore 删除子字符串前的所有字符(不包括子字符串)
func RemoveSubBefore(s, str string) string {
	index := strings.Index(s, str)
	if index == -1 {
		// If str is not found, return the original string
		return s
	}
	return s[index:]
}

// RemoveSubBefore 删除子字符串以及之前所有字符
func RemoveSubAndBefore(s, str string) string {
	index := strings.Index(s, str)
	if index == -1 {
		// If str is not found, return the original string
		return s
	}
	return s[index+len(str):]
}
