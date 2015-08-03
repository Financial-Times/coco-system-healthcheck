FROM golang

RUN go get github.com/Financial-Times/coco-system-healthcheck
RUN cp $GOPATH/bin/coco-system-healthcheck /system-healthcheck

EXPOSE 8080

WORKDIR /

CMD /system-healthcheck --hostPath=$HOST_DIR

