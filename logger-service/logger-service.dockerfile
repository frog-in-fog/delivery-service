# build
FROM golang:1.22.1-alpine AS builder

WORKDIR /app

RUN apk --no-cache add bash git make

COPY ["logger-service/go.mod", "logger-service/go.sum", "./"]

RUN go mod download

COPY logger-service .

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/app cmd/*.go

# run
FROM alpine as runner

COPY --from=builder /app/bin/app /

EXPOSE 8081

CMD ["/app"]
