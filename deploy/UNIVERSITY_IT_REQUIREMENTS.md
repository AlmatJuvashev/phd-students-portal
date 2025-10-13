# Требования для IT отдела КазНМУ
# PhD Portal - Развёртывание на серверах университета

## 📋 Краткое описание

**Система:** Портал для управления процессом докторантуры
**Технологии:** Go (backend) + React (frontend) + PostgreSQL
**Цель:** Автоматизация документооборота для докторантов PhD

---

## 🖥️ Минимальные требования к серверу

### Аппаратные требования:
- **CPU:** 2+ ядра (рекомендуется 4)
- **RAM:** 4GB минимум (рекомендуется 8GB)
- **Диск:** 20GB SSD минимум (рекомендуется 50GB для логов и backup)
- **Сеть:** Статический IP-адрес или доступ через обратный прокси

### Программное обеспечение:
- **ОС:** Ubuntu 22.04 LTS (или RHEL 8+, CentOS Stream 9+)
- **База данных:** PostgreSQL 14+ (или предоставить доступ к существующему кластеру)
- **Опционально:** Docker + Docker Compose (упрощает развёртывание)

---

## 🌐 Требования к сети и домену

### Вариант 1 (Рекомендуемый): Поддомен
**URL:** `https://phd.kaznmu.edu.kz` или `https://doctoral.kaznmu.edu.kz`

**Требования:**
1. Создать DNS A-запись: 
   - `phd.kaznmu.edu.kz` → IP сервера приложения
2. Выдать SSL сертификат (Let's Encrypt или корпоративный)
3. Открыть порты:
   - 80 (HTTP → redirect на HTTPS)
   - 443 (HTTPS)

### Вариант 2: Путь на основном домене
**URL:** `https://kaznmu.edu.kz/phd-portal/`

**Требования:**
1. Настроить обратный прокси (Nginx/Apache) на веб-сервере университета
2. Перенаправление `/phd-portal/*` на сервер приложения
3. Доступ к конфигурации веб-сервера для настройки правил

---

## 🔐 Требования безопасности

### Доступы:
1. **SSH доступ** для развёртывания и обновлений (ключи SSH, не пароли)
2. **Доступ к базе данных:**
   - Создать пользователя PostgreSQL с правами на создание БД
   - Или предоставить готовую БД с правами INSERT/UPDATE/DELETE/SELECT
3. **Backup доступ** (опционально):
   - Место для автоматических backup БД
   - Рекомендуется: ежедневный backup + недельное хранение

### Сетевая безопасность:
- Firewall правила: разрешить входящие на 80/443 (или только внутренний IP)
- Опционально: VPN доступ для администрирования
- Рекомендуется: Fail2ban для защиты от bruteforce

---

## 📦 Варианты развёртывания

### Вариант A: Docker Compose (Рекомендуется - проще обслуживать)

**Что нужно установить:**
```bash
# Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Docker Compose
sudo apt install docker-compose-plugin
```

**Развёртывание (1 команда):**
```bash
git clone https://github.com/AlmatJuvashev/phd-students-portal.git
cd phd-students-portal
cp .env.example .env
# Отредактировать .env с настройками университета
docker compose up -d
```

**Обновление приложения:**
```bash
git pull
docker compose up -d --build
```

### Вариант B: Ручная установка (Традиционный способ)

**Требуемое ПО:**
- Go 1.21+ (`sudo apt install golang-go`)
- Node.js 18+ (`curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -`)
- PostgreSQL 14+ (`sudo apt install postgresql-14`)
- Nginx (`sudo apt install nginx`)

**Скрипт автоматической установки предоставлен** в `deploy/scripts/install.sh`

---

## 🗄️ База данных

### Вариант 1: Создать новую БД на сервере приложения
```sql
CREATE DATABASE phd_portal;
CREATE USER phd_admin WITH ENCRYPTED PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE phd_portal TO phd_admin;
```

### Вариант 2: Использовать существующий кластер PostgreSQL университета
**Требуется:**
- Connection string в формате: 
  ```
  postgresql://user:password@db-server:5432/phd_portal?sslmode=require
  ```
- Права: CREATE TABLE, INSERT, UPDATE, DELETE, SELECT
- ~100MB начальный размер БД (будет расти с данными)

---

## 📊 Миграции базы данных

Приложение включает автоматические миграции:

```bash
# Применить все миграции
cd backend && make migrate-up

# Откатить последнюю миграцию (если нужно)
make migrate-down
```

Файлы миграций находятся в `backend/db/migrations/`

---

## 🔄 Backup и восстановление

### Автоматический backup (рекомендуется):

Создать cron задачу:
```bash
# Открыть crontab
crontab -e

# Добавить ежедневный backup в 2:00 ночи
0 2 * * * /usr/local/bin/pg_dump -h localhost -U phd_admin phd_portal | gzip > /backup/phd_portal_$(date +\%Y\%m\%d).sql.gz

# Удалять backup старше 7 дней
0 3 * * * find /backup -name "phd_portal_*.sql.gz" -mtime +7 -delete
```

### Восстановление из backup:
```bash
gunzip < /backup/phd_portal_20251013.sql.gz | psql -h localhost -U phd_admin phd_portal
```

---

## 📈 Мониторинг и логи

### Логи приложения:
- **Backend:** `/var/log/phd-portal/backend.log` (если ручная установка)
- **Docker:** `docker compose logs -f backend`

### Health check endpoint:
```bash
curl https://phd.kaznmu.edu.kz/api/health
# Ответ: {"status":"ok"}
```

### Рекомендуемый мониторинг:
- Проверка доступности каждые 5 минут
- Оповещение при недоступности >5 минут
- Мониторинг использования диска (логи и БД)

---

## 🔧 Обслуживание и обновления

### Обновление приложения:

**Docker версия:**
```bash
cd /opt/phd-portal
git pull origin main
docker compose down
docker compose up -d --build
```

**Ручная установка:**
```bash
cd /opt/phd-portal
git pull origin main
cd backend && make build
sudo systemctl restart phd-portal-backend
cd ../frontend && npm run build
# Nginx автоматически подхватит новые файлы
```

### Перезапуск сервисов:
```bash
# Docker
docker compose restart

# Systemd
sudo systemctl restart phd-portal-backend
sudo systemctl restart nginx
```

---

## 👥 Контакты для технической поддержки

**Разработчик:** Almat Juvashev  
**Email:** [your-email@example.com]  
**Telegram:** [@your-telegram] (опционально)  

**Репозиторий:** https://github.com/AlmatJuvashev/phd-students-portal  
**Документация:** См. `deploy/DEPLOYMENT_GUIDE.md` для детальных инструкций

---

## ✅ Чек-лист для IT отдела

- [ ] Выделен сервер с минимальными требованиями
- [ ] Создан поддомен `phd.kaznmu.edu.kz` (или путь на основном домене)
- [ ] Настроен SSL сертификат
- [ ] Создана база данных PostgreSQL
- [ ] Предоставлен SSH доступ разработчику
- [ ] Настроены firewall правила (80, 443)
- [ ] Настроен автоматический backup БД
- [ ] Настроен мониторинг доступности
- [ ] Предоставлены учётные данные для production (БД, SMTP и т.д.)

---

## 📞 Процесс развёртывания

1. **IT отдел** предоставляет доступ к серверу
2. **Разработчик** производит установку (2-4 часа)
3. **Совместное тестирование** (1-2 дня)
4. **Обучение администраторов** (опционально, 1-2 часа)
5. **Запуск в production**

**Ориентировочное время:** 3-5 рабочих дней от получения доступа до запуска

---

## 💰 Ориентировочные затраты

### Если использовать существующую инфраструктуру университета:
- **Дополнительные затраты:** Минимальные (только время IT специалиста)
- **Сервер:** Использование существующего
- **Домен:** Поддомен на kaznmu.edu.kz (бесплатно)
- **SSL:** Let's Encrypt (бесплатно)
- **БД:** Использование существующего кластера

### Если нужен отдельный сервер:
- **VPS/Dedicated:** от 10,000₸/месяц (зависит от провайдера)
- **Или облако (временно):** Railway/Heroku ~$10-20/месяц

---

## 📄 Приложения

1. Детальное руководство по развёртыванию: `deploy/DEPLOYMENT_GUIDE.md`
2. Docker Compose конфигурация: `docker-compose.yml`
3. Скрипты установки: `deploy/scripts/`
4. Конфигурация Nginx: `deploy/nginx.conf`

---

**Дата создания:** 13 октября 2025  
**Версия документа:** 1.0  
**Статус:** Готов к развёртыванию
