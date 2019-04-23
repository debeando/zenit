package config

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/go-yaml/yaml"

	"github.com/swapbyt3s/zenit/common"
)

var File Config

// Init does any initialization necessary for the module.
func init() {
	File = Config {
		Path: "/etc/zenit/zenit.yaml",
		IPAddress: common.IPAddress(),
	}

	log.SetOutput(os.Stdout)
}

func (c *Config) Load() {
	source, err := ioutil.ReadFile(c.Path)
	if err != nil {
		source, err = ioutil.ReadFile("zenit.yaml")
		if err != nil {
			log.Printf("Fail to read config file: %s or %s", c.Path, "./zenit.yaml")
			os.Exit(1)
		}
	}

	source = []byte(os.ExpandEnv(string(source)))

	err = yaml.Unmarshal(source, &c)
	if err != nil {
		log.Printf("Imposible to parse config file - %s", err)
		os.Exit(1)
	}
}

// SanityCheck verify the minimum config settings and set default values to start.
func (c *Config) SanityCheck() {
	if c.General.Interval < 3 {
		log.Println("W! Config - general.interval: Use positive value, and minimun start from 3 seconds.")
		log.Println("W! Config - general.interval: Using default 30 seconds.")
		c.General.Interval = 30
	}

	if len(c.General.Hostname) == 0 {
		log.Println("W! Config - general.hostname: Custom value is not set, using current.")
		c.General.Hostname = common.Hostname()
	}
}
