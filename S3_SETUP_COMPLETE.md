# ✅ S3 File Upload - Готово к тестированию!

## 🎯 Что сделано

### 1. ✅ MinIO настроен и работает

- **S3 API:** http://localhost:9000
- **Console:** http://localhost:9001 (minioadmin / minioadmin)
- **Bucket:** phd-portal создан
- **CORS:** настроен для localhost:5173

### 2. ✅ Backend настроен

- **Порт:** 8280
- **S3 Client:** подключен к MinIO
- **Валидация файлов:**
  - Размер: макс. 100MB
  - Типы: PDF, DOC, DOCX
- **Endpoints готовы:**
  - `POST /api/journey/nodes/:nodeId/uploads/presign` - получить presigned URL
  - `POST /api/journey/nodes/:nodeId/uploads/attach` - подтвердить загрузку
  - `GET /api/admin/students/:id/nodes/:nodeId/files` - список файлов

### 3. ✅ Playbook настроен

Узел **S1_antiplag** (Антиплагиат):

```json
{
  "id": "S1_antiplag",
  "title": "Антиплагиат ≥ 85%",
  "requirements": {
    "uploads": [
      {
        "key": "antiplag_report",
        "label": "Антиплагиат: отчёт/документ (PDF/DOCX)",
        "required": true,
        "mime": [
          "application/pdf",
          "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
        ]
      }
    ]
  }
}
```

### 4. ✅ Доступы готовы

**Admin:**

- Email: (see `ADMIN_EMAIL` in backend/.env)
- Password: (see `ADMIN_PASSWORD` in backend/.env)
- Может просматривать файлы всех студентов

**Student:**

- Email: `student@kaznmu.kz`
- Password: `password`
- Может загружать файлы в свой journey

## 🚀 Как тестировать

### Шаг 1: Убедитесь, что всё запущено

```bash
# Проверить Docker
docker-compose ps

# Должно быть:
# ✅ postgres - Up
# ✅ minio - Up
```

### Шаг 2: Запустите приложение

```bash
# Терминал 1: Backend (уже запущен)
# cd backend && go run ./cmd/server

# Терминал 2: Frontend
cd frontend
npm run dev
```

### Шаг 3: Тест как студент

1. Откройте http://localhost:5173
2. Войдите как студент
3. Перейдите в Journey
4. Найдите узел "Антиплагиат ≥ 85%"
5. Загрузите PDF или DOCX файл
6. Проверьте, что файл появился в списке вложений

### Шаг 4: Тест как admin

1. Выйдите и войдите как admin
2. Admin → Monitor Students
3. Выберите студента
4. Просмотрите его Journey
5. Найдите узел S1_antiplag
6. Проверьте, что видны загруженные файлы
7. Скачайте файл

### Шаг 5: Проверка в MinIO Console

1. Откройте http://localhost:9001
2. Войдите: minioadmin / minioadmin
3. Buckets → phd-portal
4. Проверьте структуру папок:
   ```
   nodes/
     └── {userID}/
         └── S1_antiplag/
             └── antiplag_report/
                 └── 2025-11-15-{uuid}-filename.pdf
   ```

## 📊 Архитектура загрузки

```
Frontend (Student)
    ↓
    1. POST /api/journey/nodes/S1_antiplag/uploads/presign
    ↓
Backend
    ↓
    2. Генерирует Presigned URL (valid 15 min)
    ↓
Frontend
    ↓
    3. PUT {presigned_url} (Direct upload to MinIO)
    ↓
MinIO S3
    ↓
    4. Сохраняет файл
    ↓
Frontend
    ↓
    5. POST /api/journey/nodes/S1_antiplag/uploads/attach
    ↓
Backend
    ↓
    6. Сохраняет метаданные в PostgreSQL
```

## 🔍 Проверка логов

### Backend logs

Должны быть видны:

```
[S3] Presigning PUT: bucket=phd-portal key=nodes/123/S1_antiplag/... expires=15m0s
```

### Frontend DevTools

Network tab должен показать:

1. POST `/api/journey/nodes/S1_antiplag/uploads/presign` → 200 OK
2. PUT `http://localhost:9000/phd-portal/nodes/...` → 200 OK
3. POST `/api/journey/nodes/S1_antiplag/uploads/attach` → 201 Created

## 🎉 Ожидаемый результат

✅ Student может загружать PDF/DOCX файлы  
✅ Файлы сохраняются в MinIO  
✅ Метаданные сохраняются в PostgreSQL  
✅ Admin видит все загруженные файлы  
✅ Admin может скачать файлы  
✅ Advisor может оставить комментарий (если UI поддерживает)

## 📖 Документация

- **[QUICK_START_S3.md](./QUICK_START_S3.md)** - Быстрый старт
- **[docs/TESTING_FILE_UPLOAD.md](./docs/TESTING_FILE_UPLOAD.md)** - Подробное тестирование
- **[docs/S3_STORAGE.md](./docs/S3_STORAGE.md)** - Полное руководство по S3

## 🛠️ Конфигурация

**Docker Compose:** MinIO на портах 9090 (API) / 9091 (Console)  
**Backend .env:** S3 credentials для MinIO  
**Playbook:** Узел S1_antiplag с поддержкой uploads  
**Frontend:** API client настроен на localhost:8280

---

**Всё готово для тестирования загрузки файлов! 🚀**

Откройте http://localhost:5173 и начните тестировать!
