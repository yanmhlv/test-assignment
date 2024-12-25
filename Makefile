.PHONY: build
build:
	go build -o ./bin ./cmd/...

.PHONY: tests
tests:
	go test -race -count=5 ./...
