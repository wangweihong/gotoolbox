package stackerrors_test

import (
	"fmt"
	"testing"

	"github.com/wangweihong/gotoolbox/errors/stackerrors"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/wangweihong/gotoolbox/errors"
)

func TestNewDesc(t *testing.T) {
	e := stackerrors.NewDesc(EcodeOpenFile, "some thing wrong")
	fmt.Printf("%#+v\n", e)
}

func TestNewError(t *testing.T) {
	Convey("TestNewError", t, func() {
		e := stackerrors.NewDesc(EcodeOpenFile, "some thing wrong")
		e2 := stackerrors.New(EcodeWriteFile, e)
		So(len(e2.StackInfo()), ShouldEqual, 2)

		e3 := stackerrors.New(EcodeWriteFile, nil)
		fmt.Printf("%#+v\n", e3)

	})
}

func TestUpdateStack(t *testing.T) {
	var e error
	e = stackerrors.NewDesc(EcodeOpenFile, "some thing wrong")
	e = stackerrors.UpdateStack(e)
	//fmt.Printf("%#+v\n", e)

	e2 := stackerrors.UpdateStack(nil)
	fmt.Printf("%#+v\n", e2)

}

type fakeModule struct {
	name string
	ip   string
	pid  int
}

func (s fakeModule) PID() int {
	return s.pid
}

func (s fakeModule) IP() string {
	return s.ip
}

func (s fakeModule) Name() string {
	return s.name
}

func (s fakeModule) String() string {
	return fmt.Sprintf("host:%s,pid:%d,module:%s", s.IP(), s.PID(), s.Name())
}

const (
	EcodeWriteFile = 100
	EcodeOpenFile  = 101
	EcodeReadFile  = 102
)

func init() {
	errors.UpdateModuleInfo(&fakeModule{
		name: "testing",
		ip:   "127.0.0.1",
		pid:  8536,
	})
	stackerrors.MustRegister(stackerrors.NewCoder(EcodeWriteFile, map[string]string{
		stackerrors.MessageLangENKey: "WriteFileError",
		stackerrors.MessageLangCNKey: "写文件失败",
	}))
	stackerrors.MustRegister(stackerrors.NewCoder(EcodeOpenFile, map[string]string{
		stackerrors.MessageLangENKey: "OpenFileError",
		stackerrors.MessageLangCNKey: "访问文件失败",
	}))
	stackerrors.MustRegister(stackerrors.NewCoder(EcodeReadFile, map[string]string{
		stackerrors.MessageLangENKey: "ReadFileError",
		stackerrors.MessageLangCNKey: "读文件失败",
	}))
}
