# 🚀 MysticFunds Complete Deployment Guide

## Quick Deploy (Auto-Deployment on `dev` branch)

### 1. Create a `dev` branch
```bash
git checkout -b dev
git push -u origin dev
```

### 2. Set up secrets
```bash
# Generate secrets
./scripts/generate-secrets.sh

# The script will output commands like:
railway variables set JWT_SECRET="your-generated-secret"
```

### 3. Set up GitHub Secrets
Go to GitHub repo → Settings → Secrets → Actions and add:

- `RAILWAY_TOKEN` - Get from Railway dashboard
- `VERCEL_TOKEN` - Get from Vercel dashboard  
- `VERCEL_ORG_ID` - Get from Vercel
- `VERCEL_PROJECT_ID` - Get from Vercel

### 4. Deploy automatically
```bash
# Any push to dev branch triggers auto-deployment
git add .
git commit -m "Deploy to dev"
git push origin dev
```

## Manual Deployment

### Railway Backend Deployment

1. **Install CLI**:
   ```bash
   npm install -g @railway/cli
   ```

2. **Login and setup**:
   ```bash
   railway login
   railway init  # Choose "Empty Project"
   railway add postgresql
   ```

3. **Set secrets**:
   ```bash
   # Generate a secure JWT secret
   railway variables set JWT_SECRET="$(openssl rand -hex 32)"
   
   # Optional: Set custom environment
   railway variables set NODE_ENV="production"
   railway variables set LOG_LEVEL="info"
   ```

4. **Deploy**:
   ```bash
   railway up
   ```

5. **Get your backend URL**:
   ```bash
   railway domain
   # Your backend: https://yourapp.railway.app
   ```

### Vercel Frontend Deployment

1. **Install CLI**:
   ```bash
   npm install -g vercel
   ```

2. **Update backend URL in vercel.json**:
   ```json
   {
     "rewrites": [
       {
         "source": "/api/(.*)",
         "destination": "https://yourapp.railway.app/api/$1"
       }
     ]
   }
   ```

3. **Deploy**:
   ```bash
   vercel --prod
   ```

## Environment Variables & Secrets

### Required Secrets

| Secret | Where to set | Purpose |
|--------|--------------|---------|
| `JWT_SECRET` | Railway | JWT token signing |
| `RAILWAY_TOKEN` | GitHub | Auto-deployment |
| `VERCEL_TOKEN` | GitHub | Frontend deployment |

### Optional Secrets

| Secret | Purpose | Default |
|--------|---------|---------|
| `LOG_LEVEL` | Logging verbosity | `info` |
| `CORS_ORIGINS` | Allowed CORS origins | `*` |
| `RATE_LIMIT` | API rate limiting | `60/min` |

### How to get tokens:

1. **Railway Token**:
   - Go to Railway dashboard
   - Account Settings → Tokens → Create Token

2. **Vercel Tokens**:
   - Go to Vercel dashboard  
   - Settings → Tokens → Create Token
   - For Org ID: Settings → General → Organization ID
   - For Project ID: Project Settings → General → Project ID

## Branch-Based Deployment Strategy

### Recommended Git Flow:

```
main ──────────────────────── (production)
 │
 └── dev ──────────────────── (auto-deploys to Railway/Vercel)
      │
      ├── feature/auth ────── (development)
      ├── feature/jobs ────── (development)  
      └── feature/marketplace  (development)
```

### Auto-deployment triggers:

- **Push to `dev`**: Deploys to Railway + Vercel
- **Push to `main`**: Manual deployment (for production)
- **PR to `dev`**: Runs tests only

## Monitoring & Maintenance

### Check deployment status:
```bash
# Railway
railway logs --follow

# Vercel  
vercel logs

# GitHub Actions
# Go to GitHub → Actions tab
```

### Database management:
```bash
# Connect to Railway database
railway connect postgresql

# View database status
railway status
```

### Scaling:
- Railway: Auto-scales based on traffic
- Vercel: Auto-scales based on traffic
- No manual intervention needed

## Security Best Practices

### ✅ What we handle:
- Environment variables not committed to git
- Secure JWT secret generation
- HTTPS everywhere (automatic)
- Database security (Railway managed)

### ⚠️ Additional security for production:
- Enable rate limiting
- Set specific CORS origins
- Add API authentication for admin endpoints
- Monitor unusual traffic patterns

## Troubleshooting

### Common issues:

1. **"JWT secret not set"**:
   ```bash
   railway variables set JWT_SECRET="$(openssl rand -hex 32)"
   ```

2. **Database connection failed**:
   ```bash
   railway status  # Check if PostgreSQL is running
   ```

3. **Frontend can't reach backend**:
   - Check vercel.json has correct Railway URL
   - Verify CORS settings
   - Check Railway domain is accessible

4. **GitHub Actions failing**:
   - Verify all secrets are set in GitHub
   - Check Railway/Vercel tokens are valid
   - Review Actions logs for specific errors

### Debug commands:
```bash
# Check Railway variables
railway variables

# Test backend health
curl https://yourapp.railway.app/health

# Test frontend
curl https://yourapp.vercel.app

# View detailed logs
railway logs --tail 100
```

## Cost Estimation

### Free Tier Limits:
- **Railway**: $5/month credit (covers database + backend)
- **Vercel**: 100GB bandwidth, 1000 builds/month
- **GitHub Actions**: 2000 minutes/month

### Expected costs for MysticFunds:
- **Development**: $0 (within free tiers)
- **Light production use**: $5-15/month
- **Heavy production use**: $20-50/month

## Next Steps After Deployment

1. **Set up monitoring**: Railway/Vercel dashboards
2. **Configure custom domain**: Point to your Vercel frontend
3. **Set up error tracking**: Add error reporting service
4. **Database backups**: Railway handles automatically
5. **Performance optimization**: Monitor and tune as needed

## Support Resources

- **Railway**: [docs.railway.app](https://docs.railway.app)
- **Vercel**: [vercel.com/docs](https://vercel.com/docs)
- **GitHub Actions**: [docs.github.com/actions](https://docs.github.com/actions)

