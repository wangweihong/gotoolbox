package def

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"

	"github.com/wangweihong/gotoolbox/pkg/typeutil"
)

var quoteEscape = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscape.Replace(s)
}

type FilePart struct {
	Headers textproto.MIMEHeader
	Content *os.File
}

func NewFilePart(content *os.File) *FilePart {
	return &FilePart{
		Content: content,
	}
}

func NewFilePartWithContentType(content *os.File, contentType string) *FilePart {
	headers := make(textproto.MIMEHeader)
	headers.Set("Content-Type", contentType)

	return &FilePart{
		Headers: headers,
		Content: content,
	}
}

func (f FilePart) Write(w *multipart.Writer, name string) error {
	var h textproto.MIMEHeader
	if f.Headers != nil {
		h = f.Headers
	} else {
		h = make(textproto.MIMEHeader)
	}

	filename := filepath.Base(f.Content.Name())
	if filename == "" {
		return errors.New("failed to obtain filename from: " + f.Content.Name())
	}

	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(name), escapeQuotes(filename)))

	if f.Headers.Get("Content-Type") == "" {
		h.Set("Content-Type", "application/octet-stream")
	}

	writer, err := w.CreatePart(h)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, f.Content)
	return err
}

type MultiPart struct {
	Content interface{}
}

func NewMultiPart(content interface{}) *MultiPart {
	return &MultiPart{
		Content: content,
	}
}

func (m MultiPart) Write(w *multipart.Writer, name string) error {
	err := w.WriteField(name, typeutil.ConvertInterfaceToString(m.Content))
	return err
}

type FormData interface {
	Write(*multipart.Writer, string) error
}
