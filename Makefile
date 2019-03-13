GO111MODULE=on

.PHONY: setup
setup: ## Install all the build and lint dependencies
	go get -u golang.org/x/tools/cmd/cover
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install --update

allpackages = $$(go list ./... | grep -v /vendor/)

.PHONY: test
test: ## Run all the tests
	echo 'mode: atomic' > coverage.txt && go test -v -race -covermode=atomic -coverprofile=coverage.txt -timeout=30s ./...

.PHONY: cover
cover: test ## Run all the tests and opens the coverage report
	go tool cover -html=coverage.txt

.PHONY: lint
lint: ## Run all the linters
	gometalinter --skip=gen/v2 --skip=design --vendor --disable-all \
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

.PHONY: build
build: ## Build a binary

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "%-30s %s\n", $$1, $$2}'

.DEFAULT_GOAL := build
