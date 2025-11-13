.PHONY: build test vendor run clean

build: vendor
	go build -o bin/rhdh-ai-install

test: build
	go test -v ./...

vendor:
	go mod vendor

clean:
	rm -rf bin
