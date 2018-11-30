FROM alpine:3.6

RUN adduser -D managed-services-controller
USER managed-services-controller

ADD build/_output/bin/managed-service-controller /usr/local/bin/managed-service-controller
