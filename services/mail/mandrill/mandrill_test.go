package mandrill

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
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

			Convey("Error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}
