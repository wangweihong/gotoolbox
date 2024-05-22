package log_test

import (
	"context"
	"testing"

	"github.com/wangweihong/gotoolbox/src/log"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWithFieldPair(t *testing.T) {
	defer log.Flush()

	var ctx = context.Background()
	Convey("TestWithFieldPair", t, func() {
		child := log.WithFieldPair(ctx, "a", "B")
		child2 := log.WithFieldPair(ctx, "a", "C")
		log.F(child).Info("")
		log.F(child2).Info("")

	})
}

func TestWithFields(t *testing.T) {
	defer log.Flush()

	var ctx = context.Background()
	Convey("TestWithFields", t, func() {
		log.F(ctx).Info("raw")
		f := make(map[string]interface{})
		f["a"] = "b"

		child := log.WithFields(ctx, f)
		child2 := log.WithFields(ctx, f)

		f1 := child.Value(log.FieldKeyCtx{})
		f2 := child2.Value(log.FieldKeyCtx{})

		So(f1, ShouldResemble, f2)
		f["c"] = "d"

		f1 = child.Value(log.FieldKeyCtx{})
		f2 = child2.Value(log.FieldKeyCtx{})

		So(f1, ShouldResemble, f2)

		f1m := f1.(map[string]interface{})
		So(f1m["c"], ShouldBeNil)

		child3 := log.WithFields(ctx, f)
		fm2 := make(map[string]interface{})
		fm2["a"] = "c"
		child4 := log.WithFields(ctx, fm2)

		log.F(child3).Info("")
		log.F(child4).Info("")

	})
}
