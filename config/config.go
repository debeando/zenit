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
	// Author is a const with have the name of creator and collaborators for this code.
	Author string = "Nicola Strappazzon C. <nstrappazzonc@gmail.com>"
	// Version is a const to have the latest version number for this code.
	Version string = "1.2.0"
	// ConfigFile containt the path for configuration file.
	ConfigFile string = "/etc/zenit/zenit.yaml"
)

// All is a struct to contain all configuration imported or loaded from config file.
type All struct {
	General struct {
		Hostname string        `yaml:"hostname"`
		Interval time.Duration `yaml:"interval"`
		Debug    bool          `yaml:"debug"`
	}
	MySQL struct {
		DSN string `yaml:"dsn"`
		Inputs struct {
			Indexes   bool `yaml:"indexes"`
			Overflow  bool `yaml:"overflow"`
			Slave     bool `yaml:"slave"`
			Status    bool `yaml:"status"`
			Tables    bool `yaml:"tables"`
			Variables bool `yaml:"variables"`
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
		Alerts struct {
			ReadOnly struct {
				Enable   bool `yaml:"enable"`
				Duration int  `yaml:"duration"`
			}
			Connections struct {
				Enable   bool `yaml:"enable"`
				Warning  int  `yaml:"warning"`
				Critical int  `yaml:"critical"`
				Duration int  `yaml:"duration"`
			}
			Replication struct {
				Enable   bool `yaml:"enable"`
				Warning  int  `yaml:"warning"`
				Critical int  `yaml:"critical"`
				Duration int  `yaml:"duration"`
			}
		}
	}
	ProxySQL struct {
		DSN string `yaml:"dsn"`
		Inputs struct {
			Commands bool `yaml:"commands"`
			Pool     bool `yaml:"pool"`
			Queries  bool `yaml:"queries"`
		}
		Alerts struct {
			Errors struct {
				Enable   bool `yaml:"enable"`
				Warning  int  `yaml:"warning"`
				Critical int  `yaml:"critical"`
				Duration int  `yaml:"duration"`
			}
			Status struct {
				Enable   bool `yaml:"enable"`
				Duration int  `yaml:"duration"`
			}
		}
	}
	ClickHouse struct {
		DSN string `yaml:"dsn"`
	}
	Prometheus struct {
		Enable   bool   `yaml:"enable"`
		TextFile string `yaml:"textfile"`
	}
	Slack struct {
		Enable  bool   `yaml:"enable"`
		Token   string `yaml:"token"`
		Channel string `yaml:"channel"`
	}
	OS struct {
		Inputs struct {
			CPU    bool `yaml:"cpu"`
			Disk   bool `yaml:"disk"`
			Limits bool `yaml:"limits"`
			Mem    bool `yaml:"mem"`
			Net    bool `yaml:"net"`
		}
		Alerts struct {
			CPU struct {
				Enable   bool `yaml:"enable"`
				Warning  int  `yaml:"warning"`
				Critical int  `yaml:"critical"`
				Duration int  `yaml:"duration"`
			}
			Disk struct {
				Enable   bool `yaml:"enable"`
				Warning  int  `yaml:"warning"`
				Critical int  `yaml:"critical"`
				Duration int  `yaml:"duration"`
			}
			MEM struct {
				Enable   bool `yaml:"enable"`
				Warning  int  `yaml:"warning"`
				Critical int  `yaml:"critical"`
				Duration int  `yaml:"duration"`
			}
		}
	}
	Process struct {
		Inputs struct {
			PerconaToolKitDeadlockLogger     bool `yaml:"pt_deadlock_logger"`
			PerconaToolKitKill               bool `yaml:"pt_kill"`
			PerconaToolKitOnlineSchemaChange bool `yaml:"pt_online_schema_change"`
			PerconaToolKitSlaveDelay         bool `yaml:"pt_slave_delay"`
			PerconaXtraBackup                bool `yaml:"xtrabackup"`
		}
	}
}

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
