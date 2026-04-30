## sandbox container 띄우기
```
$ cat << 'EOF'  | tee run_sandbox.sh
#!/bin/bash

# 사용법 체크
if [ $# -lt 1 ]; then
  echo "Usage: $(basename "$0") <container-name>"
  exit 1
fi

CONTAINER_NAME="$1"

nerdctl run -d --name "$CONTAINER_NAME" \
  --network host \
  -v "$(pwd):$(pwd)" \
  -w "$(pwd)" \
  -v "$HOME/.claude:/Users/$(id -un)/.claude" \
  -v "/Users/$(id -un)/Documents/works/ai:/Users/$(id -un)/Documents/works/ai" \
  claude-container

# 컨테이너가 running 상태가 될 때까지 대기
echo "Waiting for container '$CONTAINER_NAME' to be running..."
TIMEOUT=30
ELAPSED=0
while true; do
  STATUS=$(nerdctl inspect -f '{{.State.Status}}' "$CONTAINER_NAME" 2>/dev/null)
  if [ "$STATUS" = "running" ]; then
    echo "Container is running."
    break
  fi
  if [ "$ELAPSED" -ge "$TIMEOUT" ]; then
    echo "Error: Container did not start within ${TIMEOUT}s (last status: ${STATUS:-unknown})"
    exit 1
  fi
  sleep 1
  ELAPSED=$((ELAPSED + 1))
done

# .claude.json 복사 및 소유권 변경
nerdctl cp ~/.claude.json "$CONTAINER_NAME":/Users/ask.ahn/.claude.json
nerdctl exec -u 0 "$CONTAINER_NAME" chown "$(id -un):$(id -g)" /Users/ask.ahn/.claude.json

echo "Setup complete."

nerdctl exec -it "$CONTAINER_NAME" bash 2>/dev/null
EOF


$ run_sandbox.sh claude-harness-sdk
``` 



모든 구현이 완료되었습니다. 결과 요약:

  구현 완료 결과

  백엔드 (Go + Gin)

  - 9개 패키지 전체 테스트 통과
  - config, dto, handler, middleware, model, repository, router, service, server
  - 주요 기능: 인증(JWT/bcrypt), 플러그인 CRUD, 버전 관리, 설치 관리, 리뷰 시스템, 관리자 기능

  프론트엔드 (Next.js + Tailwind CSS)

  - 11개 테스트 파일, 51개 테스트 모두 통과
  - 주요 기능: API 클라이언트, Zustand 스토어, 공통 UI 컴포넌트, 인증 폼, 레이아웃, 플러그인 카드/리스트/리뷰, 대시보드, 관리자 페이지

  인프라

  - docker-compose.yml — PostgreSQL + 백엔드 서버
  - backend/Dockerfile — 멀티스테이지 빌드

  실행 방법

  # 백엔드
  docker-compose up
  # 또는
  cd backend && go run ./cmd/server/

  # 프론트엔드
  cd frontend && npm run dev

  # 테스트
  cd backend && go test ./...
  cd frontend && npx vitest run

  plan.md 44/44 Step 전부 완료 (100%) 되었습니다.