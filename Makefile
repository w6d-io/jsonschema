export GO111MODULE  := on
export PATH         := ./bin:${PATH}
export NEXT_TAG     ?=

ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

ifeq (gsed not found,$(shell which gsed))
SEDBIN=sed
else
SEDBIN=$(shell which gsed)
endif

export PATH := $(shell pwd)/bin:${PATH}

##@ General

REF        = $(shell git symbolic-ref --quiet HEAD 2> /dev/null)
VERSION   ?= $(shell basename /$(shell git symbolic-ref --quiet HEAD 2> /dev/null ) )
VCS_REF    = $(shell git rev-parse HEAD)
GOVERSION  = $(shell go env GOVERSION)
BUILD_DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

ifeq (,$(shell go env GOOS))
GOOS       = $(shell echo $OS)
else
GOOS       = $(shell go env GOOS)
endif

ifeq (,$(shell go env GOARCH))
GOARCH     = $(shell echo uname -p)
else
GOARCH     = $(shell go env GOARCH)
endif

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go get $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef

GOIMPORTS  = $(shell pwd)/bin/goimports
bin/goimports: ## Download goimports locally if necessary
	$(call go-get-tool,$(GOIMPORTS),golang.org/x/tools/cmd/goimports)

.PHONY: test
test: fmt vet
	@echo go test -v -coverpkg=./... -coverprofile=cover.out ./...
	@go test -v -coverpkg=./... -coverprofile=cover.out ./...
	@go tool cover -func cover.out | grep total

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

# Formats the code
.PHONY: format
format: bin/goimports
	goimports -w -local gitlab.w6d.io/w6d,github.com/w6dio internal pkg

.PHONY: changelog
changelog:
	git-chglog -o docs/CHANGELOG.md --next-tag $(NEXT_TAG)
