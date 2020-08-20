ifndef $(GOPATH)
    GOPATH=$(shell go env GOPATH)
    export GOPATH
endif

run: build-docker-server build-docker-migrate
	docker-compose -f ./deployments/docker-compose.yml up -d

down:
	docker-compose -f ./deployments/docker-compose.yml down

build-docker-server:
	docker build -f ./Dockerfile -t antibruteforce/server ./

build-docker-migrate:
	docker build -t antibruteforce/migrate ./migrate

build-docker-int-tests:
	docker build -f ./integration-test.Dockerfile -t antibruteforce/integration-tests ./

test:
	go test -v -count=1 -race -gcflags=-l -timeout=30s ./...

test-integration: build-docker-server build-docker-migrate build-docker-int-tests
	./scripts/run_int_tests.sh

lint:
	golangci-lint run ./...

build:
	CGO_ENABLED=0 GOARCH=amd64 go build -o ./.bin/antibf

run-local:
	source ./configs/.local.env && go run . server

migrate:
	goose -dir migrate/migrations postgres "user=antibruteforceUser password=password dbname=antibruteforce sslmode=disable" up

generate:
	protoc -Iapi -I/usr/local/include  --go_out=plugins=grpc:internal/grpc antibruteforce.proto

.PHONY: run down test test-integration lint build run-local migrate generate
