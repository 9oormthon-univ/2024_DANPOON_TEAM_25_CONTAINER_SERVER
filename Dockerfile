# 빌드 단계
FROM golang:1.22-alpine AS builder

# 필요한 빌드 도구 설치
RUN apk add --no-cache git

# 작업 디렉토리 설정
WORKDIR /app

# go.mod와 go.sum 파일을 복사하고 의존성을 다운받음
COPY go.mod ./
RUN go mod download

# 소스 코드 복사
COPY . .

# Go 애플리케이션 빌드 (최적화를 위해 정적 바이너리로 빌드)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

# 실행 단계
FROM alpine:latest

# 빌드된 파일 복사
COPY --from=builder /app/server /app/server

# 실행 포트 설정 (예: 8080)
EXPOSE 50051

# 컨테이너가 시작될 때 실행할 명령어
CMD ["/app/server"]
