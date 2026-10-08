# Every Go module in the workspace (see go.work).
MODULES := . ./pkg/rules_engine ./pkg/card_db ./pkg/record

# Extra build tags, e.g. TAGS=webkit2_41 on Ubuntu 24.04.
TAGS ?=

.PHONY: test vet build wasm

test:
	@for m in $(MODULES); do (cd $$m && go test -tags "$(TAGS)" ./...) || exit 1; done
	@go test -tags "$(TAGS) embed" ./...

vet:
	@for m in $(MODULES); do (cd $$m && go vet -tags "$(TAGS)" ./...) || exit 1; done
	@go vet -tags "$(TAGS) embed" ./...

build:
	@for m in $(MODULES); do (cd $$m && go build -tags "$(TAGS)" ./...) || exit 1; done

# The engine and the website must always build for the browser (A-09).
wasm:
	@cd pkg/rules_engine && GOOS=js GOARCH=wasm go build ./...
	@GOOS=js GOARCH=wasm go build -o /dev/null ./cmd/website
