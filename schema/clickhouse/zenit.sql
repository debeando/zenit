CREATE DATABASE IF NOT EXISTS zenit;

DROP TABLE IF EXISTS zenit.mysql_audit_log;
CREATE TABLE IF NOT EXISTS zenit.mysql_audit_log (
  _time DateTime default now(),
  _date Date default toDate(_time),
  host_ip UInt32,
  name String,
  command_class String,
  connection_id UInt64,
  status UInt64,
  sqltext String,
  sqltext_digest String,
  user String,
  priv_user String,
  os_login String,
  proxy_user String,
  host String,
  os_user String,
  ip String,
  db String
) ENGINE = MergeTree(_date,(_time,host,user), 8192);

DROP TABLE IF EXISTS zenit.mysql_slow_log;
CREATE TABLE IF NOT EXISTS zenit.mysql_slow_log (
  _time DateTime default now(),
  _date Date default toDate(_time),
  host_ip UInt32,
  bytes_sent UInt64,
  killed UInt64,
  last_errno UInt64,
  lock_time Float64,
  query String,
  query_time Float64,
  query_digest String,
  rows_affected UInt64,
  rows_examined UInt64,
  rows_read UInt64,
  rows_sent UInt64,
  schema String,
  thread_id UInt64,
  user_host String
) ENGINE = MergeTree(_date,(_time,user_host), 8192);
