# 🚀 Railway Deployment Checklist

## ✅ Pre-Deploy Checklist

- [x] All code committed to GitHub
- [x] Deployment documentation created
- [x] Environment variables template ready
- [x] Railway configuration files added
- [ ] **Push to GitHub** ⬅️ DO THIS NOW

## 📝 Deployment Steps

### Step 1: Push to GitHub

```bash
git push origin main
```

### Step 2: Deploy Backend to Railway

1. **Go to https://railway.app**
   - Sign up / Login with GitHub
2. **Create New Project**
   - Click "New Project"
   - Select "Deploy from GitHub repo"
   - Choose `phd-students-portal`
3. **Add PostgreSQL Database**
   - Click "New" → "Database" → "Add PostgreSQL"
   - Railway automatically sets `DATABASE_URL`
4. **Set Environment Variables**
   Click your backend service → "Variables" tab → Add:

   ```
   JWT_SECRET=super-secret-jwt-key-change-this-now
   GIN_MODE=release
   PORT=8080
   CORS_ORIGINS=http://localhost:5173
   ADMIN_EMAIL=admin@kaznmu.edu.kz
   ADMIN_PASSWORD=Admin2024!
   ```

5. **Wait for Deploy** (2-5 minutes)

   - Railway will build and deploy automatically
   - Copy your backend URL (e.g., `https://xxx.railway.app`)

6. **Migrations will run automatically!**

   - After you push updated files (see Step 1)
   - Railway will run migrations automatically via `Procfile` `release` command
   - Check logs: Deployments → View Logs → Look for "✅ All migrations applied"

   **If you need to run manually:** See [`deploy/MIGRATIONS_GUIDE.md`](deploy/MIGRATIONS_GUIDE.md)

phd-students-portal-production.up.railway.app

### Step 3: Deploy Frontend to Vercel

1. **Go to https://vercel.com**
   - Sign up / Login with GitHub
2. **Import Project**
   - Click "Add New..." → "Project"
   - Select `phd-students-portal`
3. **Configure Build**
   - Framework Preset: **Vite**
   - Root Directory: **`frontend`**
   - Build Command: `npm run build`
   - Output Directory: `dist`
4. **Add Environment Variable**
   ```
   VITE_API_BASE_URL=https://your-backend.railway.app/api
   ```
   (Replace with your actual Railway backend URL)
5. **Deploy**
   - Click "Deploy"
   - Wait 2-3 minutes
   - Copy your frontend URL (e.g., `https://xxx.vercel.app`)

### Step 4: Update CORS

1. **Go back to Railway**
   - Your backend service → "Variables"
   - Update `CORS_ORIGINS`:
     ```
     CORS_ORIGINS=https://your-frontend.vercel.app,http://localhost:5173
     ```
   - Service will auto-redeploy

### Step 5: Test!

1. **Open your Vercel URL**
2. **Login with:**
   - Email: `admin@kaznmu.edu.kz`
   - Password: `Admin2024!` (or what you set)
3. **Test features:**
   - Create a user
   - Navigate journey map
   - Upload documents
   - Add comments

---

## 🔧 Troubleshooting

### Backend won't start?

**Check Railway logs:**

- Go to service → "Deployments" → Click latest → "View Logs"
- Look for errors

**Common issues:**

- `DATABASE_URL not set` → PostgreSQL not added
- `Port already in use` → Change PORT variable
- `Migration failed` → Run migrations manually

### Frontend can't connect to backend?

**Check:**

1. `VITE_API_BASE_URL` is correct in Vercel
2. CORS includes your Vercel URL
3. Backend is actually running (check Railway logs)

### Database migrations not applied?

**Option A: Railway CLI**

```bash
# Install Railway CLI
npm i -g @railway/cli

# Login
railway login

# Link project
railway link

# Run migrations
railway run cd backend && make migrate-up
```

**Option B: Manual psql**

```bash
# Get DATABASE_URL from Railway variables
psql "postgresql://user:pass@host:port/db" -f backend/db/migrations/0001_init.up.sql
psql "postgresql://user:pass@host:port/db" -f backend/db/migrations/0002_comments.up.sql
```

---

## 📊 Demo URLs

After deployment, you'll have:

- **Frontend:** `https://kaznmu-phd-portal.vercel.app`
- **Backend:** `https://kaznmu-phd-portal-production.up.railway.app`
- **API Health:** `https://backend-url/api/health`

---

## 💰 Costs

- **Railway:** $5 free/month (enough for demo)
- **Vercel:** Free forever for frontend
- **Total:** $0-5/month

---

## 📞 Next Steps After Demo

1. **Show to university management**

   - Share Vercel URL
   - Demo all features
   - Get feedback

2. **Present to IT department**

   - Give them `deploy/UNIVERSITY_IT_REQUIREMENTS.md`
   - Request server access
   - Schedule migration meeting

3. **Migration to production**
   - Export/backup Railway database
   - Deploy to university server
   - Import data
   - Setup `phd.kaznmu.edu.kz` domain

---

## 🆘 Need Help?

**Railway Documentation:** https://docs.railway.app  
**Vercel Documentation:** https://vercel.com/docs  
**Your Repository:** https://github.com/AlmatJuvashev/phd-students-portal

---

**Created:** October 13, 2025  
**Status:** Ready to Deploy! 🚀
