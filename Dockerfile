FROM golang:alpine AS builder

RUN apk add --no-cache build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app ./cmd

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /root/
COPY --from=builder /app/app .

EXPOSE 4000

CMD ["./app"]
