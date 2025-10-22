# Temporary CORS Fix (if Railway env vars don't work)

Если Railway не применяет переменные окружения, можно временно захардкодить значение:

## В backend/internal/config/config.go

Замените строку:

```go
FrontendBase:  get("FRONTEND_BASE", "http://localhost:5173"),
```

На:

```go
FrontendBase:  getOrOverride("FRONTEND_BASE", "https://phd-students-portal.vercel.app"),
```

И добавьте функцию:

```go
func getOrOverride(k, override string) string {
	if v := os.Getenv(k); v != "" && v != "http://localhost:5173" {
		return v
	}
	// Force production value if env is release
	if os.Getenv("GIN_MODE") == "release" {
		return override
	}
	return "http://localhost:5173"
}
```

**После исправления:**

1. Коммит и push
2. Дождитесь деплоя
3. Проверьте `/api/debug/cors` — должно показать правильный URL
4. Логин должен работать

**Важно:** Это временное решение для диагностики. После того как найдём причину, вернём нормальную логику через env vars.

## Проверка Railway Variables через Railway CLI

Если у вас установлен Railway CLI:

```bash
railway login
railway link  # выберите ваш проект и backend service
railway variables
```

Должно показать все переменные, включая `FRONTEND_BASE`.

Если переменной нет в выводе:

```bash
railway variables set FRONTEND_BASE=https://phd-students-portal.vercel.app
```

## Проверка через Railway API

Можно проверить, какие переменные видит контейнер:

1. Railway → Backend → Settings → General → скопируйте Service ID
2. Используйте Railway GraphQL API для получения переменных (требуется токен)

Или проще: добавьте временный endpoint в `api.go`:

```go
api.GET("/debug/env", func(c *gin.Context) {
    c.JSON(200, gin.H{
        "FRONTEND_BASE": os.Getenv("FRONTEND_BASE"),
        "PORT":          os.Getenv("PORT"),
        "APP_PORT":      os.Getenv("APP_PORT"),
        "GIN_MODE":      os.Getenv("GIN_MODE"),
    })
})
```

Откройте `https://<railway-url>/api/debug/env` и увидите все переменные напрямую.
