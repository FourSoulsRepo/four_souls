# Every Go module in the workspace (see go.work).
MODULES := . ./pkg/rules_engine ./pkg/card_db ./pkg/record

.PHONY: test vet build

test:
	@for m in $(MODULES); do (cd $$m && go test ./...) || exit 1; done

vet:
	@for m in $(MODULES); do (cd $$m && go vet ./...) || exit 1; done

build:
	@for m in $(MODULES); do (cd $$m && go build ./...) || exit 1; done
