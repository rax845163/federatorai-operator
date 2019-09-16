
# Image URL to use all building/pushing image targets
IMG ?= admission-controller
COMMIT := $(shell git rev-parse --short HEAD)

.PHONY: all test admission-controller 
all: test admission-controller

# Run tests
test: generate fmt vet
	go test . -coverprofile cover.out

# Build admission-controller binary
admission-controller: generate fmt vet
	go build -ldflags "-X main.VERSION=`git rev-parse --abbrev-ref HEAD`-`git rev-parse --short HEAD``git diff --quiet || echo '-dirty'` -X 'main.BUILD_TIME=`date`' -X 'main.GO_VERSION=`go version`'" -o bin/admission-controller github.com/containers-ai/alameda/admission-controller/cmd

.PHONY: run

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet
	go run ./cmd/main.go

.PHONY: fmt vet generate docker-build docker-push

# Run go fmt against code
fmt:
	go fmt .

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate:
	go generate .

# Build the docker image
docker-build: test
	docker build ./.. -t ${IMG}:${COMMIT} -f Dockerfile
