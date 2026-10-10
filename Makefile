# Every Go module in the workspace (see go.work).
MODULES := . ./pkg/rules_engine ./pkg/card_db ./pkg/record

# Extra build tags, e.g. TAGS=webkit2_41 on Ubuntu 24.04.
TAGS ?=

# golangci-lint binary; must be built with a Go version >= the local toolchain.
GOLANGCI ?= golangci-lint
ROOT := $(CURDIR)

# Fuzzing (docs/wiki/Fuzzing.md): how long to fuzz, how long to shrink a
# failing input, which cards fuzz-focused keeps in play (comma-separated;
# empty: the test's own list), how many games sims plays.
FUZZTIME ?= 10m
MINIMIZE ?= 5s
FOCUS ?=
SIMS ?= 3000

.DEFAULT_GOAL := help
.PHONY: help test vet build wasm lint frontend e2e screenshots fuzz fuzz-engine fuzz-focused fuzz-sets sims

help: ## Show this help
	@echo "Usage: make <target> [VAR=value]"
	@echo
	@awk 'BEGIN {FS = ":.*## "} /^[a-z-]+:.*## / {printf "  %-14s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo
	@echo "Variables: TAGS, GOLANGCI, FUZZTIME=$(FUZZTIME), MINIMIZE=$(MINIMIZE), FOCUS=$(FOCUS), SIMS=$(SIMS)"

test: ## Run tests in every Go module, plus embed mode
	@for m in $(MODULES); do (cd $$m && go test -tags "$(TAGS)" ./... && go test -tags "$(TAGS) embed" ./...) || exit 1; done

vet: ## Run go vet in every Go module, plus embed mode
	@for m in $(MODULES); do (cd $$m && go vet -tags "$(TAGS)" ./... && go vet -tags "$(TAGS) embed" ./...) || exit 1; done

build: ## Build every Go module
	@for m in $(MODULES); do (cd $$m && go build -tags "$(TAGS)" ./...) || exit 1; done

# The engine and the website must always build for the browser (A-09).
wasm: ## Build the engine and the website for js/wasm
	@cd pkg/rules_engine && GOOS=js GOARCH=wasm go build ./...
	@GOOS=js GOARCH=wasm go build -o /dev/null ./cmd/website

# Strict linters with security rules (PR-03); one config for every module.
lint: ## Run golangci-lint in every Go module, plus embed mode
	@for m in $(MODULES); do (cd $$m && $(GOLANGCI) run --config $(ROOT)/.golangci.yml --build-tags "$(TAGS)" ./... && $(GOLANGCI) run --config $(ROOT)/.golangci.yml --build-tags "$(TAGS) embed" ./...) || exit 1; done

# Frontend checks for every commit (ADR 008).
frontend: ## Typecheck, lint and unit-test the frontend
	@cd frontend && npm run typecheck && npm run lint && npm test

# Playwright lives in e2e/, so frontend installs never download Chromium.
# Fake Wails bindings via ?fake; no Go process needed.
E2E := cd e2e && npm install --no-audit --no-fund && npx playwright install chromium &&

e2e: ## Run end-to-end UI tests (Playwright)
	@$(E2E) npx playwright test

screenshots: ## Run e2e and save UI screenshots into .cache/screenshots
	@$(E2E) SCREENSHOTS=1 npx playwright test

# Go fuzzes one target at a time, so fuzz runs them one after another.
fuzz: fuzz-engine fuzz-sets ## Run fuzz-engine, then fuzz-sets (FUZZTIME each)

fuzz-engine: ## Fuzz the engine with its test cards (FuzzAllCards)
	@cd pkg/rules_engine && go test -run '^$$' -fuzz=FuzzAllCards -fuzztime=$(FUZZTIME) -fuzzminimizetime=$(MINIMIZE) .

fuzz-focused: ## Fuzz with FOCUS cards in every play area (FuzzFocused)
	@cd pkg/rules_engine && FOUR_SOULS_FOCUS=$(FOCUS) go test -run '^$$' -fuzz=FuzzFocused -fuzztime=$(FUZZTIME) -fuzzminimizetime=$(MINIMIZE) .

fuzz-sets: ## Fuzz with every real card set (FuzzSets)
	@cd pkg/rules_engine && go test -run '^$$' -fuzz=FuzzSets -fuzztime=$(FUZZTIME) -fuzzminimizetime=$(MINIMIZE) ./cards/

sims: ## Play SIMS whole random games with the real sets
	@cd pkg/rules_engine && FOUR_SOULS_SIMS=$(SIMS) go test -count=1 -timeout 60m -run '^TestSimulations$$' -v ./cards/
