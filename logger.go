package gcpinstancesinfo

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	// Set the output to Standard Err
	log.Out = os.Stderr

	// Set the log level
	log.SetLevel(logrus.ErrorLevel)

	// Set the log format
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}
