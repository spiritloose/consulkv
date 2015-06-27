all: fmt lint vet build

build:
	go build ./

fmt:
	goimports -w ./

lint:
	golint ./...

vet:
	go vet ./...

.PHONY: all build lint vet fmt
