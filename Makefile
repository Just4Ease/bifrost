GOPATH:=$(shell go env GOPATH)

.PHONY: proto
proto :
	protoc -I proto/ -I${PWD}/ --go_out=plugins=grpc:proto proto/*.proto

.PHONY: build
build: proto
	go build -o srv *.go

.PHONY: test
test:
	go test --race -v ./... -cover

.PHONY: docker
docker:
	docker build . -t ms-kyc:alpine

local:
	go run .