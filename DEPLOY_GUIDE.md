# 🚀 Deployment Guide

## Quick Deploy to Render.com (Recommended)

### Why Render?
- ✅ **100% Free tier** (750 hours/month)
- ✅ **Zero configuration** - works out of the box
- ✅ **Auto-deploy** from GitHub
- ✅ **HTTPS included**
- ✅ **WebSocket support**
- ✅ **No credit card required**

### Step-by-Step Instructions

#### 1. Prepare Your Code

Make sure all changes are committed:

```bash
git add .
git commit -m "Ready for Render deployment"
git push origin main
```

#### 2. Create Render Account

1. Go to https://render.com
2. Click "Get Started"
3. Sign up with GitHub (easiest option)
4. Authorize Render to access your repositories

#### 3. Create Web Service

1. From Render Dashboard, click **"New +"** button
2. Select **"Web Service"**
3. Connect your **GitHub repository** (goroom)
4. Render will automatically detect:
   - Language: **Go**
   - Build Command: `go build -o chat-service`
   - Start Command: `./chat-service`

5. Configure the service:
   - **Name:** `goroom-chat` (or any name you like)
   - **Region:** Choose closest to you
   - **Branch:** `main`
   - **Plan:** Select **Free**

6. Click **"Create Web Service"**

#### 4. Wait for Deploy

- Build takes 2-3 minutes ☕
- You'll see build logs in real-time
- When you see "Your service is live 🎉", it's ready!

#### 5. Access Your App

Your app will be available at:
```
https://goroom-chat.onrender.com
```

WebSocket endpoint:
```
wss://goroom-chat.onrender.com/ws
```

### Important Notes

**Free Tier Limitations:**
- ⚠️ Service **sleeps after 15 minutes** of inactivity
- ⏰ First request after sleep takes **~30 seconds** to wake up
- ✅ Perfect for demos and testing
- ✅ Upgrade to paid ($7/mo) to keep it always running

**Auto-Deploy:**
- 🔄 Every `git push` to `main` triggers new deployment
- 📧 You'll get email notifications about deploys

---

## Alternative: Railway.app

### Why Railway?
- ✅ **$5 free credit** every month
- ✅ **Never sleeps** (even on free tier)
- ✅ **Faster** than Render
- ⚠️ Requires credit card (won't charge on free tier)

### Deploy Steps

1. **Install Railway CLI:**
   ```bash
   npm install -g @railway/cli
   ```

2. **Login:**
   ```bash
   railway login
   ```

3. **Initialize:**
   ```bash
   cd real-time-service
   railway init
   ```

4. **Deploy:**
   ```bash
   railway up
   ```

5. **Get URL:**
   ```bash
   railway open
   ```

Done! Your app is live.

---

## Alternative: Fly.io

### Why Fly?
- ✅ Free tier with credit card
- ✅ Global CDN
- ✅ Never sleeps
- ⚠️ Slightly more complex

### Deploy Steps

1. **Install flyctl:**
   ```bash
   # Windows
   powershell -Command "iwr https://fly.io/install.ps1 -useb | iex"
   
   # Mac/Linux
   curl -L https://fly.io/install.sh | sh
   ```

2. **Login:**
   ```bash
   fly auth login
   ```

3. **Launch:**
   ```bash
   cd real-time-service
   fly launch
   ```

4. **Follow prompts:**
   - Choose app name
   - Select region
   - Don't create database
   - Deploy now: Yes

Done!

---

## Updating Your Deployed App

### Render.com
Just push to GitHub:
```bash
git add .
git commit -m "Update features"
git push origin main
```
Render auto-deploys!

### Railway
```bash
railway up
```

### Fly.io
```bash
fly deploy
```

---

## Troubleshooting

### Build Fails on Render

**Issue:** `go.mod` not found or build errors

**Solution:**
1. Make sure `go.mod` is in repository root
2. Run `go mod tidy` locally
3. Commit and push

### WebSocket Not Working

**Issue:** WebSocket connections failing

**Solution:**
- ✅ Use `wss://` (not `ws://`) for HTTPS sites
- ✅ Update JS client to use production URL
- ✅ Check Render logs for errors

### App Sleeping Too Often (Render)

**Solutions:**
1. **Free:** Use https://uptimerobot.com to ping every 5 minutes
2. **Paid:** Upgrade to Render paid plan ($7/mo)
3. **Alternative:** Switch to Railway (doesn't sleep)

---

## Environment Variables

### On Render

Add in Render Dashboard → Environment:

```
PORT=8080           # Auto-set by Render
GIN_MODE=release    # Production mode
```

### On Railway

```bash
railway variables set GIN_MODE=release
```

### On Fly.io

Edit `fly.toml`:
```toml
[env]
  GIN_MODE = "release"
```

---

## Custom Domain (Optional)

### Render
1. Go to Settings → Custom Domain
2. Add your domain
3. Update DNS records as shown
4. SSL is automatic!

### Railway
```bash
railway domain
```

### Fly.io
```bash
fly certs add yourdomain.com
```

---

## Monitoring & Logs

### View Logs

**Render:**
- Dashboard → Your Service → Logs tab
- Real-time log streaming

**Railway:**
```bash
railway logs
```

**Fly.io:**
```bash
fly logs
```

---

## Cost Comparison

| Platform | Free Tier | Sleeps? | Requires Card? | Paid Price |
|----------|-----------|---------|----------------|------------|
| **Render** | 750h/mo | Yes (15min) | No | $7/mo |
| **Railway** | $5 credit | No | Yes | Pay-as-go |
| **Fly.io** | Limited | No | Yes | $1.94/mo |
| **Heroku** | None | - | - | $7/mo |

---

## Production Checklist

Before deploying to production:

- [ ] Set `GIN_MODE=release` in environment
- [ ] Configure proper logging
- [ ] Add rate limiting
- [ ] Set up monitoring (UptimeRobot)
- [ ] Configure custom domain
- [ ] Enable HTTPS (auto on all platforms)
- [ ] Test WebSocket connections
- [ ] Add error tracking (Sentry)

---

## Support

**Issues with deployment?**
- Check Render docs: https://render.com/docs
- Railway docs: https://docs.railway.app
- Fly.io docs: https://fly.io/docs

**App-specific issues?**
- Check GitHub Issues
- Review application logs

---

**Happy Deploying! 🚀**
