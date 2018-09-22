Changelog
=========

See [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).


2.0.0
-----

### Breaking Changes
- Package name changed from `httpsimplified` to `httpsimp`.
- Added `*http.Client` argument to `Get`, `Post`, `Put` and `Perform`.
- Renamed `EncodeBody` into `EncodeForm`.
- Renamed `BasicAuth` into `BasicAuthValue`.
- Parsers can no longer be used directly; use `Parse` function instead.
- A single non-public error is used instead of `StatusError` and `ContentTypeError`.

### New Features
- Added support for multiple body parsers.
- Added support for customizing status codes and content types expected by body parsers.


1.1.0 — 2018-09-21
------------------

### Added
- Added `PlainText` result parser.


1.0.0 — 2017-12-22
------------------

- First public stable release.