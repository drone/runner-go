The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased
### [1.5.1] - 2019-12-10
- not trimming pipeline history causing memory leak

## [1.5.0] - 2019-12-09
- support for global environment variables
- support for external environment variables from an external service
- abstraction for pipeline execution

## [1.4.0] - 2019-11-05
### Added
- function to encode registry credentials in docker config.json format

## [1.3.1] - 2019-11-01
### Fixed
- check if last exit code greater than 0 in powershell

## [1.3.0] - 2019-10-31
### Fixed
- text overflow for long commit messages
- error in step should bubble up to stage

### Added
- support for legacy CI_ environment variables
- support for registry plugins
- support for concurrency limits in yaml
- support for nodes in yaml
- helpers for working with docker auth config files
- helpers for tagging containers with labels

## [1.2.2] - 2019-09-28
### Fixed
- powershell scripts should check last exit code

### Added
- support for cron events in the dashboard screen
- support for promote events in the dashboard screen
- support for rollback events in the dashboard screen

## [1.2.1] - 2019-07-27
### Fixed
- close and already closed channel in livelog causes panic

## [1.2.0] - 2019-07-27
### Added
- semver environment variables

## [1.1.0] - 2019-07-14
### Added
- logrus hook to store recent system logs
- handler to visualize recent system logs
- handler to visualize pipeline steps
- disable dashbaord if no password set

## [1.0.0] - 2019-07-01
### Added
- defined runner manifest schema
- defined runner remote protocol
- helpers for generating environment variables
- helpers for generating clone scripts
- helpers for generating shell scripts
- support for encrypted secrets
- support for static secrets
- support for remote secrets
- support for buffered log streaming
- handler to provide healtcheck support
- handler to provide runner dashboard
