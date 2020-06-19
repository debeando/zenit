package log

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/swapbyt3s/zenit/config"
	"github.com/sirupsen/logrus"
)

var debug bool

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	if flag.Lookup("test.v") != nil {
		logrus.SetOutput(ioutil.Discard)
	} else {
		logrus.SetOutput(os.Stdout)
	}

	logrus.SetLevel(logrus.InfoLevel)

	if flag.Lookup("debug") == nil {
		flag.BoolVar(&debug, "debug", false, "Enable debug mode.")
	}
	flag.Parse()

	if flag.Lookup("debug").Value.(flag.Getter).Get().(bool) {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if config.File.General.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func Info(m string, f map[string]interface{}) {
	logrus.WithFields(f).Info(m)
}

func Warning(m string, f map[string]interface{}) {
	logrus.WithFields(f).Warning(m)
}

func Error(m string, f map[string]interface{}) {
	logrus.WithFields(f).Error(m)
}

func Debug(m string, f map[string]interface{}) {
	logrus.WithFields(f).Debug(m)
}
