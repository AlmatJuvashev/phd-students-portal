# 🚀 Быстрый чеклист деплоя (Railway + Vercel)

## ✅ Что сделать прямо сейчас

### 1. Получите Railway Backend URL

```bash
# В Railway UI:
Project → Backend Service → Settings → Networking → Public Domain
# Пример: https://phd-backend-production-abc123.up.railway.app
```

**Скопируйте этот URL** — он понадобится на следующих шагах.

### 2. Добавьте переменную в Vercel

1. Откройте ваш проект на [vercel.com](https://vercel.com/dashboard)
2. Settings → Environment Variables → Add New
3. Заполните:
   - **Name**: `VITE_API_URL`
   - **Value**: `<Railway URL из шага 1>/api` (например: `https://phd-backend-production-abc123.up.railway.app/api`)
   - **Environments**: отметьте все (Production, Preview, Development)
4. Нажмите **Save**

### 3. Redeploy на Vercel

1. Deployments → последний деплой → три точки (⋯) → **Redeploy**
2. Дождитесь завершения (1-2 минуты)
3. Скопируйте **Vercel URL** (например: `https://phd-students-portal.vercel.app`)

### 4. Добавьте FRONTEND_BASE в Railway

1. Railway → Backend Service → Variables
2. Добавьте/обновите:
   - **Name**: `FRONTEND_BASE`
   - **Value**: `<Vercel URL из шага 3>` (например: `https://phd-students-portal.vercel.app`)
3. Сервис автоматически перезапустится (~30 сек)

### 5. Проверьте работу

1. Откройте Vercel URL в браузере: `https://<ваш-проект>.vercel.app/login`
2. Откройте DevTools (F12) → Console
3. Попробуйте войти (admin/admin123)
4. **Должно быть:**
   - ✅ Нет CORS ошибок
   - ✅ Запросы идут на Railway backend URL
   - ✅ Успешный логин и редирект на `/`

## 🔍 Если что-то не работает

### CORS ошибка ("No 'Access-Control-Allow-Origin'")

- **Проверьте**: Railway Variables → `FRONTEND_BASE` точно совпадает с Vercel URL (без `/` в конце)
- **Решение**: обновите `FRONTEND_BASE`, подождите перезапуск (~30 сек), обновите страницу в браузере

### Frontend подключается к localhost

- **Проверьте**: Vercel → Settings → Environment Variables → есть ли `VITE_API_URL`?
- **Решение**: добавьте переменную, redeploy на Vercel

### 404 при перезагрузке страницы (/login, /journey и т.д.)

- **Проверьте**: есть ли файл `frontend/vercel.json` в репозитории?
- **Решение**: он уже создан в этом коммите, запушьте и redeploy

### Backend не запускается / миграции падают

- **Проверьте**: Railway → Backend → Logs — есть ли ошибки?
- **Частая причина**: `DATABASE_URL` не задан — добавьте Postgres Plugin в проект
- **Решение**: Railway → Add Plugin → PostgreSQL → Connect to Backend

## 📋 Итоговая конфигурация

**Railway Backend Variables:**
```bash
DATABASE_URL=postgresql://...  # auto from Postgres plugin
JWT_SECRET=<secure-random-string>
PORT=8280  # or leave default
GIN_MODE=release
FRONTEND_BASE=https://phd-students-portal.vercel.app
```

**Vercel Frontend Variables:**
```bash
VITE_API_URL=https://<railway-backend>.up.railway.app/api
```

## 🎯 Следующие шаги

1. ✅ Проверьте логин/регистрацию
2. ✅ Откройте карту (`/journey`) — должна грузиться из бэкенда
3. ✅ Создайте тестового пользователя через Админ-панель (`/admin/users`)
4. 📝 Подготовьте документацию для IT-отдела университета (см. `UNIVERSITY_IT_REQUIREMENTS.md`)
5. 🔐 Настройте production secrets (JWT_SECRET, DB пароли)
6. 📧 Настройте SMTP для email-уведомлений (опционально)
7. 📦 Настройте S3 для загрузки файлов (опционально)

## ⚡ Быстрые команды (если нужно локально проверить)

```bash
# Backend (локально)
cd backend
export DATABASE_URL="postgresql://localhost:5432/phd_portal"
export JWT_SECRET="dev-secret"
export FRONTEND_BASE="http://localhost:5173"
make run

# Frontend (локально)
cd frontend
export VITE_API_URL="http://localhost:8280/api"
npm run dev
```

## 🆘 Нужна помощь?

- Railway логи: Service → Deployments → View Logs
- Vercel логи: Deployments → кликните на деплой → Runtime Logs
- Frontend ошибки: DevTools (F12) → Console / Network

---

**Готово!** 🎉 Ваше приложение должно работать на Vercel + Railway.
