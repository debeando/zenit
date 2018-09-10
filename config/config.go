package config

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/go-yaml/yaml"

	"github.com/swapbyt3s/zenit/common"
)

const (
	AUTHOR  string = "Nicola Strappazzon C. <swapbyt3s@gmail.com>"
	VERSION string = "1.1.7"
)

var File All

type All struct {
	General struct {
		Hostname string        `yaml:"hostname"`
		Interval time.Duration `yaml:"interval"`
		Debug    bool          `yaml:"debug"`
	}
	MySQL struct {
		DSN       string `yaml:"dsn"`
		Overflow  bool   `yaml:"overflow"`
		Slave     bool   `yaml:"slave"`
		Status    bool   `yaml:"status"`
		Tables    bool   `yaml:"tables"`
		Variables bool   `yaml:"variables"`
		AuditLog struct {
			Enable        bool   `yaml:"enable"`
			Format        string `yaml:"format"`
			LogPath       string `yaml:"log_path"`
			BufferSize    int    `yaml:"buffer_size"`
			BufferTimeOut int    `yaml:"buffer_timeout"`
		}
		SlowLog struct {
			Enable        bool   `yaml:"enable"`
			LogPath       string `yaml:"log_path"`
			BufferSize    int    `yaml:"buffer_size"`
			BufferTimeOut int    `yaml:"buffer_timeout"`
		}
	}
	ProxySQL struct {
		DSN         string `yaml:"dsn"`
		QueryDigest bool   `yaml:"query_digest"`
	}
	ClickHouse struct {
		DSN string `yaml:"dsn"`
	}
	Prometheus struct {
		TextFile string `yaml:"textfile"`
	}
	OS struct {
		CPU    bool `yaml:"cpu"`
		Disk   bool `yaml:"disk"`
		Limits bool `yaml:"limits"`
		Mem    bool `yaml:"mem"`
		Net    bool `yaml:"net"`
	}
	Process struct {
		PerconaToolKitKill           bool `yaml:"pt_kill"`
		PerconaToolKitDeadlockLogger bool `yaml:"pt_deadlock_logger"`
		PerconaToolKitSlaveDelay     bool `yaml:"pt_slave_delay"`
	}
}

// Define default variables and initialize structs.
var (
	ConfigFile string = "/etc/zenit/zenit.yaml"
	IpAddress  string = ""
)

// Init does any initialization necessary for the module.
func init() {
	IpAddress = common.IpAddress()
}

// Loading settings from config file and set into struct.
func Load() {
	source, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		source, err = ioutil.ReadFile("zenit.yaml")
		if err != nil {
			log.Printf("Fail to read config file: %s or %s", ConfigFile, "./zenit.yaml")
			os.Exit(1)
		}
	}

	err = yaml.Unmarshal(source, &File)
	if err != nil {
		log.Printf("Imposible to parse config file - %s", err)
		os.Exit(1)
	}
}

// Check minimun config settings and set default values to start.
func SanityCheck() {
	if File.General.Interval < 5 {
		log.Println("W! Config - general.interval: Use positive value, and minimun start from 5 seconds.")
		log.Println("W! Config - general.interval: Using default 30 seconds.")
		File.General.Interval = 30
	}

	if len(File.General.Hostname) == 0 {
		log.Println("W! Config - general.hostname: Custom value is not set, using current.")
		File.General.Hostname = common.Hostname()
	}
}
