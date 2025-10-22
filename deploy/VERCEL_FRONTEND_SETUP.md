# Настройка Frontend на Vercel

## Быстрый старт

### 1. Деплой на Vercel

1. Перейдите на [vercel.com](https://vercel.com) и войдите через GitHub
2. Нажмите **New Project**
3. Выберите репозиторий `phd-students-portal`
4. Настройте проект:
   - **Framework Preset**: Vite
   - **Root Directory**: `frontend`
   - **Build Command**: `npm run build` (по умолчанию)
   - **Output Directory**: `dist` (по умолчанию)

### 2. Добавьте переменные окружения в Vercel

В настройках проекта → **Settings** → **Environment Variables**:

| Variable Name  | Value                                              | Environments                     |
| -------------- | -------------------------------------------------- | -------------------------------- |
| `VITE_API_URL` | `https://<ваш-railway-backend>.up.railway.app/api` | Production, Preview, Development |

**Где взять Railway backend URL:**

1. Откройте ваш проект на Railway
2. Выберите сервис backend
3. Во вкладке **Settings** → **Networking** → скопируйте **Public Domain**
4. Добавьте `/api` в конец URL (например: `https://phd-backend-production-abc123.up.railway.app/api`)

### 3. Настройте CORS на Backend (Railway)

В Railway → ваш backend сервис → **Variables** добавьте/обновите:

| Variable        | Value                                                                   |
| --------------- | ----------------------------------------------------------------------- |
| `FRONTEND_BASE` | `https://phd-students-portal.vercel.app` (замените на ваш Vercel домен) |

**Где взять Vercel URL:**

- После деплоя Vercel покажет URL вида `https://<project-name>.vercel.app`
- Скопируйте этот URL и вставьте в `FRONTEND_BASE` на Railway

### 4. Redeploy (если нужно)

- **Railway**: после изменения `FRONTEND_BASE` сервис перезапустится автоматически
- **Vercel**: после добавления `VITE_API_URL` нажмите **Redeploy** последнего деплоя

## Проверка

1. Откройте `https://<ваш-проект>.vercel.app/login`
2. Попробуйте войти (admin/admin123 или докторант/докторант123)
3. В DevTools → Network должны быть успешные запросы к `https://<railway-backend>/api/auth/login`
4. Не должно быть CORS ошибок

## Troubleshooting

### CORS ошибка: "No 'Access-Control-Allow-Origin' header"

**Причина**: Railway backend не знает Vercel домен.

**Решение**:

1. Проверьте, что `FRONTEND_BASE` на Railway содержит **точный** Vercel URL (без `/` в конце)
2. Перезапустите backend на Railway
3. Очистите кеш браузера и перезагрузите страницу

### Frontend подключается к localhost вместо Railway

**Причина**: `VITE_API_URL` не задан в Vercel.

**Решение**:

1. Vercel → Settings → Environment Variables
2. Добавьте `VITE_API_URL` со значением Railway backend URL + `/api`
3. Redeploy на Vercel

### 404 на /login или других страницах при перезагрузке

**Причина**: SPA routing не настроен в Vercel.

**Решение**: создайте `frontend/vercel.json`:

```json
{
  "rewrites": [
    {
      "source": "/(.*)",
      "destination": "/index.html"
    }
  ]
}
```

Закоммитьте, запушьте — Vercel автоматически пересоберёт.

## Custom Domain (опционально)

Если хотите использовать свой домен (например, `phd.kaznmu.edu.kz`):

1. Vercel → Settings → Domains → Add domain
2. Следуйте инструкциям Vercel для добавления DNS записей
3. После активации домена обновите `FRONTEND_BASE` на Railway на новый домен

## Vercel CLI (для продвинутых)

```bash
# Установите Vercel CLI
npm install -g vercel

# Войдите
vercel login

# Деплой из папки frontend
cd frontend
vercel

# Production деплой
vercel --prod
```

## Дополнительно

- **Preview Deployments**: Vercel автоматически создаёт preview для каждого PR
- **Analytics**: включите Web Analytics в Settings → Analytics
- **Monitoring**: проверяйте логи в Deployments → Logs

## Следующий шаг

После успешного деплоя frontend переходите к настройке production-окружения на серверах университета (см. `UNIVERSITY_IT_REQUIREMENTS.md`).
