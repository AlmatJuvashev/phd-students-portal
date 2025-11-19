# Quick Start: S3 File Upload Testing

## 🚀 Быстрый запуск

### 1. Запустить инфраструктуру

```bash
# Запустить PostgreSQL и MinIO
docker-compose up -d

# Создать bucket и настроить CORS
./scripts/setup-minio.sh
```

### 2. Запустить приложение

```bash
# Терминал 1: Backend
cd backend && go run ./cmd/server

# Терминал 2: Frontend
cd frontend && npm run dev
```

### 3. Протестировать загрузку

**Как студент:**

1. Войти: http://localhost:5173

   - Email: `student@kaznmu.kz`
   - Password: `password`

2. Journey → Найти узел **"Антиплагиат ≥ 85%"**

3. Загрузить PDF или DOCX файл (макс. 100MB)

**Как admin:**

1. Войти: http://localhost:5173

   - Email: (see `ADMIN_EMAIL` in backend/.env)
   - Password: (see `ADMIN_PASSWORD` in backend/.env)

2. Admin → Students → Выбрать студента

3. Просмотреть загруженные файлы в разделе Journey

## 📊 MinIO Console

- URL: http://localhost:9091
- Username: `minioadmin`
- Password: `minioadmin`

## ✅ Что работает

- ✅ Загрузка файлов студентами (PDF/DOC/DOCX)
- ✅ Валидация размера (макс. 100MB)
- ✅ Валидация типа файла
- ✅ Просмотр файлов advisor/admin
- ✅ Скачивание файлов
- ✅ Локальное S3-хранилище (MinIO)

## 📖 Подробная документация

- [S3 Storage Guide](./S3_STORAGE.md) - Полное руководство по S3
- [Testing File Upload](./TESTING_FILE_UPLOAD.md) - Детальная инструкция по тестированию

## 🔧 Настройки

Файл `backend/.env`:

```properties
S3_ENDPOINT=http://localhost:9090
S3_BUCKET=phd-portal
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_MAX_FILE_SIZE_MB=100
S3_PRESIGN_EXPIRES_MINUTES=15
```

## 🐛 Устранение проблем

**MinIO не запущен:**

```bash
docker-compose restart minio
```

**Bucket не создан:**

```bash
./scripts/setup-minio.sh
```

**Backend не подключается к MinIO:**

- Проверьте, что MinIO работает: `docker-compose ps`
- Проверьте логи: `docker-compose logs minio`
- Проверьте `.env` файл

**CORS ошибки:**

```bash
mc alias set local http://localhost:9090 minioadmin minioadmin
mc anonymous set download local/phd-portal
```
