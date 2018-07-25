# Todo List

Ideas or issues to work:

- In collect/mysql/audit.go:
  > Define percentage level to collect data, us a good way to optimize parser
    performance.
  > Have a filter to collect parcer, for example: collect all diff users like it "!zen_*".
    parser-filter=""
  > Make bulk for audit.log to ClickHouse.
- Add check for logrotate maybe in "/etc/logrotate.d/mysql.conf".
- Refactoring collect/percona/process.go:
  > Use a loop / array with all this common.PGrep("mysqld"). mysqld,proxysql...
  > Use this style: os.process{name="mysqld"} 1
- Add log path: >> /var/log/zenit.log in deamonize process.
- Change pid files path into: /var/run/zenit/*
- Implement infinite loop for collect: --interval=30s
- Catch all panic error, logging and continue running app without exit.
- Parser:
  > General Log
  > Error Log
  > Binarylogs
- Create a docker swarm with MySQL, ProxySQL, ClickHouse, InfluxDB for testing.
- Collect data from ProxySQL to ClickHouse:
  > stats_mysql_query_digest
  > mysql_servers
  > mysql_replication_hostgroups
- Execute specific query and depend of result, send Slack Notifications.
- Add debug parameter.
