# 🗄️ Railway Database Migrations Guide

## Проблема

Вы видите ошибку:
```
stat cmd/migrate/main.go: no such file or directory
```

Это потому что в проекте используется `golang-migrate` CLI tool, а не собственный файл миграций.

---

## ✅ Решение: 3 способа запустить миграции

### Способ 1: Автоматически при деплое (Рекомендуется)

**Уже настроено!** После коммита обновлённых файлов, миграции будут запускаться автоматически при каждом деплое.

```bash
# Просто запушьте изменения
git add .
git commit -m "Add automatic migrations"
git push origin main
```

Railway запустит:
1. `release` команду → миграции
2. `web` команду → запуск сервера

---

### Способ 2: Railway CLI (Ручной запуск)

#### Шаг 1: Установите golang-migrate

**macOS:**
```bash
brew install golang-migrate
```

**Linux:**
```bash
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/
```

**Windows:**
```bash
choco install golang-migrate
```

Или через Go:
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

#### Шаг 2: Получите DATABASE_URL из Railway

1. Railway → Ваш backend service → **Variables**
2. Скопируйте значение `DATABASE_URL`
   ```
   postgresql://postgres:password@region.railway.app:5432/railway
   ```

#### Шаг 3: Запустите миграции локально

```bash
cd backend
migrate -database "postgresql://postgres:password@region.railway.app:5432/railway" -path db/migrations up
```

**Результат:**
```
0001_init.up.sql (2024/10/13 10:30:00)
0002_comments.up.sql (2024/10/13 10:30:01)
0003_add_graduation_year.up.sql (2024/10/13 10:30:02)
0004_add_playbook_locked.up.sql (2024/10/13 10:30:03)
```

---

### Способ 3: Прямое SQL подключение

#### Шаг 1: Установите psql

**macOS:**
```bash
brew install postgresql
```

**Ubuntu/Debian:**
```bash
sudo apt-get install postgresql-client
```

#### Шаг 2: Подключитесь к базе

```bash
psql "postgresql://postgres:password@region.railway.app:5432/railway"
```

#### Шаг 3: Запустите миграции вручную

```bash
# Из папки backend
psql "your-database-url" < db/migrations/0001_init.up.sql
psql "your-database-url" < db/migrations/0002_comments.up.sql
psql "your-database-url" < db/migrations/0003_add_graduation_year.up.sql
psql "your-database-url" < db/migrations/0004_add_playbook_locked.up.sql
```

---

## 🧪 Проверка успешности

### 1. Проверьте таблицы

```bash
psql "your-database-url" -c "\dt"
```

**Ожидаемый результат:**
```
           List of relations
 Schema |        Name         | Type  
--------+---------------------+-------
 public | comments            | table
 public | journey_states      | table
 public | node_submissions    | table
 public | playbooks           | table
 public | schema_migrations   | table
 public | users               | table
```

### 2. Проверьте backend логи

Railway → Backend service → Deployments → View Logs

Ищите:
```
✅ All migrations applied successfully!
```

### 3. Проверьте API

Откройте в браузере:
```
https://your-backend.railway.app/api/health
```

Должно вернуть:
```json
{"status":"healthy"}
```

---

## 🚨 Troubleshooting

### "dirty database version"

**Проблема:** Миграция была прервана посередине.

**Решение:**
```bash
# Узнайте текущую версию
migrate -database "your-db-url" -path db/migrations version

# Если версия "dirty" (например, 2), очистите её
migrate -database "your-db-url" -path db/migrations force 2

# Попробуйте снова
migrate -database "your-db-url" -path db/migrations up
```

---

### "no change" или миграции уже применены

**Это нормально!** Миграции идемпотентны - безопасно запускать повторно.

---

### "permission denied"

**Проблема:** Недостаточно прав у пользователя БД.

**Решение:** Убедитесь что DATABASE_URL использует правильного пользователя с правами CREATE TABLE.

---

## 📋 Что изменено в проекте

1. **`Procfile`** - Добавлена `release` команда для автоматических миграций
2. **`backend/scripts/migrate.sh`** - Скрипт запуска миграций
3. **`nixpacks.toml`** - Установка golang-migrate при сборке

---

## 🔄 После успешных миграций

Вернитесь к **`DEPLOY_CHECKLIST.md`** и продолжите с **Step 3: Deploy Frontend**.

---

**Создано:** 13 октября 2025  
**Обновлено:** Исправлена команда миграций для Railway
