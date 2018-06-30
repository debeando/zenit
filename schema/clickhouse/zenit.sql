CREATE DATABASE IF NOT EXISTS zenit;

CREATE TABLE IF NOT EXISTS zenit.mysql_audit_log (
  _time DateTime default now(),
  _date Date default toDate(_time),
  name String,
  command_class String,
  connection_id UInt32,
  status UInt32,
  sqltext String,
  user String,
  priv_user String,
  os_login String,
  proxy_user String,
  host String,
  os_user String,
  ip String,
  db String
) ENGINE = MergeTree(_date,(_time,host,user), 8192);
