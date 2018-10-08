package log

import (
	"flag"
	"log"
	"strings"

	"github.com/swapbyt3s/zenit/config"
)

func Info(m string) {
	if flag.Lookup("test.v") == nil {
		log.Printf("I! - %s\n", m)
	}
}

func Error(m string) {
	if flag.Lookup("test.v") == nil {
		log.Printf("E! - %s\n", m)
	}
}

func Debug(m string) {
	if config.File.General.Debug {
		if flag.Lookup("test.v") == nil {
			m = strings.Replace(m, "\n", ",", 0)
			log.Printf("D! - %s\n", m)
		}
	}
}
