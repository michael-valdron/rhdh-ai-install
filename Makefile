.PHONY: build test vendor run clean

build: vendor
	go build -o bin/rhdh-ai-install

test: vendor
	go test -v ./...

vendor:
	go mod vendor

run: vendor
	go run main.go

clean:
	rm -rf bin
