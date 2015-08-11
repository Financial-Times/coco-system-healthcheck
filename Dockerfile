FROM gliderlabs/alpine:3.2
ADD coco-system-healthcheck /coco-system-healthcheck
EXPOSE 8080
CMD /coco-system-healthcheck --hostPath=$HOST_DIR

