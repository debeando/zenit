package audit

import (
	"sync"

	"github.com/debeando/zenit/common/log"
	"github.com/debeando/zenit/config"
	"github.com/debeando/zenit/plugins/inputs"
	"github.com/debeando/zenit/plugins/inputs/parsers/mysql/audit/xml-old"
	"github.com/debeando/zenit/plugins/outputs/clickhouse"
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
			log.Error("MySQLAuditLog", map[string]interface{}{"error": err})
		}
	}()

	if config.File.Parser.MySQL.AuditLog.Enable {
		log.Debug("MySQLAuditLog", map[string]interface{}{"slow_log_path": config.File.Parser.MySQL.AuditLog.LogPath})

		if !clickhouse.Check() {
			log.Error("MySQLAuditLog", map[string]interface{}{"error": "AuditLog require active connection to ClickHouse."})
		}

		if config.File.Parser.MySQL.AuditLog.Format == "xml-old" {
			xmlold.Collect()
		}
	}
}

func init() {
	inputs.Add("InputMySQLAuditLog", func() inputs.Input { return &MySQLAuditLog{} })
}
