# ![Zenit](https://raw.githubusercontent.com/debeando/zenit/master/assets/images/zenit_logo.png)
[![Build Status](https://travis-ci.org/debeando/zenit.svg?branch=master)](https://travis-ci.org/debeando/zenit) [![Coverage Status](https://coveralls.io/repos/github/debeando/zenit/badge.svg?branch=master)](https://coveralls.io/github/debeando/zenit?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/debeando/zenit)](https://goreportcard.com/report/github.com/debeando/zenit) [![Gitter chat](https://badges.gitter.im/Zenit-Agent/Lobby.png)](https://gitter.im/Zenit-Agent/Lobby)

Zenit is a daemon collector for metrics and log parsers for dedicated host for MySQL/Percona/MariaDB Servers and
ProxySQL. Maybe not requires many another agents for this purpose, but with this one you'll find an excellent tool for database administration.

- The name [Zenit](https://en.wikipedia.org/wiki/Zenit_(satellite)) is inspired by a russian spy satellite.
- This project is sponsored by [zinio](https://www.zinio.com) and [The Hotel Networks](https://www.thehotelsnetwork.com).

Why can you use it, this tool is make by DBA for DBA, other tools collect basic information, while this it collector low level information and variety than others not, all in one.

## Description:

This agent collect all basic metrics from the hardware and more details from MySQL or ProxySQL services. And read logs in real time, each event is parsed to analyse later, the logs is send to [ClickHouse](https://github.com/yandex/ClickHouse/) because is very easy to analyse with SQL and have great performance. And the metrics is send only to [InfluxDB](https://github.com/influxdata/influxdb) and you can analize and monitoring with [Grafana](https://grafana.com/) for the moment.

### Advantage

- Auto discover database servers on Amazon Web Services.
- Centralize all logs in a single point of view, you have more details for debugging and analyse, with this you can optimize queries, understand what happen inside, and more.
- Audit database security access and identify possible risk.
- One agent for all, easy to install and configure, low memory consumption and high performance.

### Compatibility

This tool collect stats data from:

- **MySQL:** Collect typical metrics; variables, status, slave status, primary key overflow, tables sizes. The parser Slow and Audit Logs is only tested on MySQL 5.5, the rest of the features work fine with any version.
- **ProxySQL:** Collect for the moment query digest only. For the moment is tested in ProxySQL 1.4
- **AWS RDS Aurora:** Basic metrics; IOPS, CPU, and Replica Lag.
- **Percona ToolKit:** Verify is running specific tools, for the moment only check follow tools; pt-kill, pt-deadlock-logger and pt-slave-delay.
- **Linux OS:** Collect basic metrics of CPU, RAM, DISK, NET, and System Limits.

And this is ingested on:

- **ClickHouse:** This a columnar database to save all log parsers to analyze them.
- **InfluxDB:** Scalable datastore for metrics, events, and real-time analytics.

The numeric values has represent time has in microseconds.

## Warnings

- The activation of the Audit and Slow Log compromise the writing performance on disk, and another resources, use another disk for logs and have the necessary resources to support this process.
- The parse files with very high QPS does big CPU consumption and compromise the server performance. Ensure that you have one available core for this process.

## Risks

Zenit is not mature, but all database tools can pose a risk to the system and the database server.
Before using this tool, please:

- Read the tool's documentation.
- Review the tool’s known "BUGS".
- Test the tool on a non-production server.

**Like most, you should not be surprised.**

## Limitations

- The audit log cut long query.
- ClickHouse no have retention policy for data storage.

## Install agent

For the moment, this tool only run in any Linux distribution with 64 bits. Paste that at a Terminal prompt:

```bash
bash < <(curl -s https://raw.githubusercontent.com/debeando/zenit/master/scripts/install.sh)
```

For more details, please visit the [wiki](https://github.com/debeando/zenit/wiki/Install-agent).

## How to use it:

See usage with:

```
zenit --help
```
