package errors_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/wangweihong/gotoolbox/pkg/errors"
)

const (
	ErrEOF = iota + 1000
	ErrCall
)

func init() {
	errors.Register(errors.NewCoder(ErrEOF, 200, map[string]string{errors.MessageLangCNKey: "输入终止", errors.MessageLangENKey: "End of input"}))
	errors.Register(errors.NewCoder(ErrCall, 500, map[string]string{errors.MessageLangCNKey: "请求失败", errors.MessageLangENKey: "call error"}))
}

func TestWithStack(t *testing.T) {
	_, err := os.Stat("noexist")
	fmt.Printf("%+v\n", err)
	// CreateFile noexist: The system cannot find the file specified.
	err = errors.WithStack(err)
	fmt.Printf("%v\n", err)
	// CreateFile noexist: The system cannot find the file specified.

	fmt.Printf("%q\n", err)
	// "CreateFile noexist: The system cannot find the file specified."

	fmt.Printf("%+v\n", err)
	//CreateFile noexist: The system cannot find the file specified.
	//github.com/wangweihong/gotoolbox/pkg/errors.TestWithStack
	//	C:/goprogram/src/github.com/wangweihong/gotoolbox/pkg/errors/example_test.go:25
	//testing.tRunner
	// 	C:/Users/Administrator/go/go1.20.12/src/testing/testing.go:1576
	//runtime.goexit
	// C:/Users/Administrator/go/go1.20.12/src/runtime/asm_amd64.s:1598
}

func TestNew(t *testing.T) {
	e := errors.New("error gen by errors.New()")
	// error gen by errors.New()
	fmt.Printf("%v\n", e)

	// "error gen by errors.New()"
	fmt.Printf("%q\n", e)

	//error gen by errors.New()
	//github.com/wangweihong/gotoolbox/pkg/errors.TestNew
	//	C:/goprogram/src/github.com/wangweihong/gotoolbox/pkg/errors/example_test.go:31
	//testing.tRunner
	//	C:/Users/Administrator/go/go1.20.12/src/testing/testing.go:1576
	//runtime.goexit
	//	C:/Users/Administrator/go/go1.20.12/src/runtime/asm_amd64.s:1598
	fmt.Printf("%+v\n", e)

}

func fn() error {
	e1 := errors.New("error")
	e2 := errors.Wrap(e1, "inner")
	e3 := errors.Wrap(e2, "middle")
	return errors.Wrap(e3, "outer")
}

func TestCause(t *testing.T) {
	err := fn()
	// outer: middle: inner: error
	fmt.Println(err)
	// outer: middle: inner: error
	fmt.Printf("%v\n", err)
	// error
	fmt.Println(errors.Cause(err))

	//error
	//github.com/wangweihong/gotoolbox/pkg/errors_test.fn
	//	C:/goprogram/src/github.com/wangweihong/gotoolbox/pkg/errors/example_test.go:49
	//github.com/wangweihong/gotoolbox/pkg/errors_test.TestCause
	//	C:/goprogram/src/github.com/wangweihong/gotoolbox/pkg/errors/example_test.go:56
	//testing.tRunner
	//	C:/Users/Administrator/go/go1.20.12/src/testing/testing.go:1576
	//runtime.goexit
	//	C:/Users/Administrator/go/go1.20.12/src/runtime/asm_amd64.s:1598

	//inner
	//github.com/wangweihong/gotoolbox/pkg/errors_test.fn
	//	C:/goprogram/src/github.com/wangweihong/gotoolbox/pkg/errors/example_test.go:50
	//github.com/wangweihong/gotoolbox/pkg/errors_test.TestCause
	//	C:/goprogram/src/github.com/wangweihong/gotoolbox/pkg/errors/example_test.go:56
	//testing.tRunner
	//	C:/Users/Administrator/go/go1.20.12/src/testing/testing.go:1576
	//
	//runtime.goexit
	//	C:/Users/Administrator/go/go1.20.12/src/runtime/asm_amd64.s:1598
	//
	//middle
	//github.com/wangweihong/gotoolbox/pkg/errors_test.fn
	//	C:/goprogram/src/github.com/wangweihong/gotoolbox/pkg/errors/example_test.go:51
	//github.com/wangweihong/gotoolbox/pkg/errors_test.TestCause
	//	C:/goprogram/src/github.com/wangweihong/gotoolbox/pkg/errors/example_test.go:56
	//testing.tRunner
	//	C:/Users/Administrator/go/go1.20.12/src/testing/testing.go:1576
	//
	//runtime.goexit
	//	C:/Users/Administrator/go/go1.20.12/src/runtime/asm_amd64.s:1598
	//
	//outer
	//github.com/wangweihong/gotoolbox/pkg/errors_test.fn
	//	C:/goprogram/src/github.com/wangweihong/gotoolbox/pkg/errors/example_test.go:52
	//github.com/wangweihong/gotoolbox/pkg/errors_test.TestCause
	//	C:/goprogram/src/github.com/wangweihong/gotoolbox/pkg/errors/example_test.go:56
	//testing.tRunner
	//	C:/Users/Administrator/go/go1.20.12/src/testing/testing.go:1576
	//
	//runtime.goexit
	//	C:/Users/Administrator/go/go1.20.12/src/runtime/asm_amd64.s:1598
	fmt.Printf("%+v\n", err)

	//error
	//github.com/wangweihong/gotoolbox/pkg/errors_test.fn
	//	C:/goprogram/src/github.com/wangweihong/gotoolbox/pkg/errors/example_test.go:49
	//github.com/wangweihong/gotoolbox/pkg/errors_test.TestCause
	//	C:/goprogram/src/github.com/wangweihong/gotoolbox/pkg/errors/example_test.go:56
	//testing.tRunner
	//	C:/Users/Administrator/go/go1.20.12/src/testing/testing.go:1576
	//runtime.goexit
	//	C:/Users/Administrator/go/go1.20.12/src/runtime/asm_amd64.s:1598
	fmt.Printf("%+v\n", errors.Cause(err))

}

func TestWithCode(t *testing.T) {
	var err error

	err = errors.WithCode(ErrEOF, "this is an error message")
	fmt.Printf("%v\n", err)
	fmt.Println("---------------------")
	fmt.Printf("%+v\n", err)
	fmt.Println("---------------------")
	fmt.Printf("%+#v\n", err)
	fmt.Println("---------------------")

	fmt.Printf("%-#v\n", err)
	fmt.Println("---------------------")

	err = errors.Wrap(err, "this is a wrap error message with error code not change")
	fmt.Printf("%v\n", err)
	fmt.Println("---------------------")

	fmt.Printf("%+v\n", err)
	fmt.Println("---------------------")

	fmt.Printf("%+#v\n", err)
	fmt.Println("---------------------")

	fmt.Printf("%-#v\n", err)
	fmt.Println("---------------------")

	err = errors.WrapCodeF(err, ErrCall, "this is a wrap error message with new error code")
	fmt.Printf("%v\n", err)
	fmt.Println("---------------------")

	fmt.Printf("%+v\n", err)
	fmt.Println("---------------------")

	fmt.Printf("%+#v\n", err)
	fmt.Println("---------------------")

	fmt.Printf("%-#v\n", err)

}

func wrapCode() error {
	var err error

	err = errors.New("this is error")
	err = errors.WrapCodeF(err, ErrEOF, "some message")
	err = errors.WrapCode(err, ErrCall)
	return err
}

func TestWrapCode(t *testing.T) {
	err := wrapCode()
	fmt.Printf("%s\n", err)

	fmt.Printf("%v\n", err)
	fmt.Println("---------------------")
	fmt.Printf("%+v\n", err)
	fmt.Println("---------------------")
	fmt.Printf("%+#v\n", err)
	fmt.Println("---------------------")

	fmt.Printf("%-#v\n", err)
	fmt.Println("---------------------")
}

func TestWrapErrorWithCode(t *testing.T) {
	err := wrapCode()
	err = errors.Wrap(err, "wrap with withcode")
	fmt.Printf("%s\n", err)

	fmt.Printf("%v\n", err)
	fmt.Println("---------------------")
	fmt.Printf("%+v\n", err)
	fmt.Println("---------------------")
	fmt.Printf("%+#v\n", err)
	fmt.Println("---------------------")

	fmt.Printf("%-#v\n", err)
	fmt.Println("---------------------")
}

func a() error {
	return fmt.Errorf("error a")
}

func b() error {
	err := a()
	return errors.Wrap(err, "error b")
}

func c() error {
	err := a()
	return errors.WithStack(err)
}

func d() error {
	_, err := os.Stat("noexist")
	return errors.WithStack(err)
}

func e() error {
	return d()
}

func TestDDD(t *testing.T) {
	e := d()
	fmt.Printf("%+v\n", e)
}

func TestAAA(t *testing.T) {
	err := b()

	// error b: error a
	fmt.Println(err)
	// error a
	fmt.Println(errors.Cause(err))

	// error a
	err = c()
	fmt.Println(err)

	//	github.com/wangweihong/gotoolbox/pkg/httpcli/def.c
	//C:/goprogram/src/github.com/wangweihong/gotoolbox/pkg/httpcli/def/multipart_test.go:23
	//	github.com/wangweihong/gotoolbox/pkg/httpcli/def.TestAAA
	//C:/goprogram/src/github.com/wangweihong/gotoolbox/pkg/httpcli/def/multipart_test.go:35
	//	testing.tRunner
	//C:/Users/Administrator/go/go1.20.12/src/testing/testing.go:1576
	//	runtime.goexit
	//C:/Users/Administrator/go/go1.20.12/src/runtime/asm_amd64.s:1598
	fmt.Printf("%+v\n", err)

	fmt.Printf("%+v\n", errors.Cause(err))

	fmt.Printf("%v\n", errors.New("message 1"))

	de := d()
	fmt.Printf("%+v", de)

	//	github.com/wangweihong/gotoolbox/pkg/errors.d
	//C:/goprogram/src/github.com/wangweihong/gotoolbox/pkg/errors/origin_test.go:29
	//	github.com/wangweihong/gotoolbox/pkg/errors.e
	//C:/goprogram/src/github.com/wangweihong/gotoolbox/pkg/errors/origin_test.go:33
	//	github.com/wangweihong/gotoolbox/pkg/errors.TestAAA
	//C:/goprogram/src/github.com/wangweihong/gotoolbox/pkg/errors/origin_test.go:61
	//	testing.tRunner
	//C:/Users/Administrator/go/go1.20.12/src/testing/testing.go:1576
	//	runtime.goexit

	// 会自动记录e的调用栈, e不需要执行()操作。
	ee := e()
	fmt.Printf("%+v\n", ee)

	/*
	   message1
	   github.com/wangweihong/gotoolbox/pkg/errors.TestAAA
	   	C:/goprogram/src/github.com/wangweihong/gotoolbox/pkg/errors/origin_test.go:75
	   testing.tRunner
	   	C:/Users/Administrator/go/go1.20.12/src/testing/testing.go:1576
	   runtime.goexit
	   	C:/Users/Administrator/go/go1.20.12/src/runtime/asm_amd64.s:1598
	*/
	we := errors.New("message1")
	fmt.Printf("%+v\n", we)

}

func TestBB(t *testing.T) {
	e := errors.New("error1")
	//// error1
	//fmt.Printf("%v\n", e)
	////"error1"
	//fmt.Printf("%q\n", e)
	//fmt.Println("---------------------------")
	//
	//fmt.Printf("%+v\n", e)
	//fmt.Println("---------------------------")
	//
	//fmt.Printf("%+#v\n", e)
	//fmt.Println("---------------------------")

	e2 := errors.Wrap(e, "error2")
	// error2: error1
	fmt.Printf("%v\n", e2)

	fmt.Println("---------------------------")
	//	error1
	//	github.com/wangweihong/gotoolbox/pkg/errors.TestBB
	//C:/goprogram/src/github.com/wangweihong/gotoolbox/pkg/errors/origin_test.go:94
	//	testing.tRunner
	//C:/Users/Administrator/go/go1.20.12/src/testing/testing.go:1576
	//	runtime.goexit
	//C:/Users/Administrator/go/go1.20.12/src/runtime/asm_amd64.s:1598
	//	error2
	//	github.com/wangweihong/gotoolbox/pkg/errors.TestBB
	//C:/goprogram/src/github.com/wangweihong/gotoolbox/pkg/errors/origin_test.go:104
	//	testing.tRunner
	//C:/Users/Administrator/go/go1.20.12/src/testing/testing.go:1576
	//	runtime.goexit
	//C:/Users/Administrator/go/go1.20.12/src/runtime/asm_amd64.s:1598
	fmt.Printf("%+v\n", e2)
	fmt.Println("---------------------------")
	//	error1
	//	github.com/wangweihong/gotoolbox/pkg/errors.TestBB
	//C:/goprogram/src/github.com/wangweihong/gotoolbox/pkg/errors/origin_test.go:94
	//	testing.tRunner
	//C:/Users/Administrator/go/go1.20.12/src/testing/testing.go:1576
	//	runtime.goexit
	//C:/Users/Administrator/go/go1.20.12/src/runtime/asm_amd64.s:1598
	//	error2
	//	github.com/wangweihong/gotoolbox/pkg/errors.TestBB
	//C:/goprogram/src/github.com/wangweihong/gotoolbox/pkg/errors/origin_test.go:104
	//	testing.tRunner
	//C:/Users/Administrator/go/go1.20.12/src/testing/testing.go:1576
	//	runtime.goexit
	//C:/Users/Administrator/go/go1.20.12/src/runtime/asm_amd64.s:1598
	fmt.Printf("%+#v\n", e2)
	fmt.Println("---------------------------")

	ue := errors.Unwrap(e2)
	//error2: error1
	fmt.Printf("%v\n", ue)
	fmt.Println("---------------------------")

	//	error1
	//	github.com/wangweihong/gotoolbox/pkg/errors.TestBB
	//C:/goprogram/src/github.com/wangweihong/gotoolbox/pkg/errors/origin_test.go:94
	//	testing.tRunner
	//C:/Users/Administrator/go/go1.20.12/src/testing/testing.go:1576
	//	runtime.goexit
	//C:/Users/Administrator/go/go1.20.12/src/runtime/asm_amd64.s:1598
	//	error2
	fmt.Printf("%+v\n", ue)
	fmt.Println("---------------------------")
	me := errors.WithMessage(e, "message1")
	// message1: error1
	fmt.Printf("%v\n", me)
	//	error1
	//	github.com/wangweihong/gotoolbox/pkg/errors.TestBB
	//C:/goprogram/src/github.com/wangweihong/gotoolbox/pkg/errors/origin_test.go:94
	//	testing.tRunner
	//C:/Users/Administrator/go/go1.20.12/src/testing/testing.go:1576
	//	runtime.goexit
	//C:/Users/Administrator/go/go1.20.12/src/runtime/asm_amd64.s:1598
	//	message1
	fmt.Printf("%+v\n", me)
}
