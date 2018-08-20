# Zenit [![Build Status](https://travis-ci.org/swapbyt3s/zenit.svg?branch=master)](https://travis-ci.org/swapbyt3s/zenit)

Zenit is a daemon collector for metrics and log parsers for dedicated host for MySQL/Percona/Mariadb Servers and
ProxySQL. Maybe not requires many another agents for this purpose, but with this one you'll find an excellent tool for database administration.

The name [Zenit](https://en.wikipedia.org/wiki/Zenit_(satellite)) is inspired by a russian spy satellite.

## Description:

This agent collect all basic metrics from the hardware and more details from MySQL or ProxySQL services.
And read logs in real time, each event is parsed to analize later, the logs is send to [ClickHouse](https://github.com/yandex/ClickHouse/)
because is very easy to analize with SQL and have great performance. And the metrics is send only to [Prometheus](https://github.com/prometheus/prometheus)
for the moment.

### Advantage

- Centralize all logs in a single point of view.
- Each event on logs take the query and digest to help to identify bad or malformed queries.
- Improve security to prevent user access into server.
- Provider useful information for developers to help optimization queries.

### Compatibility

This tool collect stats data from:

- **Linux OS (CentOS):** Collect basic metrics of CPU, RAM, DISK, NET, and System Limits.
- **MySQL:** Collect tipical metrics; variables, status, slave status, primary key overflow, tables sizes. And parser Slow and Audit Logs. For the moment is tested on MySQL 5.5
- **Percona ToolKit:** Verify is running specific tools, for the moment only check follow tools; pt-kill, pt-deadlock-logger and pt-slave-delay.
- **ProxySQL:** Collect for the moment query digest only. For the moment is tested in ProxySQL 1.4

And this is ingested on:

- **Prometheus:** This another metric tools, good for alerts by metrics generated with zenit.
- **ClickHouse:** This a columnar database to save all log parsers to analyze them.

The numeric values has represent time has in microseconds.

## Warnings

- The parse files with very high QPS does big CPU consumption and compromise the server performance. Ensure that you have
available core for this process.
- The activation of the Audit and Slow Log compromise the writing performance on disk, use another disk for logs.

## Risks

Zenit is not mature, but all database tools can pose a risk to the system and the database server.
Before using this tool, please:

- Read the tool's documentation.
- Review the tool’s known "BUGS".
- Test the tool on a non-production server.

## Install zenit agent

For the moment, this tool only run in any Linux distribution with 64 bits. Paste that at a Terminal prompt:

```bash
bash < <(curl -s https://raw.githubusercontent.com/swapbyt3s/zenit/master/scripts/install.sh)
```

### Configure zenit agent

By default configuration file are in `/etc/zenit/zenit.ini`.

#### Agent Configuration

The configuration is very intuitive, please see the example [config file](https://github.com/swapbyt3s/zenit/blob/master/zenit.ini).

## How to use it:

See usage with:

```
./zenit --help
```

#### Run zenit in quiet mode:

```
./zenit --quiet
```

#### Run zenit in daemon mode:

Runs in the background and detach from bash.

```
./zenit --start
```

#### Stop zenit in daemon mode:

```
./zenit --stop
```

## Configure ClickHouse

First, check you have connection to ClickHouse server, for this example the server it is in `127.0.0.1`. Try the follow command:

```bash
$ curl -s -d 'SELECT 1' 'http://127.0.0.1:8123/?database=system'
```

With user and password:

```bash
$ curl -s -d 'SELECT 1' 'http://127.0.0.1:8123/?user=admin&password=admin&database=system'
```

If all is well, the server will respond with the value one (1).If you have a problem, check the [ClickHouse settings](https://clickhouse.yandex/docs/en/operations/access_rights/).

Second, you will need to create the database and the tables into ClickHouse using this [sql script](https://github.com/swapbyt3s/zenit/blob/master/assets/schema/clickhouse/zenit.sql).

```bash
cat assets/schema/clickhouse/zenit.sql | clickhouse-client --multiline
```

## Exploring query problems

In ClickHouse you can find bad or malformed queries, or access log, and group by similar queries digested to find the more long execution time. You are free and use your imagination to find problem, please see this examples:

- [SlowLog](https://github.com/swapbyt3s/zenit/blob/master/assets/examples/slow.sql)
- [AuditLog](https://github.com/swapbyt3s/zenit/blob/master/assets/examples/audit.sql)
