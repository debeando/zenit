-- Query examples for ClickHouse on zenit.mysql_slow_log to analyze data.
--
-- Queries causing the most load:
SELECT
    host_name AS "HostName",
    query_digest,
    count(*) AS count_query,
    round(avg(query_time), 3) AS latency,
    round(quantile(0.99)(query_time), 3) AS latency_percentage,
    round((latency * count_query) / (max(_time) - min(_time)), 3) AS load
FROM zenit.mysql_slow_log
GROUP BY host_name, query_digest
HAVING count_query > 1
  AND minus(max(_time), min(_time)) > 0
ORDER BY host_name, load DESC
LIMIT 10 BY host_name;

-- Queries slow in the last 24h:
SELECT schema AS "Schema",
       halfMD5(query_digest) AS "Fingerprint",
       count() AS "Count",
       MAX(query_time) AS "MaxTime",
       AVG(query_time) AS "AvgTime", query_digest AS "QueryDigest",
       ROUND(AVG(rows_examined), 0) AS "RowsExamined",
       MAX(_time) AS "LastSeen"
FROM zenit.mysql_slow_log
WHERE schema LIKE 'zen\_%'
  AND user_host LIKE 'zen\_%'
  AND query_digest LIKE 'select%'
  AND query_time > 0.1
  AND host_name LIKE '%-prod-%'
  AND _time >= (NOW() - (60 * 60 * 24))
GROUP BY halfMD5(query_digest), query_digest, host_name, schema
LIMIT 10 BY schema
LIMIT 100;

-- Slowest query:
SELECT
    host_name AS "HostName",
    query_digest,
    round(query_time, 4) AS query_time
FROM zenit.mysql_slow_log
ORDER BY host_name, query_time DESC
LIMIT 10 BY host_name
LIMIT 100;

-- Dumping data:
SELECT
    host_name,
    query,
    rows_sent,
    bytes_sent
FROM zenit.mysql_slow_log
WHERE rows_sent > 1000
order by rows_sent desc
limit 10 by host_name;

-- List all hosts and the last log inserted:
SELECT
    host_name,
    MAX(_time),
    count()
FROM zenit.mysql_slow_log
WHERE _time >= (NOW() - (60 * 60 * 24))
GROUP BY host_name;

-- Example 1:
SELECT IPv4NumToString(host_ip) AS "IPAddress",
       host_name AS "HostName",
       halfMD5(query_digest) AS "HashQuery",
       COUNT() AS "Count"
FROM zenit.mysql_slow_log
WHERE _time >= (NOW() - (60 * 60 * 24))
GROUP BY host_ip, host_name, query_digest
ORDER BY host_name ASC, COUNT() DESC;

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
ORDER BY lock_time DESC
LIMIT 10;
