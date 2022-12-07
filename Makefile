build:
	@go build -o bin/gobankapi

run: build
	@./bin/gobankapi

test:
	@go test -v ./...