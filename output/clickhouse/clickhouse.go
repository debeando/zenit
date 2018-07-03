package clickhouse

import (
  "fmt"
  "gitlab.com/swapbyt3s/zenit/common"
  "gitlab.com/swapbyt3s/zenit/config"
)

func SendMySQLAuditLog(event <-chan map[string]string) {
  go func() {
    for e := range event {
      sql := fmt.Sprintf("INSERT INTO zenit.mysql_audit_log " +
                             "(_time,name,command_class,connection_id,status,sqltext,user,host,os_user,ip) " +
                             "VALUES ('%s','%s','%s',%s,%s,'%s','%s','%s','%s','%s')",
                              common.ToDateTime(e["TIMESTAMP"]),
                              e["NAME"],
                              e["COMMAND_CLASS"],
                              e["CONNECTION_ID"],
                              e["STATUS"],
                              common.Escape(e["SQLTEXT"]),
                              e["USER"],
                              e["HOST"],
                              e["OS_USER"],
                              e["IP"],
          )

      common.HTTPPost(config.DSN_CLICKHOUSE, sql)
    }
  }()
}
