all: fmt lint vet build

build:
	go build ./

fmt:
	goimports -w ./

lint:
	golint ./...

vet:
	go vet ./...

test:
	go test -v ./...

xc:
	goxc

clean:
	go clean

.PHONY: all build fmt lint vet test xc clean
