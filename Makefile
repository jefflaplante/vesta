.PHONY: build test clean install lint fmt vet

BINARY := vesta
GO := go

build:
	$(GO) mod verify
	CGO_ENABLED=0 $(GO) build -trimpath -ldflags="-s -w" -o $(BINARY) .

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
