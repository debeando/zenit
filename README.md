# Zenit

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
- Each event on logs take the query and digest to help to identify bad queries.
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

## Install

For the moment, this tool only run in any Linux distribution with 64 bits. Paste that at a Terminal prompt:

```bash
bash < <(curl -s https://raw.githubusercontent.com/swapbyt3s/zenit/master/scripts/install.sh)
```

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


## Configuration

On most systems, the default locations are `/etc/zenit/zenit.ini` for the main configuration file.

### Agent Configuration

The configuration is very intuitive, please see the example [config file](https://github.com/swapbyt3s/zenit/blob/master/zenit.ini).
