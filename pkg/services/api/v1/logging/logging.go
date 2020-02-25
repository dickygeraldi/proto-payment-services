package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

func SetLogging(api, dataIp, cases, fatal, level, message, time string) bool {
	var log = logrus.New()
	log.Out = os.Stdout

	switch level {
	case "info":
		log.WithFields(logrus.Fields{
			"api":   api,
			"cases": cases,
			"fatal": fatal,
			"ip":    dataIp,
			"time":  time,
		}).Info(message)
	case "debug":
		log.WithFields(logrus.Fields{
			"api":   api,
			"cases": cases,
			"fatal": fatal,
			"ip":    dataIp,
			"time":  time,
		}).Debug(message)
	case "panic":
		log.WithFields(logrus.Fields{
			"api":   api,
			"cases": cases,
			"fatal": fatal,
			"ip":    dataIp,
			"time":  time,
		}).Panic(message)
	case "warning":
		log.WithFields(logrus.Fields{
			"api":   api,
			"cases": cases,
			"fatal": fatal,
			"ip":    dataIp,
			"time":  time,
		}).Warn(message)
	case "fatal":
		log.WithFields(logrus.Fields{
			"api":   api,
			"cases": cases,
			"fatal": fatal,
			"ip":    dataIp,
			"time":  time,
		}).Fatal(message)

	}

	return true
}
