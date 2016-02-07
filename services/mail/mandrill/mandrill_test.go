package mandrill

import (
	"testing"

	. "github.com/freeusd/solebtc/Godeps/_workspace/src/github.com/smartystreets/goconvey/convey"
	"github.com/freeusd/solebtc/errors"
)

func TestMandrill(t *testing.T) {
	Convey("Given valid mandrill client", t, func() {
		m := New("SANDBOX_SUCCESS", "from email", "from name")

		Convey("When send email", func() {
			err := m.SendEmail([]string{"email@address.com"}, "subject", "html")

			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given invalid mandrill client", t, func() {
		m := New("SANDBOX_ERROR", "from email", "from name")

		Convey("When send email", func() {
			err := m.SendEmail([]string{"email@address.com"}, "subject", "html")

			Convey("Error should be mandrill error", func() {
				So(err.ErrCode, ShouldEqual, errors.ErrCodeMandrill)
			})
		})
	})
}
