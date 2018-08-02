commitSHA := $(shell git describe --dirty --always)
dateStr := $(shell date +%s)

.PHONY: build
build:
	go build -ldflags "-X main.commit=$(commitSHA) -X main.date=$(dateStr)" ./cmd/aws-operator

.PHONY: release
release:
	rm -fr dist
	goreleaser

.PHONY: install-bindata
install-bindata:
	go get -u github.com/jteeuwen/go-bindata/...

.PHONY: install-aws-codegen
install-aws-codegen:
	go get -u github.com/christopherhein/aws-operator-codegen

# .PHONY: update-bindata
# update-bindata:
# 	go generate ./pkg/cloudformation/

.PHONY: aws-codegen
aws-codegen:
	aws-operator-codegen process

.PHONY: k8s-codegen
k8s-codegen:
	./codegen.sh

.PHONY: codegen
codegen: aws-codegen k8s-codegen

.PHONY: rebuild
rebuild: codegen build
