-include .env

export

build:
	go build .

run: 
	go run .

lint:
	golangci-lint run

forms:
	golangci-lint fmt