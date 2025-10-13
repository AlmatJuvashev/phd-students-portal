# Railway Deployment Guide for KazNMU PhD Portal

## Quick Deploy to Railway (Demo Version)

### Prerequisites
- GitHub account with your repository
- Railway.app account (sign up at https://railway.app)

### Step 1: Push to GitHub

```bash
# If not already done
git remote add origin https://github.com/AlmatJuvashev/phd-students-portal.git
git branch -M main
git push -u origin main
```

### Step 2: Deploy Backend to Railway

1. **Go to Railway.app** → Login with GitHub
2. **Click "New Project"** → "Deploy from GitHub repo"
3. **Select your repository**: `phd-students-portal`
4. **Railway will auto-detect** Go backend

5. **Add PostgreSQL Database:**
   - Click "New" → "Database" → "Add PostgreSQL"
   - Railway will automatically set `DATABASE_URL` environment variable

6. **Set Environment Variables:**
   Go to your service → Variables → Add these:
   ```
   JWT_SECRET=your-random-secret-key-change-this
   GIN_MODE=release
   PORT=8080
   CORS_ORIGINS=https://your-frontend-url.vercel.app
   ```

7. **Get your backend URL:**
   - Copy the public URL (e.g., `https://phd-portal-production.up.railway.app`)

### Step 3: Deploy Frontend to Vercel (Free & Fast)

1. **Go to vercel.com** → Login with GitHub
2. **Import your repository**
3. **Configure:**
   - Framework Preset: Vite
   - Root Directory: `frontend`
   - Build Command: `npm run build`
   - Output Directory: `dist`
   
4. **Environment Variables:**
   ```
   VITE_API_BASE_URL=https://your-backend-url.railway.app/api
   ```

5. **Deploy** → Copy frontend URL

6. **Update Backend CORS:**
   - Go back to Railway → Backend service → Variables
   - Update `CORS_ORIGINS` with your Vercel URL

### Step 4: Run Database Migrations

Railway provides a terminal for your service:

1. Go to your backend service → "Terminal" tab (or use Railway CLI)
2. Run migrations:
   ```bash
   cd backend
   make migrate-up
   ```

Or manually:
```bash
psql $DATABASE_URL -f db/migrations/0001_init.up.sql
psql $DATABASE_URL -f db/migrations/0002_comments.up.sql
```

### Step 5: Test Your Demo

Visit your Vercel URL and test the application!

---

## Alternative: All-in-One Railway Deploy

If you want both frontend and backend on Railway:

### Create `nixpacks.toml` in root:

```toml
[phases.setup]
nixPkgs = ["nodejs-18_x", "go_1_21"]

[phases.install]
cmds = [
  "cd frontend && npm ci",
  "cd backend && go mod download"
]

[phases.build]
cmds = [
  "cd frontend && npm run build",
  "cd backend && go build -o bin/server cmd/server/main.go"
]

[start]
cmd = "cd backend && ./bin/server"
```

---

## Costs

- **Railway:** $5 free credit/month → ~$0.01/hour after
- **Vercel:** Unlimited free for frontend
- **Total Demo Cost:** ~$0-5/month

---

## After Demo Approval → Production Migration

See [`UNIVERSITY_IT_REQUIREMENTS.md`](./UNIVERSITY_IT_REQUIREMENTS.md) for moving to university servers.
