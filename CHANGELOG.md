# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### 1.2.0 (Beta)

- Implement alerts system.
- Add alert for Errors and Status connections on ProxySQL.
- Add basic alerts for OS: Disk, Mem, CPU.
- Add basic alerts for MySQL: Max Connection, Lagging, Replication, ReadOnly.
- Add basic alerts for ProxySQL: Errors, Commands, Queries and Status by hosts.
- Refactoring logging.
- Allow environment variables in config file.
- Refactoring config skeleton to show hierarchy for collect and alerts.
- Send alerts notifications to Slack.
- Prometheus exporter unexpected end of input stream.
- Replace own OS collector to github.com/shirou/gopsutil

## [1.1.9] 2018/09/19

- Add logrotate for zenit.log and zenit.err.
- Add collect index data from all tables.
- Fix stderr output for all log to stdout.

## [1.1.8] 2018/09/15

- Collect stats about connection pool from ProxySQL.
- Add test for outputs.prometheus.

## [1.1.7] 2018/09/10

- Replace format config file INI to Yaml.
- Add option to enable or disable features in .yaml
- Implement github.com/kardianos/service.

## [1.1.6] 2018/09/07

- Replace all values inside of IN() lists and convert all digest query to lowercase.
- Allow retry process for read log file when if exist.
- Remove all panic method and replace by log.

## [1.1.5] 2018/08/30

- Refactoring slow parser without regular expression.
- Add compatibility for slow logs make by MariaDB.
- Fix slice bounds out of range when parsing slow log.
- Fix get value from one line with one key and value from slow log.

## [1.1.4] 2018/08/24

- Replace own tail by github.com/hpcloud/tail.

## [1.1.3] 2018/08/23

- Bad string quoted backslashes in parser log.
- Fix exit when tail file not exist.

## [1.1.2] 2018/08/20

- Fix daemonize process with nohup.

## [1.1.1] 2018/08/20

- Fix panic: close of closed channel

## [1.1.0] 2018/08/19

- Allow send log parser when is not complete after x time.
- Allow with one generic function send data from Audit Log and Slow Log to ClickHouse.
- Fix no close child process.
