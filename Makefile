.PHONY: all
all: clean test build

.PHONY: clean
clean:
	rm -rf ./dist

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: lint
lint:
	CGO_ENABLED=0 golangci-lint run --enable-all -D exhaustivestruct

.PHONY: test
test:
	CGO_ENABLED=0 go test -coverprofile cover.out ./...
	go tool cover -html cover.out -o cover.html

.PHONY: build
build:
	goreleaser build --snapshot --rm-dist
