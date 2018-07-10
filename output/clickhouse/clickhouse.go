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
                              e["IP"])

      common.HTTPPost(config.DSN_CLICKHOUSE, sql)
    }
  }()
}

func SendMySQLSlowLog(event <-chan map[string]string) {
  go func() {
    for e := range event {
      sql := fmt.Sprintf("INSERT INTO zenit.mysql_slow_log " +
                             "(_time,bytes_sent,killed,last_errno,lock_time,query,query_time,rows_affected,rows_examined,rows_read,rows_sent,schema,thread_id,user_host) " +
                             "VALUES ('%s',%s,%s,%s,%s,'%s',%s,%s,%s,%s,%s,'%s',%s,'%s')",
                              common.ISO8601V2toRFC3339(e["time"]),
                              e["bytes_sent"],
                              e["killed"],
                              e["last_errno"],
                              e["lock_time"],
                              common.Escape(e["query"]),
                              e["query_time"],
                              e["rows_affected"],
                              e["rows_examined"],
                              e["rows_read"],
                              e["rows_sent"],
                              e["schema"],
                              e["thread_id"],
                              e["user_host"])

      common.HTTPPost(config.DSN_CLICKHOUSE, sql)
    }
  }()
}
