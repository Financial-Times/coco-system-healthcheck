FROM golang

RUN go get github.com/Financial-Times/coco-system-healthcheck
RUN cd $GOPATH/src/github.com/Financial-Times/coco-system-healthcheck && git checkout dockerize && go install
RUN cp $GOPATH/bin/coco-system-healthcheck /system-healthcheck

EXPOSE 8080

WORKDIR /

CMD /system-healthcheck --hostPath=$HOST_DIR

