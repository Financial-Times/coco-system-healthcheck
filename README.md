#coco-system-healthcheck
[![CircleCI](https://circleci.com/gh/Financial-Times/coco-system-healthcheck.svg?style=shield)](https://circleci.com/gh/Financial-Times/coco-system-healthcheck)[![Coverage Status](https://coveralls.io/repos/github/Financial-Times/coco-system-healthcheck/badge.svg)](https://coveralls.io/github/Financial-Times/coco-system-healthcheck)
## Building:
```
CGO_ENABLED=0 go build -a -installsuffix cgo -o coco-system-healthcheck .

docker build -t coco/coco-system-healthcheck .
```

## Actual checks:
* Root disk space
* AWS EBS mounts disk space
* Memory load
* CPU load average 
* NTP sync
* TCP connection available in 8080
* K8S API server certificate doesn't expire in one month