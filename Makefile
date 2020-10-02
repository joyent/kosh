VERSION ?= $(shell git describe --tags --abbrev=0 | sed 's/^v//')

build: vendor clean test all ## Test and build binaries for local architecture
first: tools

.PHONY: clean
clean: ## Remove build products from bin/ and release/
	rm -rf bin
	rm -rf release

.PHONY: tools
tools: ## Download and install all dev/code tools
	@echo "==> Installing dev tools"
	go get -u honnef.co/go/tools/cmd/staticcheck

vendor: ## Install dependencies
	go mod vendor

.PHONY: deps
deps: ## Update dependencies to latest version
	go mod verify

.PHONY: test
test: ## Ensure that code matches best practices
	staticcheck ./...
	go test

.PHONY: help
help: ## Display this help message
	@echo "GNU make(1) targets:"
	@grep -E '^[a-zA-Z_.-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: docker_test
docker_test: ## run a test build in docker
	docker/test.bash

.PHONY: docker_release
docker_release: ## Build all release binaries and checksums in docker
	docker/release.bash


################################
# Dynamic Fanciness            #
# aka Why GNU make Is Required #
################################

PLATFORMS  := darwin-amd64 linux-amd64 solaris-amd64 freebsd-amd64 openbsd-amd64 linux-arm
BINARIES   := kosh
RELEASE_BINARIES := kosh

BINS       := $(foreach bin,$(BINARIES),bin/$(bin))
RELEASES   := $(foreach bin,$(RELEASE_BINARIES),release/$(bin))

GIT_REV    := $(shell git describe --always --abbrev --dirty --long)
LD_FLAGS   := -ldflags="-X main.Version=$(VERSION) -X main.GitRev=$(GIT_REV)"
BUILD      := CGO_ENABLED=0 go build $(LD_FLAGS)

####

all: $(BINS) ## Build all binaries

.PHONY: release
release: vendor test $(RELEASES) ## Build release binaries with checksums

bin/%:
	@mkdir -p bin
	@echo "> Building bin/$(subst bin/,,$@)"
	@$(BUILD) -o bin/$(subst bin/,,$@) *.go

os   = $(firstword $(subst -, ,$1))
arch = $(lastword $(subst -, ,$1))

define release_me
	$(eval BIN:=$(subst release/,,$@))
	$(eval GOOS:=$(call os, $(platform)))
	$(eval GOARCH:=$(call arch, $(platform)))
	$(eval RPATH:=release/$(BIN)-$(GOOS)-$(GOARCH))

	@echo "> Building $(RPATH)"
	@GOOS=$(GOOS) GOARCH=$(GOARCH) $(BUILD) -o $(RPATH) *.go
	shasum -a 256 $(RPATH) > $(RPATH).sha256
endef


release/%:
	@mkdir -p release
	$(foreach platform,$(PLATFORMS),$(call release_me))


