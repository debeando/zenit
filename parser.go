package main

// log/audit

import (
  "fmt"
  "regexp"
  "gitlab.com/swapbyt3s/zenit/common"
)

const (
  REGEX_AUDIT_GRP = `(?s)<AUDIT_RECORD(.*?)/>`
  REGEX_AUDIT_ATT = `<AUDIT_RECORD` +
                    `\W+NAME="(?P<name>.*)"` +
                    `\W+RECORD="(?P<record>.*)"` +
                    `\W+TIMESTAMP="(?P<timestamp>.*)"` +
                    `(\W+COMMAND_CLASS="(?P<command_class>.*)")?` +
                    `\W+CONNECTION_ID="(?P<connection_id>.*)"` +
                    `\W+STATUS="(?P<status>.*)"` +
                    `(\W+SQLTEXT="(?P<sqltext>.*)")?` +
                    `\W+(USER="(?P<user>.*)")?` +
                    `(\W+PRIV_USER="(?P<priv_user>.*)")?` +
                    `(\W+OS_LOGIN="(?P<os_login>.*)")?` +
                    `(\W+PROXY_USER="(?P<proxy_user>.*)")?` +
                    `(\W+HOST="(?P<host>.*)")?` +
                    `(\W+OS_USER="(?P<os_user>.*)")?` +
                    `\W+IP="(?P<ip>.*)"` +
                    `(\W+DB="(?P<db>.*)")?\n/>`
)

func main() {
  audit_log := `
<AUDIT_RECORD
  NAME="Connect"
  RECORD="68_2017-11-27T13:25:59"
  TIMESTAMP="2017-11-27T13:26:00 UTC"
  CONNECTION_ID="3529526"
  STATUS="0"
  USER="zen_catalog"
  PRIV_USER="zen_catalog"
  OS_LOGIN=""
  PROXY_USER=""
  HOST=""
  IP="10.15.38.154"
  DB="zen_catalog"
/>
<AUDIT_RECORD
  NAME="Query"
  RECORD="69_2017-11-27T13:25:59"
  TIMESTAMP="2017-11-27T13:26:00 UTC"
  COMMAND_CLASS="error"
  CONNECTION_ID="3529526"
  STATUS="0"
  SQLTEXT="SELECT c0_.id AS id0, c0_.object_type AS object_type1, c0_.object_id AS object_id2, c0_.country_code AS country_code3, c0_.remarks AS remarks4, c0_.created_at AS created_at5 FROM country_restrictions c0_ WHERE c0_.object_id = '7874' AND c0_.object_type = '3' ORDER BY c0_.id DESC LIMIT 200 OFFSET 0"
  USER="zen_catalog[zen_catalog] @  [10.15.38.154]"
  HOST=""
  OS_USER=""
  IP="10.15.38.154"
/>
`

  reGroup := common.ExtRegexp{regexp.MustCompile(REGEX_AUDIT_GRP)}
  line  := reGroup.FindAllString(audit_log, -1)
  for i := range line {
    reAttributes := common.ExtRegexp{regexp.MustCompile(REGEX_AUDIT_ATT)}
    if attributes := reAttributes.FindStringSubmatchMap(line[i]); attributes != nil {
      for k, v := range attributes {
        fmt.Printf("-- > %s = %s\n", k, v)
      }
      // fmt.Printf("%#v\n", attributes["user"])
    }
    fmt.Println("--> ----")
  }
}
