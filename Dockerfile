FROM alpine
MAINTAINER Christopher Hein <heichris@amazon.com>

RUN apk --no-cache add openssl musl-dev ca-certificates libc6-compat
COPY aws-service-operator /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/aws-service-operator"]
