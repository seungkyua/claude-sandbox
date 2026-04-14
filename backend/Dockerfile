# 빌드 스테이지
FROM golang:1.19-alpine AS builder

WORKDIR /app

# 의존성 설치
COPY go.mod go.sum ./
RUN go mod download

# 소스 복사 및 빌드
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server/

# 실행 스테이지
FROM alpine:3.18

WORKDIR /app

# 필요한 패키지 설치
RUN apk --no-cache add ca-certificates tzdata

# 바이너리 및 설정 파일 복사
COPY --from=builder /app/server .
COPY --from=builder /app/config.yaml.example ./config.yaml.example

# 업로드 디렉토리 생성
RUN mkdir -p /app/uploads

# 비루트 사용자 설정
RUN adduser -D -g '' appuser
USER appuser

EXPOSE 8080

CMD ["./server"]
