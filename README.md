# PhD Student Portal — Monorepo

This repository contains:
- `backend/` — Go + Gin + sqlx + Postgres + JWT
- `frontend/` — React Router v7 + TanStack Query + Tailwind
- `mailserver/` — Mailpit (SMTP + Web UI)

## Quickstart

1. **Mailpit (emails)**
   ```bash
   cd mailserver && docker compose up -d
   ```
2. **Database**
   - Provide a Postgres instance and set `DATABASE_URL` in `backend/.env` (copy from `.env.example`).
3. **Backend**
   ```bash
   cd backend
   make migrate-up
   make run
   ```
4. **Frontend**
   ```bash
   cd ../frontend
   npm i
   VITE_API_URL=http://localhost:8080/api npm run dev
   ```

## Authentication
- Email + password login.
- JWT expiry ~ 6 months (configurable via `JWT_EXP_DAYS`).
- Password reset flow via Mailpit.

## Admin
- Create users with auto username & temp password (copy-once).
- Soft remove via `is_active` flag.
- Admin can change others’ passwords except for **superadmin**.


## v4 Upgrades
- Auto S3/local upload detection with pre-signed PUT
- Threaded comments with @mentions
- Minimal shadcn-style components + Framer Motion polish
- Backend user listing endpoint for mentions


## v5 Upgrades
- Role-based route guards (JWT-decoded role) with TanStack Router `beforeLoad`
- Vendored shadcn/ui-style components and theme tokens (mini design system)


## v6 Upgrades
- Backend now on **8280**; root `docker-compose.yml` spins up Postgres, Redis, Mailpit, MinIO, and the backend.
- New `/api/me` endpoint; Redis-backed caching for user context.
- Structured logging helpers and more comments.
- Frontend role-aware top navigation; `/me`-driven auth; common folders for hooks/config/lib.
- Added toast system; forms use `react-hook-form` + `zod`.


## v6 Upgrades
- Role-aware top nav using `/me` (no client-side JWT decode)
- `/me` endpoint + user hydration into request context (Redis-cached)
- Structured logs middleware
- Redis service + caching helpers
- Toast system + react-hook-form + zod
- Mobile vertical progress bar for students
- Root docker-compose with Postgres, Redis, Mailpit, Backend (8280), Frontend (5173)
