FROM golang:1.14-alpine AS builder

WORKDIR /app
COPY . .

#RUN sleep 20
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/antibf

FROM scratch

COPY --from=builder /go/bin/antibf /go/bin/antibf
ENTRYPOINT ["/go/bin/antibf"]
CMD ["server"]
EXPOSE 8081
