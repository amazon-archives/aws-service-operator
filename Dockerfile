FROM alpine
MAINTAINER Christopher Hein <me@christopherhein.com>

RUN apk --no-cache add openssl musl-dev ca-certificates
COPY aws-operator /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/aws-operator"]
