# 빌드 단계
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git docker-cli

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/grpc

FROM alpine:latest

COPY --from=builder /app/server /app/server

COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

EXPOSE 50051

CMD ["/bin/sh", "/entrypoint.sh"]
