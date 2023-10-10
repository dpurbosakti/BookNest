package errorhelper

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// add logger with error level and what it caused to failed to get detail
func FailedGetDetail(logger *logrus.Entry, err error, data string) {
	logger.WithError(err).Error(fmt.Sprintf("failed to get detail %s", data))
}

// add logger with error level and what it caused to failed to update
func FailedUpdate(logger *logrus.Entry, err error, data string) {
	logger.WithError(err).Error(fmt.Sprintf("failed to update %s", data))
}

// add logger with error level and what it caused to failed to delete
func FailedDelete(logger *logrus.Entry, err error, data string) {
	logger.WithError(err).Error(fmt.Sprintf("failed to delete %s", data))
}

// add logger with error level and what it caused to failed to get list
func FailedGetList(logger *logrus.Entry, err error, data string) {
	logger.WithError(err).Error(fmt.Sprintf("failed to get list %s", data))
}

// add logger with error level and what it caused to failed to create
func FailedCreate(logger *logrus.Entry, err error, data string) {
	logger.WithError(err).Error(fmt.Sprintf("failed to create %s", data))
}
