package mysql

import (
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
)

func Check() bool {
        // @TODO: Verify in all config in JSON to find values with any ENABLE
        // attribute to continue for this check.

	return mysql.Check(config.File.MySQL.DSN, "MySQL")
}
