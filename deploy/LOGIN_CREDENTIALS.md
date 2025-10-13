# ✅ РЕШЕНИЕ: Login с admin credentials

## Проблема
Вы пытались войти используя **email** (`juvashev@gmail.com`), но система ожидает **username**.

## ✅ Решение

При создании admin-пользователя из `ADMIN_EMAIL` система автоматически генерирует **username** из части до `@`:

**Ваши credentials:**
- Email: `<ваш-email>`
- Password: `<ваш-пароль-из-ADMIN_PASSWORD>`
- **Username (для логина)**: `<часть-email-до-@>` ← используйте это!

### Войдите так:

```
Username: <часть-вашего-email-до-@>
Password: <ваш-пароль-из-ADMIN_PASSWORD>
```

**Не используйте** full email для логина — только username!

---

## 📋 Все тестовые учётные записи

### 1. Superadmin (созданный из ADMIN_EMAIL)
```
Username: <из-вашего-ADMIN_EMAIL>
Password: <из-ADMIN_PASSWORD>
Role: superadmin
```

### 2. Дефолтный admin (из seed-скрипта, если есть)
```
Username: admin
Password: admin123
Role: admin
```

### 3. Тестовый докторант (если создан в seed)
```
Username: докторант
Password: докторант123
Role: student
```

---

## Как проверить username для любого email

Если вы забыли username, но знаете email, посмотрите в Railway logs при старте — там должно быть сообщение о создании admin:

```
Superadmin user created: username=<username> email=<email>
```

Или выполните SQL запрос через Railway Postgres:
```sql
SELECT username, email, role FROM users WHERE email = 'your-email@example.com';
```

---

## Почему так сделано?

- **Username** уникален и короток → удобно для логина
- **Email** используется для восстановления пароля и уведомлений
- Система автоматически создаёт username из email при bootstrap

---

## После успешного логина

1. ✅ Получите JWT token
2. ✅ Редирект на главную (`/`)
3. ✅ Доступ к admin-панели (`/admin/users`)
4. ✅ Можете создать других пользователей через UI

---

**Попробуйте сейчас:** используйте username (часть email до `@`) вместо полного email для логина! 🚀
