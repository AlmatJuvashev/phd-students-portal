# Mailserver (Mailpit)

Local email testing.

## Run
```bash
docker compose up -d
```

- SMTP: `localhost:1025`
- Web UI: http://localhost:8025

Configure backend `.env` to point to this SMTP for password reset emails.
