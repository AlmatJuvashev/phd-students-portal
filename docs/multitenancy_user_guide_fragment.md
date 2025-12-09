
---

## User Guide & Login

### How to Access a Tenant

Users must access the portal through the specific entry point for their organization. A user belonging to "Demo University" cannot log in to "Kazakh National Medical University".

#### 1. Via Frontend (Browser)
The application uses **Subdomains** to identify the tenant.
- **Main Portal**: `http://localhost:5173` (Defaults to `kaznmu`)
- **Demo University**: `http://demo.localhost:5173`
- **Other Tenants**: `http://<slug>.localhost:5173`

*Note: You may need to configure your `/etc/hosts` file to support subdomains locally.*

#### 2. Via API (Curl / Postman)
You must explicitly tell the API which tenant you are accessing using the `X-Tenant-Slug` header.

```bash
# Login to Demo University
curl -X POST http://localhost:8280/api/auth/login \
  -H "X-Tenant-Slug: demo" \
  -d '{"username": "demo.admin", "password": "..."}'
```

Without this header, the API defaults to the main tenant (`kaznmu`), and login will fail for users who don't belong to it.

### The "Demo" Tenant

A fully populated Demo tenant is available for testing and showcasing features.

- **Tenant Slug**: `demo`
- **Tenant Name**: Demo University
- **URL**: `http://demo.localhost:5173`

#### Demo Credentials
| Role | Username | Password |
|------|----------|----------|
| **Admin** | `demo.admin` | `demopassword123!` |
| **Advisor** | `dr.johnson` | `demopassword123!` |
| **Student** | `demo.student1` | `demopassword123!` |

*Note: `demo.student1` through `demo.student24` are available.*

### Troubleshooting Login

**Error: "У вас нет доступа к этому порталу" (You do not have access to this portal)**

- **Cause**: Your username/password is correct, but you are trying to log in to the wrong tenant.
- **Example**: You are logging in as `demo.admin` (who belongs to `demo`) on `http://localhost:5173` (which is `kaznmu`).
- **Fix**: Switch to the correct subdomain (`demo.localhost:5173`) or add the `X-Tenant-Slug: demo` header.
