-- Query examples for ClickHouse on zenit.mysql_audit_log to analyze data.

SELECT DISTINCT
    user,
    count(*)
FROM zenit.mysql_audit_log
GROUP BY user
ORDER BY count() DESC;

SELECT DISTINCT user, sqltext, _time
FROM zenit.mysql_audit_log
WHERE halfMD5(sqltext_digest) = 18043511469391647841 LIMIT 1;
