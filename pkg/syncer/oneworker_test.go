package syncer_test

import (
	"testing"
	"time"

	"github.com/wangweihong/gotoolbox/pkg/syncer"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewOneWorkerSyncer(t *testing.T) {
	Convey("onework", t, func() {
		SkipConvey("cronjob syncer", func() {
			stop := make(chan struct{}, 0)
			s := syncer.NewOneWorkerSyncer(func(arg interface{}) error {
				time.Sleep(1 * time.Second)
				return nil
			}, 3*time.Second, 1)
			s.Run(stop)
			So(s.Trigger(false), ShouldBeTrue)

			go func() {
				select {
				case <-time.After(3 * time.Second):
					close(stop)
				}
			}()
			<-stop
		})

		Convey("trigger syncer", func() {
			s := syncer.NewOneWorkerSyncer(func(arg interface{}) error {
				time.Sleep(100 * time.Millisecond)
				return nil
			}, 3*time.Second, 3)
			So(s.Trigger(false), ShouldBeFalse)
			time.Sleep(300 * time.Millisecond)
			So(s.Trigger(false), ShouldBeFalse)
			time.Sleep(300 * time.Millisecond)
			So(s.Trigger(false), ShouldBeFalse)
			time.Sleep(300 * time.Millisecond)
			So(s.Trigger(false), ShouldBeFalse)
			time.Sleep(300 * time.Millisecond)

			So(len(s.GetRecords()), ShouldEqual, 3)
		})
	})
}
