package mysql

import (
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
)

func Check() bool {
	var enable = 0

	if ( config.File.MySQL.Inputs.Indexes         ) { enable++ }
	if ( config.File.MySQL.Inputs.Overflow        ) { enable++ }
	if ( config.File.MySQL.Inputs.Slave           ) { enable++ }
	if ( config.File.MySQL.Inputs.Status          ) { enable++ }
	if ( config.File.MySQL.Inputs.Tables          ) { enable++ }
	if ( config.File.MySQL.Inputs.Variables       ) { enable++ }
	if ( config.File.MySQL.Inputs.AuditLog.Enable ) { enable++ }
	if ( config.File.MySQL.Inputs.SlowLog.Enable  ) { enable++ }

	if enable > 0 {
		return mysql.Check(config.File.MySQL.DSN, "MySQL")
	}

	return false
}
