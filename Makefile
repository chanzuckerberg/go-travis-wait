release: build ## run a release
	bff bump
	git push
	goreleaser release --rm-dist

build: ## build the binary
	go build .

install: ## install the go-travis-wait binary in $GOPATH/bin
	go install .

help: ## display help for this makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

clean: ## clean the repo
	rm go-travis-wait 2>/dev/null || true
	go clean
	rm -rf dist

.PHONY: build clean install release help
