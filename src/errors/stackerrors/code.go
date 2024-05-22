//nolint:errorlint
package stackerrors

import (
	"fmt"
	"strings"
	"sync"
)

const (
	MessageLangCNKey = "MessageCN"
	MessageLangENKey = "MessageEN"
)

var unknown = fundamental{
	code: 1,
	message: map[string]string{
		MessageLangCNKey: "未知错误码",
		MessageLangENKey: "unknown error code",
	},
}

var success = fundamental{
	code: 0,
	message: map[string]string{
		MessageLangCNKey: "成功",
		MessageLangENKey: "success",
	},
}

// fundamental is an error that has a message and a stack, but no caller.
type fundamental struct {
	code    int               // 状态码
	message map[string]string // 信息
}

// Coder defines an interface for an error code detail information.
type Coder interface {
	// External (user) facing error message.
	Message() map[string]string

	// Code returns the code of the coder
	Code() int
}

func (f fundamental) Message() map[string]string {
	return f.message
}

func (f fundamental) Code() int {
	return f.code
}

func (f fundamental) MessageCN() string {
	if f.message != nil {
		msg := f.message[MessageLangCNKey]
		return msg
	}
	return ""
}

func (f fundamental) MessageEN() string {
	if f.message != nil {
		msg := f.message[MessageLangENKey]
		return msg
	}
	return ""
}

// errorlint:ignore
// IsCode reports whether any error in err's chain contains the given error code.
func IsCode(err error, code int) bool {
	if err != nil {
		if v, ok := err.(*WithStack); ok {
			if v.Code() == code {
				return true
			}
		}
	}

	return false
}

// codes contains a map of error codes to metadata.
var (
	codes   = map[int]Coder{}
	codeMux = &sync.RWMutex{}
)

func registerPre(coder Coder) {
	if coder.Code() == unknown.code {
		panic("code `1` is reserved by `errors` as unknown error code")
	}

	if coder.Code() == success.code {
		panic("code `0` is reserved by `errors` as success code")
	}

	if coder.Message() == nil {
		panic(
			fmt.Sprintf(
				"coder `%v` has no message map with key `%v` and `%v`",
				coder.Code(),
				MessageLangCNKey,
				MessageLangENKey,
			),
		)
	}

	if v := coder.Message()[MessageLangENKey]; strings.TrimSpace(v) == "" {
		panic(fmt.Sprintf("coder `%v` has message map  key `%v` value is empty", coder.Code(), MessageLangENKey))
	}

	if v := coder.Message()[MessageLangCNKey]; strings.TrimSpace(v) == "" {
		panic(fmt.Sprintf("coder `%v` has message map  key `%v` value is empty", coder.Code(), MessageLangCNKey))
	}
}

// Register register a user define error code.
// It will override the exist code.
func Register(coder Coder) {
	registerPre(coder)

	codeMux.Lock()
	defer codeMux.Unlock()

	codes[coder.Code()] = coder
}

// MustRegister register a user define error code.
// It will panic when the same Code already exist.
func MustRegister(coder Coder) {
	registerPre(coder)

	codeMux.Lock()
	defer codeMux.Unlock()

	if _, ok := codes[coder.Code()]; ok {
		panic(fmt.Sprintf("code: %d already exist", coder.Code()))
	}

	codes[coder.Code()] = coder
}

// ParseCoder parse any error into *withCode.
// nil error will return nil direct.
// None WithStack error will be parsed as ErrUnknown.
func ParseCoder(err error) Coder {
	if err == nil {
		return nil
	}

	if v, ok := err.(*WithStack); ok {
		if coder, ok := codes[v.Code()]; ok {
			return coder
		}
	}

	return unknown
}

func Unknown() Coder {
	return unknown
}

func Success() Coder {
	return success
}

// errorlint:ignore
func MessageEn(err error) string {
	if err == nil {
		return ""
	}
	if st, ok := err.(Coder); ok {
		return st.Message()[MessageLangENKey]
	}
	return "unknown error message"
}

// errorlint:ignore
func MessageCN(err error) string {
	if err == nil {
		return ""
	}
	if st, ok := err.(Coder); ok {
		return st.Message()[MessageLangCNKey]
	}
	return "未知错误信息"
}

func Code(err error) int {
	if err == nil {
		return success.code
	}
	if st, ok := err.(Coder); ok {
		return st.Code()
	}
	return unknown.code
}

func NewCoder(code int, message map[string]string) Coder {
	return fundamental{
		code:    code,
		message: message,
	}
}

//nolint:gochecknoinits
func init() {
	codes[success.code] = success
	codes[unknown.code] = unknown
}
