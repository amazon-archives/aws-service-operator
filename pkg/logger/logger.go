package logger

import (
	"bufio"
	"os"
	"strings"

	"github.com/awslabs/aws-service-operator/pkg/config"
	"github.com/sirupsen/logrus"
)

// Configure will setup the logger for all commands
func Configure(config config.LoggingConfig) (*logrus.Entry, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	// use a file if you want
	if config.File != "" {
		f, errOpen := os.OpenFile(config.File, os.O_RDWR|os.O_APPEND, 0660)
		if errOpen != nil {
			return nil, errOpen
		}
		logrus.SetOutput(bufio.NewWriter(f))
	}

	if config.Level != "" {
		level, err := logrus.ParseLevel(strings.ToUpper(config.Level))
		if err != nil {
			return nil, err
		}
		logrus.SetLevel(level)
	}

	// always use the fulltimestamp
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    config.FullTimestamps,
		DisableTimestamp: config.DisableTimestamps,
	})

	return logrus.StandardLogger().WithField("hostname", hostname), nil
}
