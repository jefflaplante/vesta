.PHONY: build test clean install lint fmt vet

BINARY := vesta
GO := go

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE    ?= $(shell date -u +%Y-%m-%d)
LDFLAGS := -s -w -X github.com/jeff/vesta/cmd.Version=$(VERSION) \
                 -X github.com/jeff/vesta/cmd.Commit=$(COMMIT) \
                 -X github.com/jeff/vesta/cmd.BuildDate=$(DATE)

build:
	$(GO) mod verify
	CGO_ENABLED=0 $(GO) build -trimpath -ldflags="$(LDFLAGS)" -o $(BINARY) .

test:
	$(GO) test ./...

test-v:
	$(GO) test -v ./...

clean:
	rm -f $(BINARY)
	$(GO) clean

install:
	$(GO) install .

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

lint: fmt vet

all: clean build test
