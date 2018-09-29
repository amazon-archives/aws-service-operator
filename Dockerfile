FROM alpine
MAINTAINER Christopher Hein <heichris@amazon.com>

RUN apk --no-cache add openssl musl-dev ca-certificates
COPY aws-service-operator /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/aws-service-operator"]
