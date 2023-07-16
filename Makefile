SHELL = /bin/bash
.SHELLFLAGS = -o pipefail -c

.PHONY: help
help: ## Print info about all commands
	@echo "Commands:"
	@echo
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "    \033[01;32m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build all executables
	go build -ldflags "-X main._version=v0.1.4" ./atr.go
b:
	@make build

install: ## Install all executables
	go install -ldflags "-X main._version=`git tag --sort=-version:refname | head -n 1`" ./atr.go
i:
	@make install

.PHONY: fmt
fmt: ## Format all go files
	go fmt ./*.go

.PHONY: all
all: build