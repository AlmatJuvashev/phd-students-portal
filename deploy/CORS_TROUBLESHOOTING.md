# CORS Troubleshooting Guide

## Текущая проблема

```
Access to fetch at 'https://phd-students-portal-production.up.railway.app/api/auth/login'
from origin 'https://phd-students-portal.vercel.app' has been blocked by CORS policy
```

## Быстрая диагностика

### Шаг 1: Проверьте Railway Variables

1. Railway → Backend Service → **Variables**
2. Убедитесь, что `FRONTEND_BASE` существует и равен:
   ```
   https://phd-students-portal.vercel.app
   ```
   **Важно:** без `/` в конце!

### Шаг 2: Проверьте конфигурацию через debug endpoint

После деплоя коммита `709ed29` (или позже), откройте в браузере:

```
https://phd-students-portal-production.up.railway.app/api/debug/cors
```

**Должно вернуть:**

```json
{
  "frontend_base": "https://phd-students-portal.vercel.app",
  "origin": "https://phd-students-portal.vercel.app"
}
```

**Если `frontend_base` показывает другое значение:**

- Проверьте Railway Variables → переменная `FRONTEND_BASE` задана правильно
- Redeploy сервиса (Deployments → последний → Redeploy)

**Если `origin` пустой:**

- Вы открыли endpoint напрямую в браузере (это ОК)
- Попробуйте сделать fetch из DevTools:
  ```javascript
  fetch(
    "https://phd-students-portal-production.up.railway.app/api/debug/cors",
    {
      headers: { Origin: "https://phd-students-portal.vercel.app" },
    }
  )
    .then((r) => r.json())
    .then(console.log);
  ```

### Шаг 3: Проверьте деплой на Railway

1. Railway → Backend → **Deployments**
2. Последний деплой должен быть с коммитом `709ed29` или новее
3. Статус должен быть **Success** (зелёная галочка)
4. Проверьте логи:
   - Должны быть успешные миграции: `✅ All migrations applied successfully!`
   - Сервер должен запуститься: `Listening on :8280` или подобное

### Шаг 4: Проверьте preflight запрос

В DevTools на Vercel странице → Network → включите фильтр "All":

1. Попробуйте войти
2. Найдите запрос `OPTIONS /api/auth/login`
3. Проверьте Response Headers:
   - `Access-Control-Allow-Origin: https://phd-students-portal.vercel.app`
   - `Access-Control-Allow-Methods: GET, POST, PUT, PATCH, DELETE, OPTIONS`
   - `Access-Control-Allow-Headers: Origin, Content-Type, Authorization, Accept`
   - `Access-Control-Allow-Credentials: true`

**Если хедеры отсутствуют:**

- Backend не получил переменную `FRONTEND_BASE`
- Или запрос не дошёл до backend (проблема с Railway)

## Возможные причины и решения

### 1. Railway не видит FRONTEND_BASE

**Симптомы:**

- Debug endpoint показывает `frontend_base: "http://localhost:5173"`
- Или показывает пустую строку

**Решение:**

1. Railway → Variables → убедитесь `FRONTEND_BASE` добавлена
2. Если есть, удалите и добавьте заново
3. Redeploy сервиса

### 2. Старый код на Railway

**Симптомы:**

- Debug endpoint не существует (404)
- CORS хедеры не включают `Access-Control-Allow-Credentials: true`

**Решение:**

1. Проверьте, что Railway подключён к правильному репозиторию
2. Settings → GitHub → Branch должна быть `main`
3. Deployments → последний коммит должен быть `709ed29` или новее
4. Если нет — нажмите "New Deployment" и выберите последний коммит

### 3. Несоответствие доменов

**Симптомы:**

- Vercel URL изменился (например, добавили custom domain)
- Frontend всё ещё использует старый Railway URL

**Решение:**

1. Проверьте актуальный Vercel URL в Settings → Domains
2. Обновите `FRONTEND_BASE` в Railway
3. Проверьте `VITE_API_URL` в Vercel Variables

### 4. Проблема с preflight OPTIONS

**Симптомы:**

- Запрос `OPTIONS` возвращает 404 или 500
- Нет хедеров `Access-Control-*`

**Решение:**
Это может быть проблема с middleware порядком. Проверьте `api.go`:

- CORS middleware должен быть **первым** после logger
- Порядок должен быть: logger → CORS → routes

### 5. Railway networking issue

**Симптомы:**

- Все конфиги правильные, но CORS всё равно блокируется
- Railway логи показывают успешный запуск

**Решение:**

1. Попробуйте сделать curl запрос напрямую:
   ```bash
   curl -X OPTIONS \
     https://phd-students-portal-production.up.railway.app/api/auth/login \
     -H "Origin: https://phd-students-portal.vercel.app" \
     -H "Access-Control-Request-Method: POST" \
     -H "Access-Control-Request-Headers: content-type" \
     -v
   ```
2. Проверьте response headers в выводе
3. Если хедеры есть — проблема на стороне браузера/Vercel
4. Если хедеров нет — проблема в Railway/backend

## Финальная проверка (контрольный список)

- [ ] Railway Variables → `FRONTEND_BASE=https://phd-students-portal.vercel.app`
- [ ] Railway Deployments → последний коммит `709ed29` или новее
- [ ] Railway Deployments → статус Success
- [ ] `/api/debug/cors` возвращает правильный `frontend_base`
- [ ] Vercel Variables → `VITE_API_URL=https://phd-students-portal-production.up.railway.app/api`
- [ ] Vercel последний деплой содержит `vercel.json`
- [ ] DevTools Network → OPTIONS запрос возвращает CORS headers

## Временное решение (если ничего не помогает)

Если срочно нужно запустить, можно временно разрешить все origins (только для dev/demo!):

**В `api.go` временно замените:**

```go
AllowOriginFunc: func(origin string) bool {
    return true  // WARNING: allow all origins (dev only!)
},
```

**После деплоя обязательно:**

1. Проверьте, что приложение работает
2. Найдите причину исходной проблемы
3. Верните правильную CORS логику

## Контакты для помощи

- Railway логи: Service → Deployments → View Logs
- Vercel логи: Deployments → Runtime Logs
- Backend код CORS: `backend/internal/handlers/api.go:20-45`

---

**После решения проблемы:** удалите debug endpoint `/api/debug/cors` из production кода!
