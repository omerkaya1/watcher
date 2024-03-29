BUILD= $(CURDIR)/bin
$(shell mkdir -p $(BUILD))
VERSION= $(shell git rev-list HEAD --count)
export GO111MODULE=on
export GOPATH=$(go env GOPATH)

.PHONY: setup
setup: ## Install all the build and lint dependencies
	go get -u golang.org/x/tools
	go get -u golang.org/x/lint/golint

.PHONY: mod
mod: ## Runs mod
	go mod verify
	go mod vendor
	go mod tidy

.PHONY: fmt
fmt: setup ## Run goimports on all go files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do goimports -w "$$file"; done

.PHONY: lint
lint: setup ## Runs all the linters
	golint ./internal ./cmd ./

.PHONY: build
build: ## Builds the project
	go build -o $(BUILD)/watcher $(CURDIR)

.PHONY: dockerbuild
dockerbuild: mod ## Builds a docker image with a project
	docker build -t omer513/watcher:0.${VERSION} .

.PHONY: dockerpush
dockerpush: dockerbuild ## Publishes the docker image to the registry
	docker push omer513/watcher:0.${VERSION}

.PHONY: clean
clean: ## Remove temporary files
	go clean $(CURDIR)
	rm -rf $(BUILD)
	rm -rf coverage.txt

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := build
