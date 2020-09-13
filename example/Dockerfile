FROM golang:1.14.3 AS builder
WORKDIR /build
COPY *.go go.mod go.sum ./
RUN go mod verify && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -a -installsuffix tetris-example -o app .

FROM alpine:latest
RUN mkdir -p /app
WORKDIR /app
COPY --from=builder /build/app .
ENTRYPOINT ["./app"]