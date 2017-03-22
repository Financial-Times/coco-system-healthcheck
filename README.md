#coco-system-healthcheck

## Building:
```
CGO_ENABLED=0 go build -a -installsuffix cgo -o coco-system-healthcheck .

docker build -t coco/coco-system-healthcheck .
```

## Actual checks:
* Root disk space
* Persistent disk space
* Memory load
* CPU load average 
* NTP synch
* TCP connection available in 8080