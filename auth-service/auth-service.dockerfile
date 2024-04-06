# build
FROM --platform=linux/amd64 golang:1.22.1-alpine AS builder

WORKDIR /app

RUN apk --no-cache add bash git make g++

COPY ["auth-service/go.mod", "auth-service/go.sum", "./"]

RUN go mod download

COPY auth-service .

RUN GO111MODULE=on CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -tags musl --ldflags "-extldflags -static" -o ./bin/app cmd/*.go

# run
FROM alpine as runner

COPY --from=builder /app/bin/app /
COPY auth-service/.env /config.env
COPY auth-service/internal/storage/storage.db /

EXPOSE 4000

CMD ["/app"]