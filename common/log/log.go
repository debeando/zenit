package log

import (
	"flag"
	"io/ioutil"

	"github.com/swapbyt3s/zenit/config"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)

	if flag.Lookup("test.v") != nil {
		logrus.SetOutput(ioutil.Discard)
	}
}

func Info(m string, f map[string]interface{}) {
	logrus.WithFields(f).Info(m)
}

func Error(m string, f map[string]interface{}) {
	logrus.WithFields(f).Error(m)
}

func Debug(m string, f map[string]interface{}) {
	if flag.Lookup("debug") != nil {
		config.File.General.Debug = flag.Lookup("debug").Value.(flag.Getter).Get().(bool)
	}

	if config.File.General.Debug {
		logrus.WithFields(f).Debug(m)
	}
}
