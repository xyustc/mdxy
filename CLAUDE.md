# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Local dev (root)
```sh
./dev.sh {start|stop|restart|status|logs}
# Env overrides: FRONTEND_HOST, FRONTEND_PORT
```

### Frontend (`frontend/`)
```sh
npm install
npm run dev          # Vite on :5173, proxies /api to :8080
npm run build
npm run type-check   # lint gate (vue-tsc --noEmit)
```

### Backend (`backend-go/`)
```sh
go mod download
go run cmd/server/main.go
go build -o bin/server cmd/server/main.go
go build ./...       # compile check without binary output
go test ./...
go test ./internal/service -run TestName -v   # single test
```

### Docker
```sh
docker-compose up -d   # v2 stack (root docker-compose.yml)
```

No `*_test.go` or frontend test files exist yet.

### Verification after changes
```sh
cd backend-go && go build ./...              # backend compile
cd frontend && npm run type-check            # frontend types
```

## Architecture

Full-stack personal website: Go backend + Vue 3 SPA.

**Backend** (`backend-go/`): Gin + GORM + SQLite + JWT
Boot flow: `config` load → SQLite init + auto-migrate → IP geolocation init → router setup (`internal/router/router.go`)
Layering: **handler → service → repository → model** (consistent throughout)
API base: `/api/v1` — public routes + JWT-protected admin routes under `/api/v1/admin`

**Notes** are file-backed (from `notes/` directory, configured via `content.notes_dir`), not DB-backed. `NoteService` scans the filesystem. Featured state is stored in `notes-meta.json` (sibling to notes dir), not in the .md files. The meta file uses atomic writes (tmp + rename) and mutex protection.

**Analytics** are middleware-driven — `middleware.Logger` asynchronously batches `AccessLog` writes. Visitor identity uses `X-Visitor-ID` header, then cookie, then IP/UA hash fallback.

**Frontend** (`frontend/src`): Vue 3 + TypeScript + Vite
- Router splits `ClientLayout` (public) and `AdminLayout` (admin)
- Axios wrapper (`api/request.ts`) injects `Authorization` and `X-Visitor-ID` on every request; auto-redirects to `/admin/login` on 401
- API modules: `src/api/*.ts`; shared types: `src/api/types.ts`
- UI: Naive UI (public pages), Element Plus (admin pages)
- Design system: CSS tokens/themes in `src/styles/`
- Admin token stored in `localStorage` as `admin_token`

**Homepage data** is cached in `localStorage` with TTL (`stores/site.ts`). Admin pages that modify featured/visible state must call `siteStore.clearHomeCache()` to invalidate.

## Key conventions

**Response envelope**: all backend responses use `{ success: true, data }` or `{ success: false, error }` via `internal/pkg/response`. Helpers: `Success`, `BadRequest`, `Unauthorized`, `Forbidden`, `NotFound`, `InternalServerError`, `TooManyRequests`. Frontend `ApiResponse<T>` assumes this shape.

**New backend feature order**: model → repository → service → handler → register in `router.go`

**Notes path contract**: handler trims leading slash from `*path` param; `validateNotePath` enforces `.md` extension + path traversal protection; `validatePath` is the base (no `.md` requirement, used by `CreateDirectory`).

**Visitor ID**: always propagate `X-Visitor-ID` when adding new clients or request wrappers — it's required for analytics quality.

**Frontend style**: `<script setup lang="ts">`, `@` alias for `src`, scoped styles using shared CSS variables.

**Shared utilities**: `@/utils/highlight.ts` exports both `hljs` (configured instance) and `markdownHighlight` (markdown-it callback). Use `markdownHighlight` when creating MarkdownIt instances to avoid duplicating the highlight function.
