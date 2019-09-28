package audit

import (
	"fmt"
	"sync"

	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/inputs"
	"github.com/swapbyt3s/zenit/plugins/inputs/parsers/mysql/audit/xml-old"
	"github.com/swapbyt3s/zenit/plugins/outputs/clickhouse"
)

type MySQLAuditLog struct{}

var (
	instance *MySQLAuditLog
	once     sync.Once
)

func (l *MySQLAuditLog) Collect() {
	once.Do(func() {
		if instance == nil {
			instance = &MySQLAuditLog{}

			l.Parser()
		}
	})
}

func (l *MySQLAuditLog) Parser() {
	defer func() {
		if err := recover(); err != nil {
			log.Debug(fmt.Sprintf("Plugin - InputMySQLAuditLog - Panic (code %d) has been recover from somewhere.\n", err))
		}
	}()

	if config.File.Parser.MySQL.AuditLog.Enable {
		if config.File.General.Debug {
			log.Debug("Load MySQL AuditLog")
			log.Debug("Read MySQL AuditLog: " + config.File.Parser.MySQL.AuditLog.LogPath)
		}

		if !clickhouse.Check() {
			log.Error("AuditLog require active connection to ClickHouse.")
		}

		if config.File.Parser.MySQL.AuditLog.Format == "xml-old" {
			log.Info("Plugin - InputMySQLAuditLog - Load xml-old parser")
			xmlold.Collect()
		}
	}
}

func init() {
	inputs.Add("InputMySQLAuditLog", func() inputs.Input { return &MySQLAuditLog{} })
}
