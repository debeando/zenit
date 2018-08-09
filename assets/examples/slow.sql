-- Query examples for ClickHouse on zenit.mysql_slow_log to analyze data.
--
-- Example 1:
SELECT IPv4NumToString(host_ip) AS "IPAddress",
       host_name AS "HostName",
       halfMD5(query_digest) AS "HashQuery",
       COUNT() AS "Count"
FROM zenit.mysql_slow_log
WHERE _time >= (NOW() - (60 * 60 * 24))
GROUP BY host_ip, host_name, query_digest
ORDER BY host_name ASC, COUNT() DESC;

-- Example 2:
SELECT IPv4NumToString(host_ip) AS "IPAddress",
       host_name AS "HostName",
       halfMD5(query_digest) AS "HashQuery",
       substring(query_digest, 1, 80) AS "QueryDigest",
       COUNT() AS "Count",
       AVG(query_time) AS "AvgExecTime"
FROM zenit.mysql_slow_log
WHERE _time >= (NOW() - (60 * 60 * 24))
GROUP BY host_ip, host_name, query_digest
ORDER BY AVG(query_time) DESC
LIMIT 100;

-- Example 3:
SELECT IPv4NumToString(host_ip) AS "IPAddress",
       host_name AS "HostName",
       COUNT() AS "Count"
FROM zenit.mysql_slow_log
GROUP BY host_ip, host_name;

-- Example 4:
SELECT *
FROM zenit.mysql_slow_log
WHERE halfMD5(query_digest) = 11782761010365089099
LIMIT 10;
