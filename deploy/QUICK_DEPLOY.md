# 🚀 Quick Deploy Guide - KazNMU PhD Portal

## Demo Version (Railway + Vercel) - 15 minutes

### 1. Push to GitHub

```bash
git add .
git commit -m "prepare for railway deploy"
git push origin main
```

### 2. Deploy Backend to Railway

1. Go to [railway.app](https://railway.app) and login with GitHub
2. Click **"New Project"** → **"Deploy from GitHub repo"**
3. Select `phd-students-portal`
4. Railway will auto-detect and start deploying

**Add PostgreSQL:**

- Click **"New"** → **"Database"** → **"Add PostgreSQL"**
- Railway auto-sets `DATABASE_URL`

**Add Environment Variables:**

- Click your service → **"Variables"** tab
- Add:
  ```
  JWT_SECRET=my-super-secret-jwt-key-change-this-in-production
  GIN_MODE=release
  CORS_ORIGINS=https://your-app.vercel.app,http://localhost:5173
  ```

**Run Migrations:**

- Click **"Settings"** → **"Deploy"**
- Add custom build command: `cd backend && make migrate-up && make run`
- Or use Railway CLI after deploy

**Copy your backend URL** (something like `https://phd-portal-production.up.railway.app`)

---

### 3. Deploy Frontend to Vercel

1. Go to [vercel.com](https://vercel.com) and login with GitHub
2. Click **"Add New"** → **"Project"**
3. Import `phd-students-portal` repository
4. Configure:
   - **Framework Preset:** Vite
   - **Root Directory:** `frontend`
   - **Build Command:** `npm run build`
   - **Output Directory:** `dist`
5. **Environment Variables:**
   ```
   VITE_API_BASE_URL=https://your-backend.railway.app/api
   ```
6. Click **"Deploy"**

**Copy your frontend URL** (something like `https://phd-portal.vercel.app`)

---

### 4. Update Backend CORS

Go back to Railway → Your service → Variables → Update:

```
CORS_ORIGINS=https://your-frontend.vercel.app,http://localhost:5173
```

Redeploy the backend (Railway will auto-redeploy on changes)

---

### 5. ✅ Test Your Demo!

Open your Vercel URL and test the application.

**Demo credentials** (if you seed the database):

- Email: `admin@kaznmu.edu.kz`
- Password: `admin123`

---

## Production Version (University Server)

See [`deploy/UNIVERSITY_IT_REQUIREMENTS.md`](deploy/UNIVERSITY_IT_REQUIREMENTS.md) for full production deployment guide.

**Key differences:**

- Backend + Frontend on university server
- Custom domain: `phd.kaznmu.edu.kz`
- University's PostgreSQL cluster
- SSL certificate from university CA
- Backup and monitoring setup

---

## Costs Comparison

| Option                    | Frontend      | Backend        | Database | Total/month |
| ------------------------- | ------------- | -------------- | -------- | ----------- |
| **Demo (Railway+Vercel)** | Free          | $5 free credit | Included | $0-5        |
| **University Server**     | Free          | Free           | Free     | $0          |
| **Cloud (DigitalOcean)**  | Free (Vercel) | $6             | Included | $6          |

---

## Next Steps After Demo

1. **Show demo to university management**
2. **Get approval from IT department**
3. **Request server access** (see requirements doc)
4. **Migration to production** (1-2 weeks)
5. **Launch!** 🎉

---

## Support

- **Issues:** https://github.com/AlmatJuvashev/phd-students-portal/issues
- **Email:** [your-email]
- **Documentation:** See `deploy/` folder

---

**Created:** October 13, 2025  
**Version:** 1.0 Demo
