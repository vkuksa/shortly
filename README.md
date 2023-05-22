# shortly
[![License](https://img.shields.io/github/license/vkuksa/shortly)](https://github.com/vkuksa/shortly/blob/main/LICENSE)
[![Tag](https://img.shields.io/github/v/tag/vkuksa/shortly)](https://github.com/vkuksa/shortly/tags)
[![Go Report](https://goreportcard.com/badge/github.com/vkuksa/shortly)](https://goreportcard.com/report/github.com/vkuksa/shortly)
[![ci](https://github.com/vkuksa/shortly/actions/workflows/ci.yaml/badge.svg)](https://github.com/vkuksa/shortly/actions/workflows/ci.yaml)
[![Coverage](https://codecov.io/gh/vkuksa/shortly/branch/dev/graph/badge.svg)](https://codecov.io/gh/vkuksa/shortly)

---

A link shortening service.

Has integration with Prometheus, Grafana and Alertmanager.

As a datastorage can utilise bbolt or store data in memory as a map.

Configuration is performed via shortly.conf.
Currently only 2 kinds of DB are supported: bbolt and inmem.

# Installation

## Source
```console
$ git clone https://github.com/vkuksa/shortly.git && cd shortly
$ make
```

# Usage

```console
$ ./bin/shortly
```

# Testing
```console
$ make test
```

# Linting
```console
$ make lint
```

# Docker compose
Supports docker-compose

# Prometheus
Supported metrics:
- `shortly_http_request_count` (Total number of requests by route) 
- `shortly_http_request_seconds` (Total amount of request time by route, in seconds)
- `shortly_error_count` (Errors notification)

# Grafana
Prometheus supported as datasource. Refer to `./grafana/provisioning/datasources/prometheus_ds.yml` configuration file.

# Alertmanager
Triggers upon alert of shortly_error_count with "internal" appears.

### Important: Refer to `./prometheus/alertmanager.yml` for setup of receiver.
