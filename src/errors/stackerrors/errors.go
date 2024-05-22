//nolint:errorlint
package stackerrors

import (
	"fmt"
	"strings"

	"github.com/wangweihong/gotoolbox/src/errors"
)

// NewDesc generate a new WithStack error with `code` and `desc`.
func NewDesc(code int, desc string) *WithStack {
	codeMux.RLock()
	defer codeMux.RUnlock()

	errStack := &WithStack{
		stack:       []string{newStack(code, errors.Caller())},
		description: desc,
	}

	coder, exist := codes[code]
	if !exist {
		errStack.Coder = unknown
		return errStack
	}
	errStack.Coder = coder
	return errStack
}

// NewStack generate a new WithStack error with `code` , `desc`,`stack`.
func NewStack(code int, desc string, stack []string) *WithStack {
	codeMux.RLock()
	defer codeMux.RUnlock()

	stack = append(stack, newStack(code, errors.Caller()))
	errStack := &WithStack{
		stack:       stack,
		description: desc,
	}

	coder, exist := codes[code]
	if !exist {
		errStack.Coder = unknown
		return errStack
	}
	errStack.Coder = coder
	return errStack
}

// NewF generate a new WithStack error with `code` and desc `format+arg...`.
func NewF(code int, format string, args ...interface{}) *WithStack {
	codeMux.RLock()
	defer codeMux.RUnlock()

	errStack := &WithStack{
		stack:       []string{newStack(code, errors.Caller())},
		description: fmt.Sprintf(format, args...),
	}

	coder, exist := codes[code]
	if !exist {
		errStack.Coder = unknown
		return errStack
	}
	errStack.Coder = coder
	return errStack
}

// WrapError generate a new WithStack error with `code` and desc `format+arg...`
// if err not nil, inherit its stack and error message, replace origin code.
func New(code int, err error) *WithStack {
	codeMux.RLock()
	defer codeMux.RUnlock()

	errStack := &WithStack{
		stack:       []string{newStack(code, errors.Caller())},
		description: "",
	}

	coder, exist := codes[code]
	if !exist {
		errStack.Coder = unknown
		if err != nil {
			errStack.description = err.Error()
		}
		return errStack
	}

	if err != nil {
		if st, ok := err.(*WithStack); ok {
			errStack.stack = append(st.Stack(), errStack.stack...)
			errStack.description = st.description
		} else {
			errStack.description = err.Error()
		}
	}
	errStack.Coder = coder
	return errStack
}

// UpdateStack add a new layer to err's caller stack.
func UpdateStack(err error) error {
	if err != nil {
		errStack := FromError(err)
		if errStack != nil {
			errStack.stack = append(errStack.stack, newStack(errStack.Code(), errors.Caller()))
			return errStack
		}
	}
	return nil
}

var _ error = &WithStack{}

type WithStack struct {
	Coder
	// 状态栈
	stack []string
	// 描述
	description string
}

func (m *WithStack) Stack() []string {
	if m != nil {
		return m.stack
	}
	return nil
}

func (m *WithStack) Description() string {
	if m != nil {
		return m.description
	}
	return ""
}

func (m *WithStack) Detail() string {
	siList := m.StackInfo()
	callChain := ""
	lastService := ""
	for i, si := range siList {
		service := fmt.Sprintf("([%s:%s]", si.Host, si.Module)
		if i == 0 {
			callChain = fmt.Sprintf("%s<%s:%s>)", si.FuncName, si.FileName, si.Line)
		} else if service != lastService && lastService != "" {
			callChain = si.FuncName + ")->" + lastService + callChain
		} else {
			callChain = si.FuncName + "->" + callChain
		}
		lastService = service
	}
	callChain = lastService + callChain
	if m.description != "" {
		return fmt.Sprintf("stack:%s,code:%d,message:%s,desc:%s", callChain, m.Code(), m.Message(), m.description)
	}
	return fmt.Sprintf("stack:%s,code:%d,message:%s", callChain, m.Code(), m.Message())
}

type StackInfo struct {
	Host     string `json:"host"`
	PID      string `json:"pid"`
	Module   string `json:"module"`
	Code     string `json:"code"`
	FileName string `json:"file_name"`
	FuncName string `json:"func_name"`
	Line     string `json:"line"`
}

func (m *WithStack) StackInfo() []StackInfo {
	siList := make([]StackInfo, 0, len(m.stack))
	for _, str := range m.stack {
		si := StackInfo{}
		slist := strings.Split(str, ",")
		for _, s := range slist {
			if strings.HasPrefix(s, "host:") {
				si.Host = strings.TrimPrefix(s, "host:")
			}
			if strings.HasPrefix(s, "pid:") {
				si.PID = strings.TrimPrefix(s, "pid:")
			}
			if strings.HasPrefix(s, "module:") {
				si.Module = strings.TrimPrefix(s, "module:")
			}
			if strings.HasPrefix(s, "code:") {
				si.Code = strings.TrimPrefix(s, "code:")
			}
			if strings.HasPrefix(s, "file:") {
				si.FileName = strings.TrimPrefix(s, "file:")
			}
			if strings.HasPrefix(s, "func:") {
				si.FuncName = strings.TrimPrefix(s, "func:")
			}
			if strings.HasPrefix(s, "line:") {
				si.Line = strings.TrimPrefix(s, "line:")
			}
		}
		siList = append(siList, si)
	}
	return siList
}

func (m WithStack) Error() string {
	return fmt.Sprintf("%v:%v", m.Message()[MessageLangENKey], m.description)
}

func (m WithStack) ToBasicJson() map[string]interface{} {
	out := make(map[string]interface{})
	out["desc"] = m.description
	out["message"] = m.Message()
	out["code"] = m.Code()

	return out
}

func (m WithStack) ToDetailJson() map[string]interface{} {
	out := make(map[string]interface{})
	out["desc"] = m.description
	out["stack"] = m.StackInfo()
	out["message"] = m.Message()
	out["code"] = m.Code()
	return out
}

// FromError parse any error into *WithStack.
// nil error will return nil directly, caller should handle nil *WithStack.
// None WithStack error will be parsed as ErrUnknown.
// NOTE: `*WithStack is nil` doesn't equal to `error is nil`.
func FromError(err error) *WithStack {
	return fromError(err)
}

func fromError(err error) *WithStack {
	if err == nil {
		return nil
	}

	if v, ok := err.(*WithStack); ok {
		return v
	}

	return &WithStack{
		Coder:       unknown,
		stack:       []string{newStack(unknown.code, errors.Caller())},
		description: err.Error(),
	}
}

func newStack(code int, caller string) string {
	return errors.ModuleString() + fmt.Sprintf(",code:%d,", code) + caller
}

type ErrorInfo struct {
	Code    int
	Message map[string]string
	Stacks  []string
	Desc    string
}

func (m WithStack) ToErrorInfo() *ErrorInfo {
	return &ErrorInfo{
		Code:    m.Code(),
		Message: m.Message(),
		Stacks:  m.Stack(),
		Desc:    m.Description(),
	}
}
