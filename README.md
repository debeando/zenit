# Zenit

Zenit is a daemon collector for metrics and log parsers for dedicated host for MySQL/Percona/Mariadb Servers and
ProxySQL. Maybe no require many another agents for this propouse, with one does excellent tool for database administrator.

The name [Zenit](https://en.wikipedia.org/wiki/Zenit_(satellite)) is inspired a russian was spy satellite.

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

- Read the tool’s documentation.
- Review the tool’s known “BUGS”
- Test the tool on a non-production server.

## Warnings

- The parse files on very high QPS does big CPU consumption and compromise the server.
Asegurar que tiene un core disponible, la herramienta no ocupará mas de un core.
- La activación del Audit y Slow Log compromete el rendimiento de escritura del disco,
para no comprometer el disco de datos, use un disco aparte al de datos exclusivo para
los logs. No olvide rotar los logs para evitar llenar el disco.

## Beneficios

- Centralización de todos los logs en un unico punto de consulta.
- Mejoramos la seguridad evitando acceso de los usuarios a los servidores para
analizar los Logs.
- Ofrece información útil para todos desarrolladores para optimizar las consultas.


## Install

```bash
chown root. zenit
mv zenit /usr/local/bin/
```

## Configuration

### MySQL

Configure slow log:

SET GLOBAL long_query_time = 100;

### ProxySQL

Allow remote access:

```bash
mysql -u admin -padmin -h 127.0.0.1 -P 6032
SET admin-admin_credentials = "admin:admin;radminuser:radminpass";
LOAD ADMIN VARIABLES TO RUNTIME;
```

### ClickHouse

How to test connection?

```bash
curl -s -d 'SELECT 1' http://10.201.17.217:8123/?database=zenit
```

Create database and tables for clickhouse:

```bash
cat assets/schema/clickhouse/zenit.sql | clickhouse-client --multiline
```

### Prometheus



## Build

## Development

```bash
docker run -d --name some-clickhouse-server --ulimit nofile=262144:262144 yandex/clickhouse-server
docker run -it --rm --link some-clickhouse-server:clickhouse-server yandex/clickhouse-client --host clickhouse-server

GOOS=linux go build -ldflags "-s -w" -o zenit main.go && \
docker cp zenit d1c86f2f36ff:/root && \
docker exec -i -t d1c86f2f36ff /root/zenit -collect-os
```

while :; do cat test_slow.log >> /var/lib/mysql/slow.log; sleep 0.1; done
while :; do cat test_audit.log >> /var/lib/mysql/audit.log; sleep 0.1; done
