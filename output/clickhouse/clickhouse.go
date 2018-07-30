package clickhouse

import (
  "fmt"
  "strings"
  "gitlab.com/swapbyt3s/zenit/common"
  "gitlab.com/swapbyt3s/zenit/config"
)

func SendMySQLAuditLog(event <-chan map[string]string) {
  values := []string{}

  go func() {
    for e := range event {
      value := fmt.Sprintf("('%s',IPv4StringToNum('%s'),'%s','%s','%s',%s,%s,'%s','%s','%s','%s','%s','%s')",
                           e["timestamp"],
                           e["host_ip"],
                           e["host_name"],
                           e["name"],
                           e["command_class"],
                           e["connection_id"],
                           e["status"],
                           e["sqltext"],
                           e["sqltext_digest"],
                           e["user"],
                           e["host"],
                           e["os_user"],
                           e["ip"])

      values = append(values, value)

      if len(values) == 100 {
        sql := fmt.Sprintf("INSERT INTO zenit.mysql_audit_log " +
                           "(_time,host_ip,host_name,name,command_class,connection_id,status,sqltext,sqltext_digest,user,host,os_user,ip) " +
                           "VALUES %s;", strings.Join(values,","))
        // fmt.Printf("--(sql)> %s\n", sql)
        values = []string{}
        go common.HTTPPost(config.DSN_CLICKHOUSE, sql)
      }
    }
  }()
}

func SendMySQLSlowLog(event <-chan map[string]string) {
  values := []string{}

  go func() {
    for e := range event {
      value := fmt.Sprintf("(toDateTime(%s),IPv4StringToNum('%s'),'%s',%s,%s,%s,%s,'%s',%s,'%s',%s,%s,%s,%s,'%s',%s,'%s')",
        e["timestamp"],
        e["host_ip"],
        e["host_name"],
        e["bytes_sent"],
        e["killed"],
        e["last_errno"],
        e["lock_time"],
        e["query"],
        e["query_time"],
        e["query_digest"],
        e["rows_affected"],
        e["rows_examined"],
        e["rows_read"],
        e["rows_sent"],
        e["schema"],
        e["thread_id"],
        e["user_host"],
      )
      values = append(values, value)

      if len(values) == 100 {
        sql := fmt.Sprintf("INSERT INTO zenit.mysql_slow_log " +
                           "(_time,host_ip,host_name,bytes_sent,killed,last_errno,lock_time,query,query_time,query_digest,rows_affected,rows_examined,rows_read,rows_sent,schema,thread_id,user_host) " +
                           "VALUES %s;", strings.Join(values,","))
        // fmt.Printf("--(sql)> %s\n", sql)
        values = []string{}
        go common.HTTPPost(config.DSN_CLICKHOUSE, sql)
      }
    }
  }()
}
