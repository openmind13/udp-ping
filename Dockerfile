FROM golang:1.18.3-alpine3.16 AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN GOOS=linux GOARCH=amd64 go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o udp-ping cmd/udp-ping/main.go

FROM --platform=linux/amd64 alpine:3.16
WORKDIR /app
COPY --from=builder /build/udp-ping .
CMD ["/app/udp-ping"]