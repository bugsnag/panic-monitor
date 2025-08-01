# Changelog

## 2.0.0 (2025-07-31)

### Enhancements

* Add go.mod to lock bugsnag-go notifier version, separated to v2 [#6](https://github.com/bugsnag/panic-monitor/pull/6)

## 1.1.0 (2022-06-09)

### Enhancements

* Detect and report panics from crashes in Cgo pseudo-packages

## 1.0.1 (2021-01-29)

### Bug fixes

* Disable the `bugsnag.Configuration.PanicHandler` in the child process
  automatically
* Update endpoint configuration documentation to specify
  `BUGSNAG_NOTIFY_ENDPOINT` instead of `BUGSNAG_ENDPOINT`

## 1.0.0

Initial release
