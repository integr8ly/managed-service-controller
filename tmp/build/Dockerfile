FROM alpine:3.6

RUN adduser -D managed-services-operator
USER managed-services-operator

ADD tmp/_output/bin/managed-services-operator /usr/local/bin/managed-services-operator
