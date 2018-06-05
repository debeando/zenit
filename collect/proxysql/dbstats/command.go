// Package collect all stats_mysql_commands_counters stats from ProxySQL.
package dbstats

import (
  "fmt"
)

func Command() {
  fmt.Printf("--> proxysql.stats.commands.\n")
}
