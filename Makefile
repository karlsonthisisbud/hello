SERVICE_NAME = sku-aggregator

.PHONY: help
help: ## Displays the Makefile help.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: setup
setup: ## Downloads and install various libs for development.
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	go get golang.org/x/tools/cmd/goimports

.PHONY: build
build:  ## Builds project binary.
	go build -race -o ./bin/$(SERVICE_NAME) -v ./cmd/app/app.go

.PHONY: run
run: build ## Builds project binary and executes it.
	bin/$(SERVICE_NAME)