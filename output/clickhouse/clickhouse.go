package clickhouse

import (
  "fmt"
  "gitlab.com/swapbyt3s/zenit/common"
  "gitlab.com/swapbyt3s/zenit/config"
)

var ip_address string

func init() {
  ip_address = common.IpAddress()
}

func SendMySQLAuditLog(event <-chan map[string]string) {
  go func() {
    for e := range event {
      sql := fmt.Sprintf("INSERT INTO zenit.mysql_audit_log " +
                         "(_time,host_ip,name,command_class,connection_id,status,sqltext,user,host,os_user,ip) " +
                         "VALUES ('%s',IPv4StringToNum('%s'),'%s','%s',%s,%s,'%s','%s','%s','%s','%s')",
                          common.ToDateTime(e["timestamp"]),
                          ip_address,
                          e["name"],
                          e["command_class"],
                          e["connection_id"],
                          e["status"],
                          common.Escape(e["sqltext"]),
                          e["user"],
                          e["host"],
                          e["os_user"],
                          e["ip"])

      common.HTTPPost(config.DSN_CLICKHOUSE, sql)
    }
  }()
}

func SendMySQLSlowLog(event <-chan map[string]string) {
  go func() {
    for e := range event {
      sql := fmt.Sprintf("INSERT INTO zenit.mysql_slow_log " +
                         "(_time,host_ip,bytes_sent,killed,last_errno,lock_time,query,query_time,rows_affected,rows_examined,rows_read,rows_sent,schema,thread_id,user_host) " +
                         "VALUES (toDateTime(%s),IPv4StringToNum('%s'),%s,%s,%s,%s,'%s',%s,%s,%s,%s,%s,'%s',%s,'%s')",
                          e["timestamp"],
                          ip_address,
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
