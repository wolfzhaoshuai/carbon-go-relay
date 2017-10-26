package global

import . "github.com/smartystreets/goconvey/convey"
import "testing"

func TestConfig(t *testing.T) {
	ParseConfig("../cfg.json")
	config := Config()

	Convey("should be ok", t, func() {
		So(config.SendBatchSize, ShouldEqual, 2000)
	})
}
