package config

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/go-yaml/yaml"

	"github.com/swapbyt3s/zenit/common"
)

// ConfigFile containt the path for configuration file.
const ConfigFile string = "/etc/zenit/zenit.yaml"

// Define default variables and initialize structs.
var (
	IPAddress string
	File      All
)

// Init does any initialization necessary for the module.
func init() {
	IPAddress = common.IPAddress()

	log.SetOutput(os.Stdout)
}

// Load read settings from config file and set into struct.
func Load() {
	source, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		source, err = ioutil.ReadFile("zenit.yaml")
		if err != nil {
			log.Printf("Fail to read config file: %s or %s", ConfigFile, "./zenit.yaml")
			os.Exit(1)
		}
	}

	source = []byte(os.ExpandEnv(string(source)))

	err = yaml.Unmarshal(source, &File)
	if err != nil {
		log.Printf("Imposible to parse config file - %s", err)
		os.Exit(1)
	}
}

// SanityCheck verify the minimum config settings and set default values to start.
func SanityCheck() {
	if File.General.Interval < 3 {
		log.Println("W! Config - general.interval: Use positive value, and minimun start from 3 seconds.")
		log.Println("W! Config - general.interval: Using default 30 seconds.")
		File.General.Interval = 30
	}

	if len(File.General.Hostname) == 0 {
		log.Println("W! Config - general.hostname: Custom value is not set, using current.")
		File.General.Hostname = common.Hostname()
	}
}
