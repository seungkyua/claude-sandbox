# 코드 생성 계획 — KTC Claude Plugin Hub

## 프로젝트 정보
- **프로젝트명**: ktc-claude-plugin-hub
- **코드 위치**: `frontend/`, `backend/`
- **기술 스택**: Next.js + Tailwind CSS + Vitest (프론트엔드), Go + Gin + go test (백엔드), PostgreSQL
- **TDD 적용**: Red → Green → Refactor

## 코드 생성 순서

### Phase 1: 프로젝트 구조 및 기반

- [x] Step 1: 백엔드 프로젝트 초기화 — Go 모듈 초기화, Gin/GORM/JWT 의존성, 디렉토리 구조 생성
- [x] Step 2: 프론트엔드 프로젝트 초기화 — Next.js + Tailwind CSS + Vitest 설정, 디렉토리 구조 생성
- [x] Step 3: 백엔드 설정 로드 구현 — config.go, YAML 파일 로드, 환경변수 우선순위 적용
- [x] Step 4: 데이터베이스 연결 구현 — PostgreSQL 연결, GORM 초기화, 마이그레이션 자동 실행

### Phase 2: 데이터베이스 모델 및 공통 유틸리티 (backend/internal/model/, backend/internal/dto/)

- [x] Step 5: 데이터 모델 구현 — User, Category, Plugin, PluginVersion, Review, Installation 모델 정의
- [x] Step 6: DTO 구현 — 요청/응답 DTO, 페이지네이션 공통 구조체, RFC 7807 에러 응답 포맷
- [x] Step 7: 공통 미들웨어 구현 — CORS 미들웨어, JSON 에러 핸들러

### Phase 3: 인증 시스템 (backend/internal/)

- [x] Step 8: JWT 유틸리티 구현 — Access Token 생성/검증, Refresh Token 생성/검증, claims 추출
- [x] Step 9: 인증 미들웨어 구현 — Bearer Token 검증, 사용자 정보 컨텍스트 주입, 관리자 권한 검증
- [x] Step 10: 사용자 Repository 구현 — CreateUser, FindByEmail, FindByID
- [x] Step 11: 인증 Service 구현 — Register (bcrypt 해싱, 이메일 중복 검증), Login (비밀번호 검증, 토큰 발급), RefreshToken
- [x] Step 12: 인증 Handler 구현 — POST /auth/register, POST /auth/login, POST /auth/refresh, GET /me

### Phase 4: 카테고리 관리 (backend/internal/)

- [x] Step 13: 카테고리 Repository 구현 — FindAll, FindByID, Create
- [x] Step 14: 카테고리 Handler 구현 — GET /categories

### Phase 5: 플러그인 CRUD (backend/internal/)

- [x] Step 15: 플러그인 Repository 구현 — Create, FindByID, FindAll (필터/정렬/페이지네이션), Update, Delete, FindByName
- [x] Step 16: 플러그인 Service 구현 — 등록 (관리자: 즉시공개+공식, 일반: pending), 수정/삭제 (본인+관리자 권한 검증), 중복명 검증
- [x] Step 17: 플러그인 Handler 구현 — POST /plugins, GET /plugins, GET /plugins/:id, PUT /plugins/:id, DELETE /plugins/:id

### Phase 6: 버전 관리 (backend/internal/)

- [x] Step 18: 버전 Repository 구현 — Create, FindByPluginID, FindByID, 중복 버전 검증
- [x] Step 19: 버전 Service 구현 — 새 버전 업로드, 파일 저장, 파일 크기/확장자 검증
- [x] Step 20: 버전 Handler 구현 — POST /plugins/:id/versions, GET /plugins/:id/versions/:versionId/download (다운로드 횟수 증가)

### Phase 7: 설치 관리 (backend/internal/)

- [x] Step 21: 설치 Repository 구현 — Create, Delete, FindByUserID, FindByUserAndPlugin, UpdateActive
- [x] Step 22: 설치 Service 구현 — 설치 (다운로드 횟수 증가), 삭제, 활성화/비활성화 토글
- [x] Step 23: 설치 Handler 구현 — POST /plugins/:id/install, DELETE /plugins/:id/install, PATCH /plugins/:id/install, GET /me/installations

### Phase 8: 리뷰 시스템 (backend/internal/)

- [x] Step 24: 리뷰 Repository 구현 — Create, FindByPluginID, FindByUserAndPlugin, Update, Delete
- [x] Step 25: 리뷰 Service 구현 — 작성 (1인1리뷰, 본인 플러그인 불가), 수정, 삭제, 평균 평점/리뷰수 재계산
- [x] Step 26: 리뷰 Handler 구현 — POST /plugins/:id/reviews, GET /plugins/:id/reviews, PUT /plugins/:id/reviews/:reviewId, DELETE /plugins/:id/reviews/:reviewId

### Phase 9: 관리자 기능 (backend/internal/)

- [x] Step 27: 관리자 Service 구현 — 심사 대기 목록 조회, 승인 (status→approved), 반려 (status→rejected, 사유), 비공개 처리 (status→hidden)
- [x] Step 28: 관리자 Handler 구현 — GET /admin/plugins/pending, PATCH /admin/plugins/:id/approve, PATCH /admin/plugins/:id/reject, PATCH /admin/plugins/:id/hide

### Phase 10: 백엔드 라우터 통합 및 서버 엔트리포인트

- [x] Step 29: 라우터 통합 구현 — 전체 라우트 등록, 미들웨어 적용, 그룹 분리 (public/auth/admin)
- [x] Step 30: 서버 엔트리포인트 구현 — main.go, 설정 로드→DB 연결→라우터 초기화→서버 시작

### Phase 11: 프론트엔드 공통 기반 (frontend/src/)

- [x] Step 31: API 클라이언트 구현 — Axios 인스턴스, 인터셉터 (토큰 자동 첨부, 401 시 토큰 갱신), 에러 핸들링
- [x] Step 32: 타입 정의 구현 — auth.ts, plugin.ts, review.ts, api.ts (서버 응답과 동일한 타입)
- [x] Step 33: 전역 상태 관리 구현 — authStore (Zustand: 토큰, 사용자 정보, 로그인/로그아웃), themeStore (다크/라이트 모드)
- [x] Step 34: 공통 UI 컴포넌트 구현 — Button, Input, Modal, Card, Badge, Pagination

### Phase 12: 프론트엔드 인증 UI (frontend/src/app/, frontend/src/components/auth/)

- [x] Step 35: 로그인/회원가입 페이지 구현 — LoginForm, RegisterForm, 입력 검증, API 연동, 토큰 저장
- [x] Step 36: 레이아웃 구현 — Header (로고, 검색바, 로그인/프로필), Footer, 다크모드 토글, 반응형

### Phase 13: 프론트엔드 메인 및 탐색 (frontend/src/app/, frontend/src/components/plugin/)

- [x] Step 37: 메인 페이지 구현 — 공식 플러그인 섹션 (상단), 커뮤니티 플러그인 섹션 (하단), PluginCard 컴포넌트
- [x] Step 38: 검색/탐색 페이지 구현 — 키워드 검색, 카테고리 필터, 정렬 (인기순/최신순/평점순), 페이지네이션
- [x] Step 39: 플러그인 상세 페이지 구현 — 설명, 스크린샷, 버전 정보, 설치 버튼, 다운로드 수, 평점, 리뷰 목록

### Phase 14: 프론트엔드 사용자 대시보드 (frontend/src/app/dashboard/)

- [x] Step 40: 내 설치 플러그인 대시보드 구현 — 설치 플러그인 목록, 활성화/비활성화 토글, 삭제
- [x] Step 41: 개발자 콘솔 구현 — 내가 등록한 플러그인 목록, 새 플러그인 등록 폼, 버전 업데이트, 통계

### Phase 15: 프론트엔드 관리자 페이지 (frontend/src/app/admin/)

- [x] Step 42: 관리자 페이지 구현 — 심사 대기 목록, 승인/반려 (사유 입력), 비공개 처리, 공식 플러그인 관리

### Phase 16: 통합 및 빌드

- [ ] Step 43: 프론트엔드-백엔드 통합 테스트 — API 연동 E2E 검증, CORS 동작 확인
- [ ] Step 44: 빌드 및 배포 설정 — 프론트엔드 빌드 최적화, 백엔드 Docker 설정, 환경별 설정 분리

## 스토리 매핑

| 기능 요구사항 | 구현 Step |
|--------------|----------|
| 회원가입/로그인 (이메일/비밀번호) | Step 8-12, 35 |
| 플러그인 검색 및 필터링 | Step 15-17, 38 |
| 플러그인 상세 페이지 | Step 17, 39 |
| 플러그인 설치/삭제 | Step 21-23, 40 |
| 내 플러그인 대시보드 (활성화/비활성화) | Step 21-23, 40 |
| 플러그인 등록/수정/삭제 | Step 15-17, 41 |
| 버전 관리 | Step 18-20, 41 |
| 관리자 심사/승인/비공개 처리 | Step 27-28, 42 |
| 공식 플러그인 배포 (상단 노출) | Step 16-17, 37 |
| 평점 및 리뷰 시스템 | Step 24-26, 39 |
| 카테고리 분류 | Step 13-14, 38 |
| 다운로드 횟수 집계 | Step 20, 22, 39 |
| 다크 모드 지원 | Step 33, 36 |
| JWT 인증 (Access + Refresh Token) | Step 8-9, 31 |
| 보안 (bcrypt, CORS, 입력 검증) | Step 7, 11, 6 |
