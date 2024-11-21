#!/bin/bash
set -e

# Docker 소켓 마운트 여부 확인
if [ ! -S /var/run/docker.sock ]; then
    echo "Docker socket not found. Ensure it's mounted correctly."
    exit 1
fi

# Docker 호스트 환경 변수 설정
export DOCKER_HOST=unix:///var/run/docker.sock

# Go 애플리케이션 실행
exec /app/server
