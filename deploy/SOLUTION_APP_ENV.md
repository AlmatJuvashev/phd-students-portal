# 🎯 РЕШЕНИЕ: Railway показывает Env=development

## Проблема

Логи Railway показывают:
```
Config loaded: Port=8080, Env=development, FrontendBase=http://localhost:5173
```

**Причина:** Переменная `APP_ENV` не задана, поэтому код использует дефолтное значение `"development"`, и `FRONTEND_BASE` игнорируется в пользу localhost.

## ✅ Быстрое решение

### Добавьте в Railway Variables:

1. Railway → Backend Service → **Variables**
2. Нажмите **New Variable** и добавьте:

| Variable Name | Value |
|--------------|-------|
| `APP_ENV` | `production` |

3. Сервис автоматически перезапустится (~30 сек)

### Проверьте логи после перезапуска:

Должно быть:
```
Config loaded: Port=8080, Env=production, FrontendBase=https://phd-students-portal.vercel.app
```

### Проверьте /api/debug/cors:

```
https://phd-students-portal-production.up.railway.app/api/debug/cors
```

Должно показать:
```json
{
  "frontend_base": "https://phd-students-portal.vercel.app",
  "origin": ""
}
```

### Проверьте логин на Vercel:

```
https://phd-students-portal.vercel.app/login
```

Войдите (admin/admin123) → не должно быть CORS ошибок ✅

---

## 📋 Итоговые Railway Variables (полный список)

```bash
APP_ENV=production              # ← ДОБАВЬТЕ ЭТО!
APP_PORT=8080
DATABASE_URL=${{Postgres.DATABASE_URL}}
FRONTEND_BASE=https://phd-students-portal.vercel.app
GIN_MODE=release
JWT_SECRET=super-secret-jwt-key-change-this-now
ADMIN_EMAIL=juvashev@gmail.com
ADMIN_PASSWORD=<ваш-пароль>
```

**Удалите старые:**
- `CORS_ORIGINS` (больше не используется)
- `PORT` (используйте `APP_PORT` вместо этого)

---

## Почему это произошло?

В `config.go` код читает:
```go
Env: get("APP_ENV", "development"),  // ← дефолт = "development"
```

Когда `APP_ENV` не задана, код думает, что это dev окружение, и использует localhost для CORS.

---

## После исправления

1. ✅ Railway логи покажут `Env=production`
2. ✅ CORS будет разрешён для Vercel домена
3. ✅ Логин заработает без ошибок
4. ✅ Можете удалить debug endpoint `/api/debug/cors` из production кода

---

**Это всё!** После добавления `APP_ENV=production` всё заработает. 🚀
