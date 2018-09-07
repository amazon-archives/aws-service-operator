commitSHA := $(shell git describe --dirty --always)
dateStr := $(shell date +%s)

.PHONY: build
build:
	go build -ldflags "-X main.commit=$(commitSHA) -X main.date=$(dateStr)" ./cmd/aws-operator

.PHONY: release
release:
	goreleaser --rm-dist

.PHONY: dev-release
dev-release:
	goreleaser --rm-dist --snapshot --skip-publish

.PHONY: tag
tag:
	git tag -a ${VERSION} -s
	git push origin --tags

.PHONY: install-aws-codegen
install-aws-codegen:
	go get -u github.com/christopherhein/aws-operator-codegen

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
