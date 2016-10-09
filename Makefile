build:
	GOPATH=$(shell pwd) go build -o transfer src/main.go

run:
	GOPATH=$(shell pwd) go run src/main.go
