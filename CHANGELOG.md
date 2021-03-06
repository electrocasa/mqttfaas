# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

### Changed

### Deprecated

### Removed

### Fixed

### Security

## [0.0.2] - 2018-07-07

### Added 0.0.2

- Added [temperature convertor](./samples/temp_converter) example
- Updated [README](./README.md) on usage with [mqttcli](https://github.com/shirou/mqttcli)
- Functions can dictate if they should be single use only. (#2)
- Proxy settings are copied over to the container.
- Any environment with prefix `MQTT_FAAS` is available to the functions.
- Added samples on using with [dind](https://hub.docker.com/_/docker/) and [nodered](https://nodered.org/)

### Changed 0.0.2

- `stderr` of functions would be piped with an identifier [commit](https://github.com/sks/mqttfaas/commit/ea1f5fabab9ed38302095f580adc608c824dcf53)

### Deprecated 0.0.2

None

### Removed 0.0.2

None

### Fixed 0.0.2

None

### Security 0.0.2

None

## [0.0.1] - 2018-07-07

- Initial Draft version
