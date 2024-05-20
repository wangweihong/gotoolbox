package httperrors_test

import (
	"fmt"
	"testing"

	errors "github.com/wangweihong/gotoolbox/errors/httperrors"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAggregateError(t *testing.T) {
	Convey("Aggregates", t, func() {
		Convey("e", func() {
			e1 := errors.NewDesc(101, "error1")
			e2 := errors.NewDesc(101, "error2")

			fmt.Println(errors.NewAggregate(e1, e2).Error())
		})
	})
}
