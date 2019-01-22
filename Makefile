commitSHA := $(shell git describe --dirty --always)
dateStr := $(shell date +%s)
repo ?= github.com/awslab/aws-service-operator

.PHONY: build
build:
	go build -ldflags "-X main.commit=$(commitSHA) -X main.date=$(dateStr)" ./cmd/aws-service-operator

.PHONY: release
release:
	goreleaser --rm-dist

.PHONY: dev-release
dev-release:
	goreleaser --rm-dist --snapshot --skip-publish

.PHONY: test
test:
	go test -v -cover -race $(repo)/...

.PHONY: tag
tag:
	git tag -a ${VERSION} -s
	git push awslabs --tags

.PHONY: install-aws-codegen
install-aws-codegen:
	${MAKE} -C code-generation install

.PHONY: aws-codegen
aws-codegen:
	aws-service-operator-codegen process

.PHONY: k8s-codegen
k8s-codegen:
	./hack/update-codegen.sh

.PHONY: codegen
codegen: aws-codegen k8s-codegen

.PHONY: rebuild
rebuild: codegen build
