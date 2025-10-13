# KazNMU PhD Portal

Система управления процессом докторантуры для Казахского Национального Медицинского Университета имени С.Д. Асфендиярова.

## 🎯 Описание

Веб-приложение для автоматизации документооборота и отслеживания прогресса докторантов PhD программы, включая:

- 📋 Интерактивная карта пути докторанта (Journey Map)
- 📝 Управление документами и чеклистами
- 👥 Роли: Докторант, Научный руководитель, Администратор
- 🔔 Комментарии и уведомления
- 📊 Отчёты и статистика (для администраторов)
- 📱 Адаптивный дизайн (mobile-friendly)

## 🛠️ Технологии

**Frontend:** React 18 + TypeScript + Vite + TailwindCSS + Radix UI  
**Backend:** Go 1.21+ + Gin + PostgreSQL 14+ + JWT  
**Deploy:** Docker Compose, Railway, или традиционный сервер

---

## 🚀 Quickstart (Local Development)

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

---

## 📦 Deployment Options

### 🎬 Demo Version (Railway + Vercel) - 15 minutes

Quick deploy for demonstration:

```bash
# See detailed instructions
cat deploy/QUICK_DEPLOY.md
```

**Steps:**
1. Push to GitHub
2. Deploy backend to Railway (with PostgreSQL)
3. Deploy frontend to Vercel
4. Done! 🎉

### 🏛️ Production Version (University Server)

Full documentation for KazNMU IT department:

```bash
# Server requirements and deployment guide
cat deploy/UNIVERSITY_IT_REQUIREMENTS.md
```

**Integration with https://kaznmu.edu.kz:**
- Option A: Subdomain `phd.kaznmu.edu.kz` (recommended)
- Option B: Path `/phd-portal/` on main domain
- Option C: iFrame integration

---

## 📚 Documentation

- [Quick Deploy (Railway)](deploy/QUICK_DEPLOY.md)
- [IT Requirements](deploy/UNIVERSITY_IT_REQUIREMENTS.md)
- [Full Deployment Guide](deploy/DEPLOYMENT_GUIDE.md)
- [Backend API](backend/README.md)

---

**Version:** 1.0  
**Status:** Ready for Production  
**Date:** October 2025
