FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd/effective-mobile-test/main.go
RUN go build -o migrator ./cmd/migrator/main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /build/app .
COPY --from=builder /build/migrator .

COPY --from=builder /build/migrations ./migrations

EXPOSE 8080

CMD ["./app"]