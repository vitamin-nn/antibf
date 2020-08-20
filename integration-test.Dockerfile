FROM golang:1.14-alpine

WORKDIR /app
COPY . .

ENV CGO_ENABLED 0
CMD go test -v -tags integration ./test/integration
