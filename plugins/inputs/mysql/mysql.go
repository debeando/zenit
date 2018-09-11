package mysql

import (
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
)

func Check() bool {
	return mysql.Check(config.File.MySQL.DSN, "MySQL")
}
