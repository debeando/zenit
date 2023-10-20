package config

import (
	"time"
)

// All is a struct to contain all configuration imported or loaded from config file.
type Config struct {
	Path      string
	IPAddress string
	General   struct {
		Hostname  string        `yaml:"hostname"`
		Interval  time.Duration `yaml:"interval"`
		AWSRegion string        `yaml:"aws_region"`
	}
	Inputs struct {
		AWS struct {
			Discover struct {
				Enable   bool   `yaml:"enable"`
				Username string `yaml:"username"`
				Password string `yaml:"password"`
				Filter   string `yaml:"filter"`
				Plugins  struct {
					MySQL struct {
						Aurora    bool `yaml:"aurora"`
						Enable    bool `yaml:"enable"`
						InnoDB    bool `yaml:"innodb"`
						Overflow  bool `yaml:"overflow"`
						Replica   bool `yaml:"replica"`
						Status    bool `yaml:"status"`
						Tables    bool `yaml:"tables"`
						Variables bool `yaml:"variables"`
					}
				}
			}
			CloudWatch struct {
				Enable bool `yaml:"enable"`
			}
		}
		MongoDB  []MongoDB
		MySQL    []MySQL
		ProxySQL []struct {
			Hostname string `yaml:"hostname"`
			DSN      string `yaml:"dsn"`
			Enable   bool   `yaml:"enable"`
			Commands bool   `yaml:"commands"`
			Errors   bool   `yaml:"errors"`
			Global   bool   `yaml:"global"`
			Pool     bool   `yaml:"pool"`
			Queries  bool   `yaml:"queries"`
		}
		OS struct {
			CPU    bool `yaml:"cpu"`
			Disk   bool `yaml:"disk"`
			Limits bool `yaml:"limits"`
			Mem    bool `yaml:"mem"`
			Net    bool `yaml:"net"`
		}
		Process struct {
			PerconaToolKitDeadlockLogger     bool `yaml:"pt_deadlock_logger"`
			PerconaToolKitKill               bool `yaml:"pt_kill"`
			PerconaToolKitOnlineSchemaChange bool `yaml:"pt_online_schema_change"`
			PerconaToolKitSlaveDelay         bool `yaml:"pt_slave_delay"`
			PerconaXtraBackup                bool `yaml:"xtrabackup"`
		}
	}
	Outputs struct {
		InfluxDB struct {
			Enable   bool   `yaml:"enable"`
			URL      string `yaml:"url"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
			Database string `yaml:"database"`
		}
	}
}

type MySQL struct {
	Engine    string
	Hostname  string `yaml:"hostname"`
	DSN       string `yaml:"dsn"`
	Aurora    bool   `yaml:"aurora"`
	Enable    bool   `yaml:"enable"`
	InnoDB    bool   `yaml:"innodb"`
	Overflow  bool   `yaml:"overflow"`
	Replica   bool   `yaml:"replica"`
	Status    bool   `yaml:"status"`
	Tables    bool   `yaml:"tables"`
	Variables bool   `yaml:"variables"`
}

type MongoDB struct {
	Hostname     string `yaml:"hostname"`
	DSN          string `yaml:"dsn"`
	Enable       bool   `yaml:"enable"`
	ServerStatus bool   `yaml:"serverstatus"`
	Collections  bool   `yaml:"collections"`
}
