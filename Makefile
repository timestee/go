SHELL=/bin/bash
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOCOVER=$(GOCMD) tool cover
GOLINT=gometalinter


GO111MODULE=on


.PHONY: setup
setup: ## Install all the build and lint dependencies
	$(GOGET) -u golang.org/x/tools/cmd/cover
	$(GOGET) -u github.com/alecthomas/gometalinter
	$(GOLINT) --install --update


.PHONY: build
build: ## Build the library
	$(GOBUILD) -v  ./...


.PHONY: test
test: ## Run all the tests
	echo 'mode: atomic' > coverage.txt && $(GOTEST) -v -race -covermode=atomic -coverprofile=coverage.txt -timeout=30s ./...


.PHONY: clean
clean:  ## Clean 
	$(GOCLEAN)


.PHONY: cover
cover: test ## Run all the tests and opens the coverage report
	$(GOCOVER) -html=coverage.txt


.PHONY: lint
lint: ## Run all the linters
	$(GOLINT) --tests --errors --disable-all \
		--enable=deadcode \
		--enable=goconst \
		--enable=goimports \
		--enable=ineffassign \
		--enable=interfacer \
		--enable=maligned \
		--enable=misspell \
		--enable=staticcheck \
		--enable=unconvert \
		--enable=varcheck \
		--enable=vet \
		--deadline=10m \
		./...


.PHONY: ci
ci: lint test ## Run all the tests and code checks


.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "%-30s %s\n", $$1, $$2}'

.DEFAULT_GOAL := build
