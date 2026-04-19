FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd/effective-mobile-test/main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /build/app .

EXPOSE 8080

CMD ["./app"]