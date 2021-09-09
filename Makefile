.PHONY: all
all: help

clean: ## Cleans built artifacts
	go clean
	rm -rf dist/

build: ## Builds binaries with goreleaser config
	goreleaser build --snapshot

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
