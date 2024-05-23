package errors_test

import (
	"fmt"
	"testing"

	errors2 "github.com/wangweihong/gotoolbox/pkg/errors"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAggregateError(t *testing.T) {
	Convey("Aggregates", t, func() {
		Convey("e", func() {
			e1 := errors2.NewDesc(101, "error1")
			e2 := errors2.NewDesc(101, "error2")

			fmt.Println(errors2.NewAggregate(e1, e2).Error())
		})
	})
}
