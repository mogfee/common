package logger

import (
	"github.com/sirupsen/logrus"
	"io"
)

var Logger *logrus.Logger

func New(w io.Writer) {
	Logger = logrus.New()
	Logger.SetOutput(w)
	Logger.SetReportCaller(true)
	Logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
}
