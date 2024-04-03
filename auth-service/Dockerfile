# build
FROM golang:1.22.1-alpine AS builder

WORKDIR /app

RUN apk --no-cache add bash git make

COPY ["go.mod", "go.sum", "./"]

RUN go mod download

COPY . .

ENV CGO_ENABLED=0

RUN go build -o ./bin/app cmd/api/main.go

# run
FROM alpine as runner

COPY --from=builder /app/bin/app /
COPY .env /config.env

EXPOSE 4000

CMD ["/app"]