VERSION ?= v0.0.0
#VERSION ?= $(shell git describe --tags --abbrev=0 | sed 's/^v//')
GIT_REV := $(shell git describe --always --abbrev --dirty --long)

LD_FLAGS = -ldflags="-X main.Version=$(VERSION) -X main.GitRev=$(GIT_REV)"
BUILD := CGO_ENABLED=0 go build $(LD_FLAGS)

build: vendor clean test kosh ## Clean, test, and build all the things. You probably want this target
first: tools

.PHONY: clean
clean: ## Clean up the local environment
	rm -f kosh

kosh: ## Build kosh
	$(BUILD) -o kosh ./...

.PHONY: tools
tools: ## Download and install all dev/code tools
	@echo "==> Installing dev tools"
	go get -u github.com/golang/dep/cmd/dep
	go get -u honnef.co/go/tools/cmd/staticcheck

vendor: ## Install dependencies
	dep ensure -v -vendor-only

.PHONY: deps
deps: ## Update dependencies to latest version
	dep ensure -v

.PHONY: test
test: ## Ensure that code matches best practices
	staticcheck ./...

.PHONY: help
help: ## Display this help message
	@echo "GNU make(1) targets:"
	@grep -E '^[a-zA-Z_.-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'


