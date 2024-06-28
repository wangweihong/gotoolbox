package errors_test

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/json"
)

func init() {
	errors.UpdateModuleInfo(errors.NewModuleGetter("github.com/wangweihong/gotoolbox", "127.0.0.1", 123))
}

func example() error {
	e := errors.NewStatus(ErrEOF, "error example")
	return e
}

func TestErrorStack(t *testing.T) {
	Convey("check errors function line number", t, func() {
		// status
		Convey("WithStack error", func() {

			e := example()
			e = errors.WrapStatus(e, ErrEOF)
			e = errors.WrapStatus(e, ErrCall)

			s2 := errors.ToStatus(e)
			//模拟跨服务错误转换
			e2 := s2.Error()
			s3 := errors.ToStatus(e2)
			So(s3.HTTPStatus, ShouldEqual, 500)
			So(s3.Code, ShouldEqual, ErrCall)
			So(len(s3.Cause), ShouldEqual, 2)
			json.PrintObject(s3)

			/*
					{
					"HTTPStatus": 500,
					"Code": 1001,
					"Message": {
						"CN": "请求失败",
						"EN": "call error"
					},
					"Desc": "call error:call error:End of input:End of input:error example",
					"Cause": [
						{
							"Service": {
								"Host": "127.0.0.1",
								"Pid": 123,
								"Name": "github.com/wangweihong/gotoolbox"
							},
							"Stacks": [
								"github.com/wangweihong/gotoolbox/pkg/errors/status_test.go:31 github.com/wangweihong/gotoolbox/pkg/errors_test.TestErrorStack.func1.1'",
								"github.com/wangweihong/gotoolbox/pkg/errors/status_test.go:23 github.com/wangweihong/gotoolbox/pkg/errors_test.TestErrorStack.func1'",
								"github.com/wangweihong/gotoolbox/pkg/errors/status_test.go:21 github.com/wangweihong/gotoolbox/pkg/errors_test.TestErrorStack'"
							]
						},
						{
							"Service": {
								"Host": "127.0.0.1",
								"Pid": 123,
								"Name": "github.com/wangweihong/gotoolbox"
							},
							"Stacks": [
								"github.com/wangweihong/gotoolbox/pkg/errors/status_test.go:27 github.com/wangweihong/gotoolbox/pkg/errors_test.TestErrorStack.func1.1'",
								"github.com/wangweihong/gotoolbox/pkg/errors/status_test.go:26 github.com/wangweihong/gotoolbox/pkg/errors_test.TestErrorStack.func1.1'",
								"github.com/wangweihong/gotoolbox/pkg/errors/status_test.go:16 github.com/wangweihong/gotoolbox/pkg/errors_test.example'",
								"github.com/wangweihong/gotoolbox/pkg/errors/status_test.go:25 github.com/wangweihong/gotoolbox/pkg/errors_test.TestErrorStack.func1.1'",
								"github.com/wangweihong/gotoolbox/pkg/errors/status_test.go:23 github.com/wangweihong/gotoolbox/pkg/errors_test.TestErrorStack.func1'",
								"github.com/wangweihong/gotoolbox/pkg/errors/status_test.go:21 github.com/wangweihong/gotoolbox/pkg/errors_test.TestErrorStack'"
							]
						}
					]
				}
			*/
			fmt.Println(e2.Error())
			fmt.Println(errors.Cause(e))
			fmt.Println(e2.Error())

		})
	})
}

func efirst() error {
	e := errors.NewStatus(ErrEOF, "error example")
	return e
}

func esecond() error {
	e := efirst()
	if e != nil {
		return e
	}
	return nil
}

func ethird() error {
	e := esecond()
	if e != nil {
		return errors.NewStatus(ErrCall, e)
	}
	return nil
}

func TestErrorStack2(t *testing.T) {
	Convey("check errors function line number", t, func() {
		// status
		Convey("WithStack error", func() {

			//模拟跨服务错误转换
			e2 := ethird()
			s3 := errors.ToStatus(e2)
			json.PrintObject(s3)
		})
	})
}

//func TestNewDesc(t *testing.T) {
//	Convey("NewDesc", t, func() {
//		Convey("unknown code", func() {
//			e := errors.NewDesc(8888888888888, "not exist code")
//			So(e.Code(), ShouldEqual, errors.Unknown().Code())
//			So(e.HTTPStatus(), ShouldEqual, errors.Unknown().HTTPStatus())
//			So(e.Message(), ShouldEqual, errors.Unknown().Message())
//			So(e.Description(), ShouldEqual, "not exist code")
//		})
//
//		Convey("exist", func() {
//			e := errors.NewDesc(101, "NOT")
//			So(e.Code(), ShouldEqual, e.Code())
//			So(e.HTTPStatus(), ShouldEqual, e.HTTPStatus())
//			So(e.Message(), ShouldEqual, e.Message())
//			So(e.Description(), ShouldEqual, "NOT")
//			So(e.Stack(), ShouldNotEqual, "")
//		})
//	})
//}
//
////
////func TestNew(t *testing.T) {
////	Convey("New", t, func() {
////		Convey("unknown code", func() {
////			e := errors.New(8888888888888, fmt.Errorf("myError"))
////			So(e.Code(), ShouldEqual, errors.Unknown().Code())
////			So(e.HTTPStatus(), ShouldEqual, errors.Unknown().HTTPStatus())
////			So(e.Message(), ShouldEqual, errors.Unknown().Message())
////			So(e.Description(), ShouldEqual, "myError")
////		})
////
////		Convey("nil error", func() {
////			e := errors.New(101, nil)
////			So(e.Code(), ShouldEqual, e.Code())
////
////		})
////
////		Convey("exist", func() {
////			Convey("WithStack error", func() {
////				e1 := errors.NewDesc(100, "error1")
////				e2 := errors.New(101, e1)
////
////				So(e2.Code(), ShouldEqual, 101)
////				So(len(e2.StackInfo()), ShouldEqual, 2)
////			})
////
////			Convey("normal error", func() {
////				e := errors.New(101, fmt.Errorf("myError"))
////				So(e.Code(), ShouldEqual, 101)
////				So(e.HTTPStatus(), ShouldEqual, e.HTTPStatus())
////				So(e.Message(), ShouldEqual, e.Message())
////				So(e.Description(), ShouldEqual, "myError")
////				So(e.Stack(), ShouldNotEqual, "")
////			})
////		})
////	})
////}
////
////func TestFormat(t *testing.T) {
////	Convey("Format", t, func() {
////		Convey("%s", func() {
////			e := errors.NewDesc(101, "file not exist")
////			So(fmt.Sprintf("%s", e), ShouldEqual, "OpenFileError:file not exist")
////			So(fmt.Sprintf("%q", e), ShouldEqual, "\"OpenFileError:file not exist\"")
////			So(fmt.Sprintf("%v", e), ShouldEqual, "OpenFileError:file not exist")
////			//So(
////			//	fmt.Sprintf("%+v", e),
////			//	ShouldEqual,
////			// 	"OpenFileError:file not exist
////			// [host:127.0.0.1,pid:8536,module:testing,code:101,file:error_test.go,func:1,line:41]",
////			//)
////			So(
////				fmt.Sprintf("%#v", e),
////				ShouldEqual,
////				"{\"code\":101,\"desc\":\"file not exist\",\"message\":{\"MessageCN\":\"访问文件失败\",\"MessageEN\":\"OpenFileError\"}}",
////			)
////			//So(
////			//	fmt.Sprintf("%+#v", e),
////			//	ShouldEqual,
////			// 	"{\"code\":101,\"desc\":\"file not
////			// exist\",\"http\":200,\"message\":{\"cn\":\"访问文件失败\",\"en\":\"OpenFileError\"},\"stack\":[{\"host\":\"127.0.0.1\",\"pid\":\"8536\",\"module\":\"testing\",\"code\":\"101\",\"file_name\":\"error_test.go\",\"func_name\":\"1\",\"line\":\"41\"}]}",
////			//)
////		})
////	})
////}
////
////func TestFromError(t *testing.T) {
////	Convey("FormatError", t, func() {
////		Convey("error is WithStack error", func() {
////			e := errors.NewDesc(102, "some thing happen")
////			st := errors.FromError(e)
////			So(st, ShouldNotBeNil)
////			So(st.Code(), ShouldEqual, 102)
////			So(st.Error(), ShouldEqual, "ReadFileError:some thing happen")
////		})
////
////		Convey("error is simple error", func() {
////			e := fmt.Errorf("i'm not WithStack error")
////			st := errors.FromError(e)
////			So(st, ShouldNotBeNil)
////			So(st.Code(), ShouldEqual, errors.Unknown().Code())
////			So(st.Description(), ShouldEqual, "i'm not WithStack error")
////			So(st.Error(), ShouldEqual, "unknown error code:i'm not WithStack error")
////		})
////
////		Convey("error is nil", func() {
////			var e error
////			st := errors.FromError(e)
////			So(st, ShouldBeNil)
////		})
////
////	})
////}
//
//func TestUpdateStack(t *testing.T) {
//	Convey("UpdateStack", t, func() {
//		Convey("error is WithStack error", func() {
//			e := errors.NewDesc(102, "some thing happen")
//			st := errors.UpdateStack(e)
//			ss := errors.FromError(st)
//			So(ss, ShouldNotBeNil)
//			So(ss.Code(), ShouldEqual, 102)
//			So(ss.Error(), ShouldEqual, "ReadFileError:some thing happen")
//			So(len(ss.Stack()), ShouldEqual, 2)
//		})
//
//		Convey("error is simple error", func() {
//			e := fmt.Errorf("i'm not WithStack error")
//			st := errors.UpdateStack(e)
//			ss := errors.FromError(st)
//			So(ss, ShouldNotBeNil)
//			So(ss.Code(), ShouldEqual, errors.Unknown().Code())
//			So(ss.Description(), ShouldEqual, "i'm not WithStack error")
//			So(ss.Error(), ShouldEqual, "unknown error code:i'm not WithStack error")
//			So(len(ss.Stack()), ShouldEqual, 2)
//		})
//
//		Convey("when error is nil", func() {
//			e := errors.UpdateStack(nil)
//			So(e, ShouldBeNil)
//			IsError := func(err error) (ok bool) {
//				if err == nil {
//					return true
//				}
//				return false
//			}
//			So(IsError(e), ShouldBeTrue)
//		})
//	})
//}
//
//func TestIsCode(t *testing.T) {
//	Convey("isCode", t, func() {
//		e := errors.NewDesc(101, "error 1001")
//		So(errors.IsCode(e, 100), ShouldBeFalse)
//		So(errors.IsCode(e, 101), ShouldBeTrue)
//		So(errors.IsCode(e, 101222), ShouldBeFalse)
//	})
//}
//
//type fakeModule struct {
//	name string
//	ip   string
//	pid  int
//}
//
//func (s fakeModule) PID() int {
//	return s.pid
//}
//
//func (s fakeModule) IP() string {
//	return s.ip
//}
//
//func (s fakeModule) Name() string {
//	return s.name
//}
//
//func (s fakeModule) String() string {
//	return fmt.Sprintf("host:%s,pid:%d,module:%s", s.IP(), s.PID(), s.Name())
//}
//
//func init() {
//	errors.UpdateModuleInfo(&fakeModule{
//		name: "testing",
//		ip:   "127.0.0.1",
//		pid:  8536,
//	})
//	errors.MustRegister(errors.NewCoder(100, 200, map[string]string{
//		errors.MessageLangENKey: "WriteFileError",
//		errors.MessageLangCNKey: "写文件失败",
//	}))
//	errors.MustRegister(errors.NewCoder(101, 200, map[string]string{
//		errors.MessageLangENKey: "OpenFileError",
//		errors.MessageLangCNKey: "访问文件失败",
//	}))
//	errors.MustRegister(errors.NewCoder(102, 200, map[string]string{
//		errors.MessageLangENKey: "ReadFileError",
//		errors.MessageLangCNKey: "读文件失败",
//	}))
//}
//
//func TestFormatPrint(t *testing.T) {
//	Convey("Format", t, func() {
//		Convey("%s", func() {
//			//e := httperrors.NewDesc(101, "file not exist")
//			//fmt.Println(e.Error())
//			//fmt.Printf("%v\n", e)
//			//fmt.Printf("%+v\n", e)
//			//fmt.Printf("%+#v\n", e)
//			e2 := fmt.Errorf("my error")
//			e2 = errors.UpdateStack(e2)
//			e2 = errors.UpdateStack(e2)
//			e2 = errors.UpdateStack(e2)
//			e2 = errors.UpdateStack(e2)
//			fmt.Println(e2.Error())
//			fmt.Printf("%v\n", e2)
//			fmt.Printf("%+v\n", e2)
//			fmt.Printf("%+#v\n", e2)
//		})
//	})
//}
