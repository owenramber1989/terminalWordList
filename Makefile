.DEFAULT_GOAL := install

fmt:
	go fmt ./... 
.PHONY:fmt

lint: fmt
	golint ./...
.PHONY:lint

vet: fmt 
	go vet ./...
.PHONY:vet

build:vet
	go build voc.go
.PHONY:build

install: build
	cp voc /usr/local/bin/ 
.PHONY:install
