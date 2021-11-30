# Changelog

## [v1.12.0](https://github.com/drone/runner-go/tree/v1.12.0) (2021-11-30)

[Full Changelog](https://github.com/drone/runner-go/compare/v1.11.0...v1.12.0)

**Implemented enhancements:**

- Add `retries` option to clone manifest [\#22](https://github.com/drone/runner-go/pull/22) ([julienduchesne](https://github.com/julienduchesne))

## [v1.11.0](https://github.com/drone/runner-go/tree/v1.11.0) (2021-11-11)

[Full Changelog](https://github.com/drone/runner-go/compare/v1.10.0...v1.11.0)

**Implemented enhancements:**

- create card path env variable [\#19](https://github.com/drone/runner-go/pull/19) ([eoinmcafee00](https://github.com/eoinmcafee00))

**Merged pull requests:**

- Release/1.11.0 [\#21](https://github.com/drone/runner-go/pull/21) ([eoinmcafee00](https://github.com/eoinmcafee00))

## [v1.10.0](https://github.com/drone/runner-go/tree/v1.10.0) (2021-11-10)

[Full Changelog](https://github.com/drone/runner-go/compare/v1.9.0...v1.10.0)

**Implemented enhancements:**

- read & upload card data to drone server [\#16](https://github.com/drone/runner-go/pull/16) ([eoinmcafee00](https://github.com/eoinmcafee00))

**Merged pull requests:**

- Release/1.10.0 [\#18](https://github.com/drone/runner-go/pull/18) ([eoinmcafee00](https://github.com/eoinmcafee00))
- feat\(proxy\): support for all\_proxy variables [\#15](https://github.com/drone/runner-go/pull/15) ([ysicing](https://github.com/ysicing))

## [v1.9.0](https://github.com/drone/runner-go/tree/v1.9.0) (2021-08-26)

[Full Changelog](https://github.com/drone/runner-go/compare/v1.8.0...v1.9.0)

**Implemented enhancements:**

- \(feat\) add silent version of bash and powershell Script [\#13](https://github.com/drone/runner-go/pull/13) ([tphoney](https://github.com/tphoney))

**Merged pull requests:**

- \(maint\) v1.9.0 release prep [\#14](https://github.com/drone/runner-go/pull/14) ([tphoney](https://github.com/tphoney))

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.8.0] - 2021-06-24
### Fixed
- graceful shutdown of http servers
- url escape of 'machine' parameter

### Added
- environment variable for build trigger
- environment variable for pull request title

## [1.7.0] - 2021-03-01
### Fixed
- panic when registry uri parsing errors
- do not mask single-character secrets
- capture stage duration on failure
- capture dag errors
- capture oom kill and exit code
- cancel step on semaphore deadline exceeded

### Added
- support for running a single pipeline on-demand
- support for interpolating global environment variables
- function for creating netrc environment variables
- support for debug mode
- adding depends_on, image and detached fields to step

### Updated
- upgrade drone-go dependency version

## [1.6.0] - 2020-03-24
### Added
- support for username/password in docker config.json
- support for multiple external environment providers
- support for calendar version environment variables
- drain response body to ensure connection re-use

## [1.5.1] - 2019-12-10
### Fixed
- not trimming pipeline history causing memory leak

## [1.5.0] - 2019-12-09
### Added
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


\* *This Changelog was automatically generated by [github_changelog_generator](https://github.com/github-changelog-generator/github-changelog-generator)*
