# Тестирование загрузки файлов в узел Antiplagiarism

## Обзор

Узел **S1_antiplag** (Антиплагиат) в playbook настроен с поддержкой загрузки файлов. Студенты могут загружать отчёты об антиплагиате в форматах PDF/DOCX, а advisor и admin могут просматривать эти файлы и оставлять комментарии.

## Настройка локального окружения

### 1. Запуск MinIO (локальное S3-хранилище)

MinIO уже настроен в `docker-compose.yml` и готов к использованию:

```bash
# Запустить PostgreSQL и MinIO
docker-compose up -d

# Создать bucket и настроить CORS
./scripts/setup-minio.sh
```

**Порты:**

- MinIO S3 API: `http://localhost:9000`
- MinIO Console: `http://localhost:9091`
  - Username: `minioadmin`
  - Password: `minioadmin`

### 2. Проверка конфигурации

Файл `backend/.env` должен содержать:

```properties
S3_ENDPOINT=http://localhost:9000
S3_REGION=us-east-1
S3_BUCKET=phd-portal
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_USE_PATH_STYLE=true
S3_PRESIGN_EXPIRES_MINUTES=15
S3_MAX_FILE_SIZE_MB=100
```

### 3. Запуск приложения

```bash
# Терминал 1: Backend
cd backend
go run ./cmd/server

# Терминал 2: Frontend
cd frontend
npm run dev
```

## Тестирование загрузки файлов

### Как студент (Student)

1. **Войдите как студент:**

   - Email: `student@kaznmu.kz` (или другой студент из базы)
   - Password: `password` (или установленный пароль)

2. **Перейдите на страницу Journey:**

   - Откройте `/student/journey`
   - Найдите узел **"S1_antiplag"** ("Антиплагиат ≥ 85%")

3. **Загрузите файл:**

   - Узел должен показывать поле для загрузки файла с меткой:

     - RU: "Антиплагиат: отчёт/документ (PDF/DOCX)"
     - KZ: "Антиплагиат: есеп/құжат (PDF/DOCX)"
     - EN: "Antiplagiarism: report/document (PDF/DOCX)"

   - Нажмите **"Загрузить файл"** или перетащите файл
   - Поддерживаемые форматы:

     - `application/pdf` (PDF)
     - `application/vnd.openxmlformats-officedocument.wordprocessingml.document` (DOCX)
     - `application/msword` (DOC)

   - Максимальный размер: **100 MB**

4. **Проверьте загрузку:**
   - Файл должен отобразиться в списке вложений
   - Должны быть видны:
     - Название файла
     - Размер файла
     - Дата загрузки
     - Кнопка "Скачать"

### Как advisor/admin

1. **Войдите как admin:**

   - Email: `admin@example.com` (see your `.env` file for actual credentials)
   - Password: (see `ADMIN_PASSWORD` in `.env`)

2. **Найдите студента:**

   - Перейдите в `/admin/monitor/students`
   - Выберите студента, который загрузил файл
   - Кликните на карточку студента

3. **Просмотр загруженных файлов:**

   - На странице деталей студента (`/admin/students/:id`)
   - Перейдите в раздел **"Journey"** или **"Nodes"**
   - Найдите узел **S1_antiplag**
   - Должны быть видны все загруженные файлы студента

4. **Скачайте файл:**

   - Кликните на кнопку **"Скачать"** рядом с файлом
   - Файл должен скачаться из MinIO

5. **Оставьте комментарий (опционально):**
   - Если в интерфейсе есть возможность комментирования
   - Добавьте feedback для студента

## Проверка в MinIO Console

1. Откройте MinIO Console: `http://localhost:9091`
2. Войдите с credentials: `minioadmin` / `minioadmin`
3. Перейдите в **Buckets** → **phd-portal**
4. Вы должны увидеть структуру:
   ```
   nodes/
     └── {userID}/
         └── {nodeID}/
             └── {slotKey}/
                 └── {timestamp}-{uuid}-{filename}
   ```

## Типичные проблемы и решения

### 1. Ошибка "File upload failed"

**Причина:** MinIO не запущен или bucket не создан

**Решение:**

```bash
docker-compose up -d
./scripts/setup-minio.sh
```

### 2. CORS ошибка в браузере

**Причина:** CORS не настроен для MinIO

**Решение:**

```bash
mc alias set local http://localhost:9000 minioadmin minioadmin
mc anonymous set download local/phd-portal
```

### 3. "File too large" ошибка

**Причина:** Файл превышает лимит 100MB

**Решение:** Используйте файл меньше 100MB или измените `S3_MAX_FILE_SIZE_MB` в `.env`

### 4. "Unsupported file type"

**Причина:** Формат файла не в whitelist

**Решение:** Используйте PDF, DOC или DOCX файл

### 5. Presigned URL expired

**Причина:** URL устарел (истёк срок действия 15 минут)

**Решение:** Запросите новый presigned URL или увеличьте `S3_PRESIGN_EXPIRES_MINUTES`

## API endpoints для тестирования

### 1. Получить presigned URL для загрузки

```bash
POST /api/journey/nodes/S1_antiplag/uploads/presign
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "filename": "antiplagiat_report.pdf",
  "content_type": "application/pdf",
  "size_bytes": 1048576,
  "slot_key": "antiplag_report"
}
```

**Ответ:**

```json
{
  "upload_url": "http://localhost:9000/phd-portal/nodes/123/S1_antiplag/antiplag_report/...",
  "object_key": "nodes/123/S1_antiplag/antiplag_report/2025-11-15-uuid-antiplagiat_report.pdf",
  "expires_at": "2025-11-15T10:15:00Z"
}
```

### 2. Загрузить файл напрямую в S3

```bash
PUT <upload_url>
Content-Type: application/pdf

<binary file data>
```

### 3. Подтвердить загрузку

```bash
POST /api/journey/nodes/S1_antiplag/uploads/attach
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "object_key": "nodes/123/S1_antiplag/antiplag_report/...",
  "slot_key": "antiplag_report",
  "filename": "antiplagiat_report.pdf",
  "size_bytes": 1048576,
  "content_type": "application/pdf",
  "etag": "\"d41d8cd98f00b204e9800998ecf8427e\""
}
```

### 4. Список файлов студента (для admin)

```bash
GET /api/admin/students/:id/nodes/S1_antiplag/files
Authorization: Bearer <admin_jwt_token>
```

## Структура данных в базе

### Таблица `node_attachments`

```sql
SELECT * FROM node_attachments WHERE node_id = 'S1_antiplag';
```

Поля:

- `id` - UUID вложения
- `user_id` - ID студента
- `node_id` - 'S1_antiplag'
- `slot_key` - 'antiplag_report'
- `object_key` - Полный путь в S3
- `filename` - Оригинальное имя файла
- `size_bytes` - Размер файла
- `content_type` - MIME type
- `uploaded_at` - Дата загрузки
- `reviewed_status` - 'pending', 'approved', 'rejected'
- `reviewer_comment` - Комментарий от advisor/admin

## Production deployment

Для продакшена замените MinIO на AWS S3:

```properties
S3_ENDPOINT=
S3_REGION=eu-west-1
S3_BUCKET=kaznmu-phd-portal-prod
S3_ACCESS_KEY=<aws_access_key>
S3_SECRET_KEY=<aws_secret_key>
S3_USE_PATH_STYLE=false
```

См. полную документацию в `docs/S3_STORAGE.md`.

## Мониторинг

В MinIO Console вы можете:

- Просматривать все загруженные файлы
- Проверять размер bucket
- Настраивать lifecycle policies
- Анализировать метрики загрузки

**URL:** http://localhost:9091
