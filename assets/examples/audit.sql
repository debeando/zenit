-- Query examples for ClickHouse on zenit.mysql_audit_log to analyze data.

SELECT DISTINCT
  user,
  count(*)
FROM zenit.mysql_audit_log
GROUP BY user
ORDER BY count() DESC;

SELECT
  host_name,
  user,
  count(*)
FROM zenit.mysql_audit_log
GROUP BY host_name, user
ORDER BY count() DESC;

SELECT
  user,
  halfMD5(sqltext_digest),
  sqltext_digest,
  COUNT()
FROM zenit.mysql_audit_log
where name = 'Query'
  AND command_class = 'error'
GROUP BY user, halfMD5(sqltext_digest), sqltext_digest
ORDER BY COUNT() DESC;

SELECT DISTINCT
  user,
  sqltext,
  _time
FROM zenit.mysql_audit_log
WHERE halfMD5(sqltext_digest) = 18043511469391647841 LIMIT 1;

SELECT
  toDate(_time),
  toHour(_time),
  count() AS c,
  bar(c, 0, 12000, 20)
FROM zenit.mysql_audit_log
WHERE halfMD5(sqltext_digest) = 11527302759258400550
  AND toDate(_time) = '2018-08-20'
GROUP BY toDate(_time), toHour(_time)
ORDER BY toDate(_time), toHour(_time);
