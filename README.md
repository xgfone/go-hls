# Golang RFC8216 HLS

[![Build Status](https://github.com/xgfone/go-hls/actions/workflows/go.yml/badge.svg)](https://github.com/xgfone/go-hls/actions/workflows/go.yml)
[![GoDoc](https://pkg.go.dev/badge/github.com/xgfone/go-hls)](https://pkg.go.dev/github.com/xgfone/go-hls)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=flat-square)](https://raw.githubusercontent.com/xgfone/go-hls/master/LICENSE)
![Minimum Go Version](https://img.shields.io/github/go-mod/go-version/xgfone/go-hls?label=Go%2B)
![Latest SemVer](https://img.shields.io/github/v/tag/xgfone/go-hls?sort=semver)

An golang implementation of RFC8216 HLS.

## Install

```shell
$ go get -u github.com/xgfone/go-hls
```

## PlayList Tags

- **Basic Tags** [RFC 8216, 4.3.1](https://datatracker.ietf.org/doc/html/rfc8216#section-4.3.1)
  - [x] `#EXTM3U` [RFC 8216, 4.3.1.1](https://datatracker.ietf.org/doc/html/rfc8216#section-4.3.1.1)
  - [x] `#EXT-X-VERSION` [RFC 8216, 4.3.1.2](https://datatracker.ietf.org/doc/html/rfc8216#section-4.3.1.2)
- **Media Segment Tags** [RFC 8216, 4.3.2](https://datatracker.ietf.org/doc/html/rfc8216#section-4.3.2)
  - [x] `#EXTINF` [RFC 8216, 4.3.2.1](https://datatracker.ietf.org/doc/html/rfc8216#section-4.3.2.1)
  - [x] `#EXT-X-BYTERANGE` [RFC 8216, 4.3.2.2](https://datatracker.ietf.org/doc/html/rfc8216#section-4.3.2.2)
  - [x] `#EXT-X-DISCONTINUITY` [RFC 8216, 4.3.2.3](https://datatracker.ietf.org/doc/html/rfc8216#section-4.3.2.3)
  - [x] `#EXT-X-KEY` [RFC 8216, 4.3.2.4](https://datatracker.ietf.org/doc/html/rfc8216#section-4.3.2.4)
  - [x] `#EXT-X-MAP` [RFC 8216, 4.3.2.5](https://datatracker.ietf.org/doc/html/rfc8216#section-4.3.2.5)
  - [x] `#EXT-X-PROGRAM-DATE-TIME` [RFC 8216, 4.3.2.6](https://datatracker.ietf.org/doc/html/rfc8216#section-4.3.2.6)
  - [ ] `#EXT-X-DATERANGE` [RFC 8216, 4.3.2.7](https://datatracker.ietf.org/doc/html/rfc8216#section-4.3.2.7)
- **Media Playlist Tags** [RFC 8216, 4.3.3](https://datatracker.ietf.org/doc/html/rfc8216#section-4.3.3)
  - [x] `#EXT-X-TARGETDURATION` [RFC 8216, 4.3.3.1](https://datatracker.ietf.org/doc/html/rfc8216#section-4.3.3.1)
  - [x] `#EXT-X-MEDIA-SEQUENCE` [RFC 8216, 4.3.3.2](https://datatracker.ietf.org/doc/html/rfc8216#section-4.3.3.2)
  - [x] `#EXT-X-DISCONTINUITY-SEQUENCE` [RFC 8216, 4.3.3.3](https://datatracker.ietf.org/doc/html/rfc8216#section-4.3.3.3)
  - [x] `#EXT-X-ENDLIST` [RFC 8216, 4.3.3.4](https://datatracker.ietf.org/doc/html/rfc8216#section-4.3.3.4)
  - [x] `#EXT-X-PLAYLIST-TYPE` [RFC 8216, 4.3.3.5](https://datatracker.ietf.org/doc/html/rfc8216#section-4.3.3.5)
  - [x] `#EXT-X-I-FRAMES-ONLY` [RFC 8216, 4.3.3.6](https://datatracker.ietf.org/doc/html/rfc8216#section-4.3.3.6)
- **Master Playlist Tags** [RFC 8216, 4.3.4](https://datatracker.ietf.org/doc/html/rfc8216#section-4.3.4)
  - [x] `#EXT-X-MEDIA` [RFC 8216, 4.3.4.1](https://datatracker.ietf.org/doc/html/rfc8216#section-4.3.4.1)
  - [x] `#EXT-X-STREAM-INF` [RFC 8216, 4.3.4.2](https://datatracker.ietf.org/doc/html/rfc8216#section-4.3.4.2)
  - [x] `#EXT-X-I-FRAME-STREAM-INF` [RFC 8216, 4.3.4.3](https://datatracker.ietf.org/doc/html/rfc8216#section-4.3.4.3)
  - [x] `#EXT-X-SESSION-DATA` [RFC 8216, 4.3.4.4](https://datatracker.ietf.org/doc/html/rfc8216#section-4.3.4.4)
  - [x] `#EXT-X-SESSION-KEY` [RFC 8216, 4.3.4.5](https://datatracker.ietf.org/doc/html/rfc8216#section-4.3.4.5)
- **Media or Master Playlist Tags** [RFC 8216, 4.3.5](https://datatracker.ietf.org/doc/html/rfc8216#section-4.3.5)
  - [x] `#EXT-X-INDEPENDENT-SEGMENTS` [RFC 8216, 4.3.5.1](https://datatracker.ietf.org/doc/html/rfc8216#section-4.3.5.1)
  - [x] `#EXT-X-START` [RFC 8216, 4.3.5.2](https://datatracker.ietf.org/doc/html/rfc8216#section-4.3.5.2)

### Difference with RFC8216 for `#EXT-X-KEY`

When a key in one `KEYFORMAT` is updated or overwritten, all keys in other `KEYFORMAT`s must be updated simultaneously.
