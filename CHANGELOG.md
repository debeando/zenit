# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### 1.1.10 (Beta)

- Implement alerts system.
- Add alert for Errors Connections on ProxySQL.
- Add basic alerts for OS: Disk, Mem, CPU.
- Add basic alerts for MySQL: Max Connection, Lagging, Replication, ReadOnly.
- Refactoring logging.
- Allow environment variables in config file.
- Refactoring config skeleton to show hierarchy for collect and alerts.
- Send alerts notifications to Slack.
- BUG: Prometheus exporter unexpected end of input stream. (Issue #56)
