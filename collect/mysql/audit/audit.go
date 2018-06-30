// This parse is for OLD Style
// https://dev.mysql.com/doc/refman/5.5/en/audit-log-file-formats.html

package audit

import (
  "bufio"
  "fmt"
  "os"
  "regexp"
  "strings"
  "gitlab.com/swapbyt3s/zenit/common"
)

type Audit struct {
  name          string
  timestamp     string
  command_class string
  connection_id string
  status        string
  sqltext       string
  user          string
  priv_user     string
  os_login      string
  proxy_user    string
  host          string
  os_user       string
  ip            string
  db            string
}

func Record(audit_log string) []string {
  r:=regexp.MustCompile(`(?s)<AUDIT_RECORD(.*?)/>`)
  return r.FindAllString(audit_log, -1)
}

func Value(s string, attribute string) string {
  re := regexp.MustCompile(`\s` + attribute + `="(.*?)"`)
  match := re.FindStringSubmatch(s)

  if len(match) == 2 {
    return match[1]
  }

  return ""
}

func Fields(audit_log string) Audit {
  return Audit {
    name: Value(audit_log, "NAME"),
    timestamp: common.ToDateTime(Value(audit_log, "TIMESTAMP")),
    command_class: Value(audit_log, "COMMAND_CLASS"),
    connection_id: Value(audit_log, "CONNECTION_ID"),
    status: Value(audit_log, "STATUS"),
    sqltext: common.Escape(Value(audit_log, "SQLTEXT")),
    user: Value(audit_log, "USER"),
    priv_user: Value(audit_log, "PRIV_USER"),
    os_login: Value(audit_log, "OS_LOGIN"),
    proxy_user: Value(audit_log, "PROXY_USER"),
    host: Value(audit_log, "HOST"),
    os_user: Value(audit_log, "OS_USER"),
    ip: Value(audit_log, "IP"),
    db: Value(audit_log, "DB"),
  }
}

func Parse(path string) (<-chan Audit) {
  channel_record := make(chan Audit)

  file, err := os.Open(path)
  if err != nil {
    return nil
  }

  scanner := bufio.NewScanner(file)
  if err := scanner.Err(); err != nil {
    return nil
  }

  go func() {
    var record_buffer string

    for scanner.Scan() {
      record_buffer += scanner.Text()
      record := Record(record_buffer)

      if len(record) > 0 {
        record_buffer = ""
        channel_record <- Fields(strings.Join(record, ""))
      }
    }

    file.Close()
    close(channel_record)
  }()

  return channel_record
}

func ToSQL(record Audit) string {
  return fmt.Sprintf("INSERT INTO zenit.mysql_audit_log (_time,name,command_class,connection_id,status,sqltext,user,priv_user,os_login,proxy_user,host,os_user,ip,db) VALUES ('%s','%s','%s',%s,%s,'%s','%s','%s','%s','%s','%s','%s','%s','%s');",
    record.timestamp,
    record.name,
    record.command_class,
    record.connection_id,
    record.status,
    record.sqltext,
    record.user,
    record.priv_user,
    record.os_login,
    record.proxy_user,
    record.host,
    record.os_user,
    record.ip,
    record.db,
  )
}