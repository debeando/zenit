# ![Zenit](https://raw.githubusercontent.com/swapbyt3s/zenit/master/assets/images/zenit_logo.png)
[![Build Status](https://travis-ci.org/swapbyt3s/zenit.svg?branch=master)](https://travis-ci.org/swapbyt3s/zenit) [![Coverage Status](https://coveralls.io/repos/github/swapbyt3s/zenit/badge.svg?branch=master)](https://coveralls.io/github/swapbyt3s/zenit?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/swapbyt3s/zenit)](https://goreportcard.com/report/github.com/swapbyt3s/zenit) [![Gitter chat](https://badges.gitter.im/Zenit-Agent/Lobby.png)](https://gitter.im/Zenit-Agent/Lobby)

Zenit is a daemon collector for metrics and log parsers for dedicated host for MySQL/Percona/MariaDB Servers and
ProxySQL. Maybe not requires many another agents for this purpose, but with this one you'll find an excellent tool for database administration.

The name [Zenit](https://en.wikipedia.org/wiki/Zenit_(satellite)) is inspired by a russian spy satellite.

## Description:

This agent collect all basic metrics from the hardware and more details from MySQL or ProxySQL services.
And read logs in real time, each event is parsed to analyse later, the logs is send to [ClickHouse](https://github.com/yandex/ClickHouse/)
because is very easy to analyse with SQL and have great performance. And the metrics is send only to [Prometheus](https://github.com/prometheus/prometheus)
for the moment.

### Advantage

- Centralize all logs in a single point of view, you have more details for debugging and analyse, with this you can optimize queries, understand what happen inside, and more.
- Audit database security access and identify possible risk.
- Monitoring and alerting system to prevent a disaster or identify possible risk.
- One agent for all, easy to install and configure, low memory consumption and high performance.

### Compatibility

This tool collect stats data from:

- **MySQL:** Collect typical metrics; variables, status, slave status, primary key overflow, tables sizes. And parser Slow and Audit Logs. For the moment is tested on MySQL 5.5
- **ProxySQL:** Collect for the moment query digest only. For the moment is tested in ProxySQL 1.4
- **Percona ToolKit:** Verify is running specific tools, for the moment only check follow tools; pt-kill, pt-deadlock-logger and pt-slave-delay.
- **Linux OS:** Collect basic metrics of CPU, RAM, DISK, NET, and System Limits.

And this is ingested on:

- **Prometheus:** This another metric tools, good for alerts by metrics generated with zenit.
- **ClickHouse:** This a columnar database to save all log parsers to analyze them.

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

## Limitations

- The audit log cut long query.
- ClickHouse no have retention policy for data storage.

## Install zenit agent

For the moment, this tool only run in any Linux distribution with 64 bits. Paste that at a Terminal prompt:

```bash
bash < <(curl -s https://raw.githubusercontent.com/swapbyt3s/zenit/master/scripts/install.sh)
```

### Configure zenit agent

By default configuration file are in `/etc/zenit/zenit.yaml`.

The configuration file is very intuitive, please see the example [config file](https://github.com/swapbyt3s/zenit/blob/master/zenit.yaml).

**Important:** The hostname for your server cannot be "localhost." The host name should be a unique name.

However, you may need to manually restart the agent (for example, after changing your agent configuration). For Linux systems, the agent selects an init system depending on your operating system version.

For Linux systems, ensure you use the correct command for your init system. Select start, stop, restart, or status as appropriate:

How do I find out what version of Linux I'm running?

```
cat /etc/*{release,version}
```

SystemD (SLES 12, CentOS 7, Debian 8, Debian 9, RHEL 7, Ubuntu 15.04 or higher):

```
sudo systemctl <start|stop|restart|status> zenit
```

System V (Debian 7, SLES 11.4):

```
sudo /etc/init.d/zenit <start|stop|restart|status>
```

Upstart (Amazon Linux, CentOS 6, RHEL 6, Ubuntu 14.04 or lower):

```
sudo initctl <start|stop|restart|status> zenit
```
