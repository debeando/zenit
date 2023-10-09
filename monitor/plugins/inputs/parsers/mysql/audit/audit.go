package audit

import (
	"sync"

	"zenit/config"
	"zenit/monitor/plugins/inputs"
	"zenit/monitor/plugins/inputs/parsers/mysql/audit/xml-old"
	"zenit/monitor/plugins/lists/metrics"
	"zenit/monitor/plugins/outputs/clickhouse"

	"github.com/debeando/go-common/log"
)

type Plugin struct{}

var (
	instance *Plugin
	once     sync.Once
)

func (p *Plugin) Collect(name string, cnf *config.Config, mtc *metrics.Items) {
	once.Do(func() {
		if instance == nil {
			instance = &Plugin{}

			p.Parser(name, cnf)
		}
	})
}

func (p *Plugin) Parser(name string, cnf *config.Config) {
	defer func() {
		if err := recover(); err != nil {
			log.ErrorWithFields(name, log.Fields{"error": err})
		}
	}()

	if cnf.Parser.MySQL.AuditLog.Enable {
		log.DebugWithFields(name, log.Fields{"slow_log_path": cnf.Parser.MySQL.AuditLog.LogPath})

		if !clickhouse.Check() {
			log.ErrorWithFields(name, log.Fields{"error": "AuditLog require active connection to ClickHouse."})
		}

		if cnf.Parser.MySQL.AuditLog.Format == "xml-old" {
			xmlold.Collect()
		}
	}
}

func init() {
	inputs.Add("InputMySQLAuditLog", func() inputs.Input { return &Plugin{} })
}
