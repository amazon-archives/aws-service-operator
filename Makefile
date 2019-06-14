commitSHA := $(shell git describe --dirty --always)
dateStr := $(shell date +%s)
repo ?= awsoperator.io/aws-service-operator

.PHONY: compile-docs
generate-docs:
	mkdir -p website/images
	asciidoctor docs/readme.adoc -o website/index.html
	asciidoctor docs/*.adoc -D website
	cp -r docs/images/* website/images/

.PHONY: serve-docs
serve-docs:
	cd website/ && python3 -m http.server --bind 0.0.0.0 8080

.PHONY: generate
generate:
	ruby -rpry hack/generate/process.rb

.PHONY: test-generate
test-generate:
	rspec hack/generate/spec/

.PHONY: build
build:
	 go build -ldflags "-X main.commit=$(commitSHA) -X main.date=$(dateStr)" ./cmd/awsoperator

.PHONY: build-eks-cluster
build-eks-cluster:
	eksctl create cluster -f hack/eks-cluster.yaml \
		--color=fabulous \
		--kubeconfig kubeconfig

.PHONY: build-kind-cluster
build-kind-cluster:
	kind create cluster --config hack/kind-cluster.yaml \
		--name aws-service-operator

.PHONY: install-tools
install-tools:
	GO111MODULE=off go get -u github.com/jteeuwen/go-bindata/...
	gem install git pry activesupport hana

.PHONY: generate-test-data
generate-test-data:
	go generate ./pkg/testutils/

.PHONY: test
test:
	go test -v ./pkg/cloudformation/
	go test -v ./pkg/queue/
	go test -v ./pkg/encoding/cloudformation/
	go test -v ./pkg/controller-manager/
	go test -v ./pkg/apis/self/v1alpha1/