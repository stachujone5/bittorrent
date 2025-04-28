.PHONY: build lint fmt

run: build
	./bin/main

build:
	go build -o ./bin/main ./main.go

lint:
	golangci-lint run

fmt:
	golangci-lint fmt

test-bencode:
	cd bencode && go test ./...

test: test-bencode
