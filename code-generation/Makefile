commitSHA := $(shell git describe --dirty --always)
dateStr := $(shell date +%s)

.PHONY: build
build:
	go build -ldflags "-X main.commit=$(commitSHA) -X main.date=$(dateStr)" ./cmd/aws-service-operator-codegen

.PHONY: install
install: rebuild
	go install -ldflags "-X main.commit=$(commitSHA) -X main.date=$(dateStr)" ./cmd/aws-service-operator-codegen

.PHONY: install-bindata
install-bindata:
	go get -u github.com/jteeuwen/go-bindata/...

.PHONY: update-bindata
update-bindata:
	go generate ./pkg/codegen/
	go fmt ./...

.PHONY: rebuild
rebuild: update-bindata build
