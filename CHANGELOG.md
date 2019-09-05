The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
