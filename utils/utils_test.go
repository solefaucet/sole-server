package utils

import (
	"io/ioutil"

	"github.com/Sirupsen/logrus"
)

func init() {
	logrus.SetOutput(ioutil.Discard)
}
