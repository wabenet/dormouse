PLATFORMS=darwin linux
ARCHITECTURES=386 amd64

.PHONY: all
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
	CGO_ENABLED=0 go test -cover ./...

.PHONY: build
build:
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), \
	$(shell CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -v -o dormouse_$(GOOS)_$(GOARCH) ./cmd/dormouse)))
