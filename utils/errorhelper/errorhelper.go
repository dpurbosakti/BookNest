package errorhelper

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// add logger with error level with what it caused
func FailedGetDetail(logger *logrus.Entry, err error, data string) {
	logger.WithError(err).Error(fmt.Sprintf("failed to get detail %s", data))
}
