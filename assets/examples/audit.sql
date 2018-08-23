-- Query examples for ClickHouse on zenit.mysql_audit_log to analyze data.
--
-- Example 1: List all hosts and the last log inserted:
SELECT
    host_name,
    MAX(_time)
FROM zenit.mysql_audit_log
GROUP BY host_name;

-- Example 2: List all hosts with audit log count:
SELECT
    host_name,
    count(*)
FROM zenit.mysql_audit_log
GROUP BY host_name
ORDER BY count() DESC;

-- Example 3: List all users with audit log count:
SELECT DISTINCT
  user,
  count(*)
FROM zenit.mysql_audit_log
GROUP BY user
ORDER BY count() DESC;

-- Example 4: List all hosts and users with audit log count:
SELECT
  host_name,
  user,
  count(*)
FROM zenit.mysql_audit_log
GROUP BY host_name, user
ORDER BY count() DESC;

-- Example 5: List all queries with errors on execution:
SELECT
  user,
  halfMD5(sqltext_digest),
  sqltext_digest,
  COUNT()
FROM zenit.mysql_audit_log
where name = 'Query'
  AND command_class = 'error'
  AND status > 0
GROUP BY user, halfMD5(sqltext_digest), sqltext_digest
ORDER BY COUNT() DESC;

-- Example 6: List the specific hash of query digest to determine frecuency in
--            specific day.
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
