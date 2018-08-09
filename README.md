# Zenit

Zenit is a daemon collector for metrics and log parsers for dedicated host for MySQL/Percona/Mariadb Servers and
ProxySQL. Maybe not requires many another agents for this purpose, but with this one you'll find an excellent tool for database administration

The name [Zenit](https://en.wikipedia.org/wiki/Zenit_(satellite)) is inspired by a russian spy satellite.

## Description:

This tool collect stats data from:

- **Linux OS (CentOS):** Collect basic metrics of CPU, RAM, DISK, NET, and System Limits.
- **MySQL:** Collect tipical metrics; variables, status, slave status, primary key overflow, tables sizes. And parser Slow and Audit Logs.
- **Percona ToolKit:** Verify is running specific tools, for the moment only check follow tools; pt-kill, pt-deadlock-logger and pt-slave-delay.
- **ProxySQL:** Collect for the moment query digest only.

And this is ingested on:

- **Prometheus:** This another metric tools, good for alerts by metrics generated with zenit.
- **ClickHouse:** This a columnar database to save all log parsers to analyze them.

The numeric values has represent time has in microseconds.

## RISKS!

Zenit is not mature, but all database tools can pose a risk to the system and the database server.
Before using this tool, please:

- Read the tool's documentation.
- Review the tool’s known "BUGS".
- Test the tool on a non-production server.

## Warnings

- The parse files with very high QPS does big CPU consumption and compromise the server performance. Ensure that you have
available core for this process.
- The activation of the Audit and Slow Log compromise the writing performance on disk, use another disk for logs.

## Advantage

- Centralize all logs in a single point of view.
- Improve security to prevent user access into server.
- Provider useful information for developers to help optimization queries.
