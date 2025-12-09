# SMTP Configuration for Production

## Overview

The PhD Portal uses SMTP for sending email notifications including password resets, submission notifications, and deadline reminders.

---

## Gmail SMTP

### Setup

1. Enable **2-Step Verification** at [Google Account Security](https://myaccount.google.com/security)
2. Generate **App Password** at [App Passwords](https://myaccount.google.com/apppasswords)
3. Configure `.env`:

```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-16-char-app-password
SMTP_FROM="PhD Portal <your-email@gmail.com>"
```

### Limits
- 500 emails/day (free tier)
- Not recommended for high-volume production

---

## SendGrid (Recommended for Production)

### Setup

1. Create account at [SendGrid](https://sendgrid.com)
2. Generate API Key under **Settings > API Keys**
3. Configure `.env`:

```env
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USER=apikey
SMTP_PASS=SG.xxxxx  # Your SendGrid API Key
SMTP_FROM="PhD Portal <no-reply@yourdomain.com>"
```

### Benefits
- 100 emails/day free tier
- Delivery analytics
- Email templates

---

## AWS SES (Simple Email Service)

### Setup

1. Verify domain/email in AWS SES Console
2. Create SMTP credentials in **SES > SMTP Settings**
3. Configure `.env`:

```env
SMTP_HOST=email-smtp.us-east-1.amazonaws.com  # Use your region
SMTP_PORT=587
SMTP_USER=your-smtp-username
SMTP_PASS=your-smtp-password
SMTP_FROM="PhD Portal <no-reply@yourdomain.com>"
```

### Benefits
- Pay-per-use pricing (~$0.10/1000 emails)
- High deliverability
- Integrates with AWS infrastructure

---

## Mailgun

### Setup

1. Create account at [Mailgun](https://mailgun.com)
2. Add and verify your domain
3. Get SMTP credentials from dashboard
4. Configure `.env`:

```env
SMTP_HOST=smtp.mailgun.org
SMTP_PORT=587
SMTP_USER=postmaster@yourdomain.com
SMTP_PASS=your-mailgun-password
SMTP_FROM="PhD Portal <no-reply@yourdomain.com>"
```

---

## Development (Mailpit)

For local development, use Mailpit to capture emails without sending:

```env
SMTP_HOST=localhost
SMTP_PORT=1025
SMTP_USER=
SMTP_PASS=
SMTP_FROM="PhD Portal <no-reply@phd.local>"
```

**View emails:** http://localhost:8025

---

## Testing

After configuration, test by:
1. Triggering a password reset
2. Creating a submission that notifies advisors
3. Check inbox or Mailpit UI
