package config

import (
	"time"
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
				Enable   bool   `yaml:"enable"`
				Warning  uint64 `yaml:"warning"`
				Critical uint64 `yaml:"critical"`
				Duration int    `yaml:"duration"`
			}
			Replication struct {
				Enable   bool   `yaml:"enable"`
				Warning  uint64 `yaml:"warning"`
				Critical uint64 `yaml:"critical"`
				Duration int    `yaml:"duration"`
			}
		}
	}
	ProxySQL []struct {
		Hostname string `yaml:"hostname"`
		DSN      string `yaml:"dsn"`
		Inputs struct {
			Commands bool `yaml:"commands"`
			Pool     bool `yaml:"pool"`
			Queries  bool `yaml:"queries"`
		}
		Alerts struct {
			Errors struct {
				Enable   bool   `yaml:"enable"`
				Warning  uint64 `yaml:"warning"`
				Critical uint64 `yaml:"critical"`
				Duration int    `yaml:"duration"`
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
				Enable   bool   `yaml:"enable"`
				Warning  uint64 `yaml:"warning"`
				Critical uint64 `yaml:"critical"`
				Duration int    `yaml:"duration"`
			}
			Disk struct {
				Enable   bool   `yaml:"enable"`
				Warning  uint64 `yaml:"warning"`
				Critical uint64 `yaml:"critical"`
				Duration int    `yaml:"duration"`
			}
			MEM struct {
				Enable   bool   `yaml:"enable"`
				Warning  uint64 `yaml:"warning"`
				Critical uint64 `yaml:"critical"`
				Duration int    `yaml:"duration"`
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
