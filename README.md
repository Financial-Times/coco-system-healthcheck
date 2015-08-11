#coco-system-healthcheck

##Building
```
CGO_ENABLED=0 go build -a -installsuffix cgo -o coco-system-healthcheck .

docker build -t coco/coco-system-healthcheck .
```

