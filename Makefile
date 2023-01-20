#!/usr/bin/make -f

.PHONY: test build docker-test docker-build docker-image

test:
	go test ./...

build:
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -tags netgo -ldflags '-w' .

docker-test:
	docker run -t -v $$PWD:/go/src/github.com/yannh/redis-dump-go -w /go/src/github.com/yannh/redis-dump-go golang:1.18 make test

docker-build:
	docker run -t -v $$PWD:/go/src/github.com/yannh/redis-dump-go -w /go/src/github.com/yannh/redis-dump-go golang:1.18 make build

docker-image:
	docker build -t contentful-labs/terraform-diff:latest .
