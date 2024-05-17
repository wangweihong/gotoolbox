package callerutil

import (
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func CallerDepth(depth int) (string, int, string) {
	if depth < 2 {
		depth = 2
	}

	funcName := "???"

	pc, file, line, ok := runtime.Caller(depth)
	if !ok {
		return "???", 0, funcName
	}
	if fp := runtime.FuncForPC(pc); fp != nil {
		funcName = fp.Name()
	}

	dir, filename := path.Split(file)
	// show package name for error stack
	if dir != "" {
		parent := filepath.Base(dir)
		filename = filepath.Join(parent, filename)
	}

	fileList := strings.Split(funcName, ".")
	funcName = fileList[len(fileList)-1]
	return filename, line, funcName
}
