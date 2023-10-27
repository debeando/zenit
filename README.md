# ![Zenit](https://raw.githubusercontent.com/debeando/zenit/master/assets/images/zenit_logo.png)
[![Go Report Card](https://goreportcard.com/badge/github.com/debeando/zenit)](https://goreportcard.com/report/github.com/debeando/zenit)

Zenit is a daemon collector for metrics in yours Database's and Proxy's service in small environment. Maybe not requires many another agents for this purpose, but with this one you'll find an excellent tool for database administration.

Why can you use it, this tool is make by DBA for DBA, other tools collect basic information with many services and complex configs, while this it collector low level information and variety than others not, all in one and easy usage.

## Sponsored by

- [OneClick](https://oneclick.es).
- [Security Online Solutions](https://securityonlinesolutions.com).
- [The Hotel Networks](https://www.thehotelsnetwork.com).
- [Zinio](https://www.zinio.com).

## Description:

This agent collect all basic metrics from the hardware and more details from MySQL, MongoDB or ProxySQL services. And the metrics is send only to [InfluxDB](https://github.com/influxdata/influxdb) and you can analize and monitoring with [Grafana](https://grafana.com).

## Advantage

- One agent for all, easy to install and configure, low memory consumption and high performance.
- Auto discover database servers on Amazon Web Services.

## Warnings

- The activation of the Audit and Slow Log compromise the writing performance on disk, and another resources, use another disk for logs and have the necessary resources to support this process.

## Risks

Zenit is not mature, but all database tools can pose a risk to the system and the database server.
Before using this tool, please:

- Read the tool's documentation.
- Review the tool’s known "BUGS".
- Test the tool on a non-production server.

**Like most, you should not be surprised.**

## Install agent

For the moment, this tool only run in any Linux distribution with amd/aarch 64 bits. Paste that at a Terminal prompt:

```bash
bash < <(curl -s https://debeando.com/zenit.sh)
```

For more details, please visit the [wiki](https://github.com/debeando/zenit/wiki).

## How to use it:

See usage with:

```
zenit --help
```
