# shortly
[![License](https://img.shields.io/github/license/vkuksa/shortly)](https://github.com/vkuksa/shortly/blob/main/LICENSE)
[![Tag](https://img.shields.io/github/v/tag/vkuksa/shortly)](https://github.com/vkuksa/shortly/tags)
[![Go Report](https://goreportcard.com/badge/github.com/vkuksa/shortly)](https://goreportcard.com/report/github.com/vkuksa/shortly)
[![main](https://github.com/vkuksa/shortly/actions/workflows/main.yaml/badge.svg)](https://github.com/vkuksa/shortly/actions/workflows/main.yaml)

---

A link shortening service.

Intention was in creation of a simple functionality showcasing usage of Clean Architecture, Prometheus, Grafana and Alertmanager.


# Installation

## Source
```console
$ git clone https://github.com/vkuksa/shortly.git && cd shortly
$ make
```

# Testing
```console
$ make test
```

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
