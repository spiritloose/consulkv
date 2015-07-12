VERSION = $(shell awk -F\" '/^var VERSION/ { print $$2 }' main.go)

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

dep:
	go get
	go get golang.org/x/tools/cmd/goimports
	go get golang.org/x/tools/cmd/vet
	go get github.com/golang/lint/golint
	go get github.com/laher/goxc

xc:
	goxc -pv $(VERSION)

clean:
	go clean

.PHONY: all build fmt lint vet test dep xc clean
