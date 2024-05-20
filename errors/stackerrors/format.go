package stackerrors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Format implements fmt.Formatter. https://golang.org/pkg/fmt/#hdr-Printing
//
// Verbs:
//
//	%s  - Returns the user-safe error string mapped to the error code or
//	  ┊   the error message if none is specified.
//	%v      Alias for %s
//
// Flags:
//
//	#      JSON formatted output, useful for logging
//	+      Output full error stack details, useful for debugging
//
// Examples:
//
//	     %s:    OpenFileError:file not exist
//			%q:    "OpenFileError:file not exist"
//	     %v:    OpenFileError:file not exist
//			%+v:   OpenFileError:file not exist
//
// [host:127.0.0.1,pid:8536,module:testing,code:101,file:error_test.go,func:1,line:55]
//
//	%#v:   {"code":101,"desc":"file not exist","message":{"cn":"访问文件失败","en":"OpenFileError"}}
//	%#+v:  {"code":101,"desc":"file not
//
// exist","http":200,"message":{"cn":"访问文件失败","en":"OpenFileError"},"stack":[{"host":"127.0.0.1","pid":"6716","module":"testing","code":"101","file_name":"error_test.go","func_name":"1","line":"55"}]}
//
//gofmt:disable
//gofmt:enable
func (m *stackerrors.WithStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		str := bytes.NewBuffer([]byte{})

		var (
			flagDetail bool
			modeJSON   bool
		)

		if s.Flag('#') {
			modeJSON = true
		}

		if s.Flag('+') {
			flagDetail = true
		}

		if modeJSON {
			byteData, _ := json.Marshal(m.formatJson(flagDetail))
			str.Write(byteData)
		} else {
			if flagDetail {
				fmt.Fprintf(str, "%s %s", m.Error(), m.stack)
			} else {
				fmt.Fprintf(str, "%s", m.Error())
			}
		}
		fmt.Fprintf(s, "%s", strings.Trim(str.String(), "\r\n\t"))
	case 's':
		_, _ = io.WriteString(s, m.Error())
	case 'q':
		fmt.Fprintf(s, "%q", m.Error())
	}
}

func (m *stackerrors.WithStack) formatJson(detail bool) map[string]interface{} {
	data := m.ToBasicJson()
	if detail {
		data = m.ToDetailJson()
	}
	return data
}
