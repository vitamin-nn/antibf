ifndef $(GOPATH)
    GOPATH=$(shell go env GOPATH)
    export GOPATH
endif

build:
	go build -o ./.bin/antibrtfrs ./cmd/antibruteforce/main.go

test:
	go test -race ./...

install-deps:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(GOPATH)/bin v1.27.0

lint: #install-deps
	golangci-lint run ./...

run:
	go run ./cmd/antibruteforce/main.go

migrate:
	goose -dir migrations postgres "user=vdudov password=123qwe dbname=antibruteforce sslmode=disable" up

generate:
	protoc -Iapi -I/usr/local/include  --go_out=plugins=grpc:internal/grpc antibruteforce.proto
.PHONY: build
