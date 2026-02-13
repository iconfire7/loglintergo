GOLANGCI ?= ./golangci-lint/golangci-lint
CUSTOM_GCL ?= ./custom-gcl

.PHONY: tidy test custom-linter lint-golangci clean

tidy:
	go mod tidy

test:
	go test ./...

lint-standalone:
	go run ./cmd ./...

custom-linter:
	$(GOLANGCI) custom

lint-golangci: custom-linter
	$(CUSTOM_GCL) run --config .golangci.yml ./testdata/src/...

clean:
	go clean
	rm -f custom-gcl
	rm -rf .cache/golangci-lint
