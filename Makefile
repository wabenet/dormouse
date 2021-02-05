all: clean test build

.PHONY: clean
clean:
	rm -f dormouse_*

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: lint
lint:
	golangci-lint run --enable-all -D exhaustivestruct

.PHONY: test
test:
	go test -cover ./...

.PHONY: build
build:
	gox -arch="amd64" -os="darwin linux" ./...
