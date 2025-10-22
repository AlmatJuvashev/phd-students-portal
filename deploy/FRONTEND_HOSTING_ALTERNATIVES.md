# Frontend Hosting Alternatives to Vercel

## Quick Comparison

| Platform             | Free Tier    | Setup Time | Speed    | GitHub Integration |
| -------------------- | ------------ | ---------- | -------- | ------------------ |
| **Netlify**          | ✅ Generous  | 5 min      | ⚡⚡⚡   | ✅ Yes             |
| **Cloudflare Pages** | ✅ Unlimited | 5 min      | ⚡⚡⚡⚡ | ✅ Yes             |
| **Railway**          | ✅ Limited   | 3 min      | ⚡⚡     | ✅ Yes             |
| **Render**           | ✅ Good      | 5 min      | ⚡⚡⚡   | ✅ Yes             |

## Option 1: Netlify (Recommended)

### Why Netlify?

- Very similar to Vercel
- Simple setup
- Excellent free tier
- Fast CDN
- Great for Kazakhstan region

### Setup via Web UI

1. **Sign up**: https://app.netlify.com/signup
2. **Import project** from GitHub
3. **Configure build settings**:

   - Base directory: `frontend`
   - Build command: `npm run build`
   - Publish directory: `dist`

4. **Add Environment Variables**:

   ```
   VITE_API_URL=https://phd-student-portal-starter-v8-production.up.railway.app/api
   ```

5. **Deploy**: Click "Deploy site"

### Setup via CLI

```bash
# Install Netlify CLI
npm install -g netlify-cli

# Login to Netlify
netlify login

# Deploy from frontend directory
cd frontend
netlify init
netlify deploy --prod
```

The `netlify.toml` file in the project root already configures everything needed.

---

## Option 2: Cloudflare Pages

### Why Cloudflare Pages?

- Fastest CDN globally
- Unlimited bandwidth (free)
- Great performance in Asia/Kazakhstan
- Simple GitHub integration

### Setup

1. **Sign up**: https://dash.cloudflare.com/sign-up/pages
2. **Connect GitHub** repository
3. **Configure build**:

   - Framework preset: `Vite`
   - Build command: `npm run build`
   - Build output directory: `dist`
   - Root directory: `frontend`

4. **Environment Variables**:

   ```
   VITE_API_URL=https://phd-student-portal-starter-v8-production.up.railway.app/api
   NODE_VERSION=18
   ```

5. **Deploy**: Automatic on push to main

---

## Option 3: Railway (Frontend + Backend Together)

### Why Railway for Frontend?

- Already hosting your backend
- Single platform for everything
- Simple configuration
- Good for development/staging

### Setup

1. **Go to Railway Dashboard**: https://railway.app/dashboard
2. **Create New Project** or add service to existing project
3. **Select "Deploy from GitHub repo"**
4. **Configure**:

   - Root Directory: `frontend`
   - Build Command: `npm run build && npx serve -s dist -l $PORT`
   - Or use the generated service

5. **Environment Variables**:

   ```
   VITE_API_URL=https://phd-student-portal-starter-v8-production.up.railway.app/api
   PORT=3000
   ```

6. **Add `package.json` script** in frontend:

   ```json
   "scripts": {
     "serve": "serve -s dist -l $PORT"
   }
   ```

7. **Install serve**:
   ```bash
   cd frontend
   npm install serve
   ```

Railway will automatically detect the project and deploy.

---

## Option 4: Render

### Why Render?

- Simple static site hosting
- Automatic SSL
- Good free tier
- Fast builds

### Setup

1. **Sign up**: https://dashboard.render.com/register
2. **New Static Site** from GitHub
3. **Configure**:

   - Build Command: `cd frontend && npm install && npm run build`
   - Publish Directory: `frontend/dist`

4. **Environment Variables**:

   ```
   VITE_API_URL=https://phd-student-portal-starter-v8-production.up.railway.app/api
   NODE_VERSION=18
   ```

5. **Deploy**: Automatic

---

## After Deployment: Update Backend CORS

Once you deploy to any platform, update your Railway backend environment variable:

```bash
# In Railway backend service, update:
FRONTEND_BASE=https://your-new-frontend-url.netlify.app
# or
FRONTEND_BASE=https://your-site.pages.dev  # for Cloudflare
# or
FRONTEND_BASE=https://your-frontend.up.railway.app  # for Railway
```

Then redeploy the backend service.

---

## Testing Your Deployment

After deployment, test:

1. ✅ Frontend loads correctly
2. ✅ Can login with username: `juvashev`
3. ✅ No CORS errors in console
4. ✅ Can navigate between modules
5. ✅ Forms submit successfully

---

## Troubleshooting

### CORS Errors

- Check `FRONTEND_BASE` in Railway backend matches your new frontend URL
- Ensure `APP_ENV=production` is set in Railway backend
- Redeploy backend after changing variables

### Build Fails

- Check Node version is 18 or higher
- Verify `VITE_API_URL` is set correctly
- Check build logs for specific errors

### 404 on Routes

- Ensure redirects are configured (netlify.toml handles this)
- For other platforms, configure SPA fallback to `/index.html`

---

## My Recommendation

**For your case, I recommend Netlify** because:

1. ✅ Most similar to Vercel (easy migration)
2. ✅ Excellent free tier
3. ✅ Fast deployment
4. ✅ Good performance in Kazakhstan
5. ✅ `netlify.toml` already configured in your project

Simply push your code to GitHub and deploy via Netlify web UI - takes 5 minutes!
