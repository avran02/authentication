FROM golang:latest AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

ENV GOOS linux
ENV GOARCH amd64
ENV CGO_ENABLED=0

RUN go build -o auth-service .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/auth-service .
RUN chmod +x /app/auth-service

ENTRYPOINT ["./auth-service", "serve"]
