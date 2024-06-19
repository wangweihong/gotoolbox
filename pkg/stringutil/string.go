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

func TrimAnyPrefix(str string, prefixes ...string) string {
	if str != "" {
		for _, p := range prefixes {
			str = strings.TrimPrefix(str, p)
		}
	}
	return str
}

func TrimAnySuffix(str string, suffixes ...string) string {
	if str != "" {
		for _, p := range suffixes {
			str = strings.TrimPrefix(str, p)
		}
	}
	return str
}

func AddSuffixIfNotHas(str, suffix string) string {
	if !strings.HasSuffix(str, suffix) {
		str += suffix
	}
	return str
}

func AddPrefixIfNotHas(str, prefix string) string {
	if !strings.HasPrefix(str, prefix) {
		str = prefix + str
	}
	return str
}
