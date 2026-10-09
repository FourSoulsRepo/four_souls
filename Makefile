# Every Go module in the workspace (see go.work).
MODULES := . ./pkg/rules_engine ./pkg/card_db ./pkg/record

# Extra build tags, e.g. TAGS=webkit2_41 on Ubuntu 24.04.
TAGS ?=

# golangci-lint binary; must be built with a Go version >= the local toolchain.
GOLANGCI ?= golangci-lint
ROOT := $(CURDIR)

.PHONY: test vet build wasm lint screenshots

test:
	@for m in $(MODULES); do (cd $$m && go test -tags "$(TAGS)" ./... && go test -tags "$(TAGS) embed" ./...) || exit 1; done

vet:
	@for m in $(MODULES); do (cd $$m && go vet -tags "$(TAGS)" ./... && go vet -tags "$(TAGS) embed" ./...) || exit 1; done

build:
	@for m in $(MODULES); do (cd $$m && go build -tags "$(TAGS)" ./...) || exit 1; done

# The engine and the website must always build for the browser (A-09).
wasm:
	@cd pkg/rules_engine && GOOS=js GOARCH=wasm go build ./...
	@GOOS=js GOARCH=wasm go build -o /dev/null ./cmd/website

# Strict linters with security rules (PR-03); one config for every module.
lint:
	@for m in $(MODULES); do (cd $$m && $(GOLANGCI) run --config $(ROOT)/.golangci.yml --build-tags "$(TAGS)" ./... && $(GOLANGCI) run --config $(ROOT)/.golangci.yml --build-tags "$(TAGS) embed" ./...) || exit 1; done

# UI screenshots for review: .cache/screenshots/screen-*.png (git-ignored).
# Fake Wails bindings via ?screenshot=1; no Go process needed.
screenshots:
	@cd scripts/screenshots && npm install --no-audit --no-fund && npx playwright install chromium && node capture.mjs
