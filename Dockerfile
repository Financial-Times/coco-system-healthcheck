FROM gliderlabs/alpine:3.2

ADD . /system-healthcheck
RUN apk --update add go git\
  && export GOPATH=/.gopath \
  && go get github.com/Financial-Times/coco-system-healthcheck \
  && cd system-healthcheck \
  && go build \
  && mv system-healthcheck /coco-system-healthcheck \
  && apk del go git \
  && rm -rf $GOPATH /var/cache/apk/*

EXPOSE 8080
CMD /coco-system-healthcheck --hostPath=$HOST_DIR

