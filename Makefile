DBG         ?= 0
PROJECT     ?= federatorai-operator
ORG_PATH    ?= github.com/containers-ai
REPO_PATH   ?= $(ORG_PATH)/$(PROJECT)
VERSION     = $(shell git rev-parse --abbrev-ref HEAD)-$(shell git rev-parse --short HEAD)$(TMP_VERSION_SUFFIX)$(shell git diff --quiet || echo '-dirty')
LD_FLAGS    ?= -X $(REPO_PATH)/pkg/version.Raw=$(VERSION)
BUILD_DEST  ?= bin/federatorai-operator
MUTABLE_TAG ?= latest
IMAGE		= federatorai-operator
IMAGE_TAG	?= $(shell git describe --always --dirty --abbrev=7)

FIRST_GOPATH:=$(firstword $(subst :, ,$(shell go env GOPATH)))
GOBINDATA_BIN=$(FIRST_GOPATH)/bin/go-bindata

ifeq ($(DBG),1)
GOGCFLAGS ?= -gcflags=all="-N -l"
endif

.PHONY: all
all: build images check

NO_DOCKER ?= 0
ifeq ($(NO_DOCKER), 1)
  DOCKER_CMD =
  IMAGE_BUILD_CMD = imagebuilder
else
  DOCKER_CMD := docker run --rm -v "$(PWD):/go/src/$(REPO_PATH):Z" -w "/go/src/$(REPO_PATH)" openshift/origin-release:golang-1.10
  IMAGE_BUILD_CMD = docker build
endif

.PHONY: pkg/assets/bindata.go
pkg/assets/bindata.go: $(GOBINDATA_BIN)
	# Using "-modtime 1" to make generate target deterministic. It sets all file time stamps to unix timestamp 1
	cd assets && $(GOBINDATA_BIN) -pkg assets -o ../$@ \
		CustomResourceDefinition/... \
		ClusterRole/... \
		ServiceAccount/... \
		ClusterRoleBinding/... \
		Secret/... \
		ConfigMap/... \
		PersistentVolumeClaim/... \
		Service/... \
		Deployment/... \
		AlamedaScaler/... \
		Route/... \
		StatefulSet/... \
		Ingress/... \
		DaemonSet/... \
		PodSecurityPolicy/... \
		SecurityContextConstraint/... \




.PHONY: depend
depend: $(GOBINDATA_BIN)
	dep version || go get -u github.com/golang/dep/cmd/dep
	dep ensure

.PHONY: depend-update
depend-update:
	dep ensure -update

.PHONY: build
build: pkg/assets/bindata.go## build binaries
	$(DOCKER_CMD) go build $(GOGCFLAGS) -ldflags "$(LD_FLAGS)" -o "$(BUILD_DEST)" "$(REPO_PATH)/cmd/manager"

.PHONY: images
images: ## Create images
	$(IMAGE_BUILD_CMD) -t "$(IMAGE):$(IMAGE_TAG)" -t "$(IMAGE):$(MUTABLE_TAG)" ./

.PHONY: push
push:
	docker push "$(IMAGE):$(IMAGE_TAG)"
	docker push "$(IMAGE):$(MUTABLE_TAG)"

.PHONY: check
check: fmt vet lint test ## Check your code

.PHONY: test
test: ## Run unit tests
	$(DOCKER_CMD) go test -race -cover ./...

.PHONY: lint
lint: ## Go lint your code
	hack/go-lint.sh -min_confidence 0.3 $(go list -f '{{ .ImportPath }}' ./...)

.PHONY: fmt
fmt: ## Go fmt your code
	hack/go-fmt.sh .

.PHONY: vet
vet: ## Apply go vet to all go files
	hack/go-vet.sh ./...

$(GOBINDATA_BIN):
	go get -u github.com/shuLhan/go-bindata/...

.PHONY: help
help:
	@grep -E '^[a-zA-Z/0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

