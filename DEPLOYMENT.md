# MysticFunds Railway Deployment Guide

This guide walks you through deploying MysticFunds to Railway with PostgreSQL database.

## Prerequisites

1. **Railway Account**: Sign up at [railway.app](https://railway.app)
2. **Railway CLI**: Install with `npm install -g @railway/cli`
3. **Git**: Ensure your code is committed to git

## Step 1: Prepare for Deployment

Make sure all code is committed:
```bash
git add .
git commit -m "Prepare for Railway deployment"
```

## Step 2: Deploy to Railway

### Option A: Using Railway CLI (Recommended)

1. **Login to Railway**:
   ```bash
   railway login
   ```

2. **Initialize Railway project**:
   ```bash
   railway init
   # Choose: "Empty Project"
   # Project name: mysticfunds-backend
   ```

3. **Add PostgreSQL database**:
   ```bash
   railway add postgresql
   ```

4. **Set environment variables**:
   ```bash
   # JWT Secret (generate a secure random string)
   railway variables set JWT_SECRET="your-super-secure-jwt-secret-here"
   
   # Optional: Set custom domain
   railway variables set RAILWAY_STATIC_URL="https://your-domain.com"
   ```

5. **Deploy the application**:
   ```bash
   railway up
   ```

### Option B: Using Railway Dashboard

1. Go to [railway.app](https://railway.app)
2. Click "Start a New Project"
3. Choose "Deploy from GitHub repo"
4. Connect your GitHub account and select your repository
5. Add PostgreSQL service:
   - Click "Add Service" → "Database" → "PostgreSQL"
6. Set environment variables in the dashboard:
   - `JWT_SECRET`: Generate a secure random string
7. Deploy!

## Step 3: Verify Deployment

1. **Check service status**:
   ```bash
   railway status
   ```

2. **View logs**:
   ```bash
   railway logs
   ```

3. **Get your domain**:
   ```bash
   railway domain
   ```

Your backend will be available at: `https://your-app-name.railway.app`

## Step 4: Configure Frontend for Railway Backend

Update your Vercel frontend to use the Railway backend:

1. **Update vercel.json**:
   ```json
   {
     "rewrites": [
       {
         "source": "/api/(.*)",
         "destination": "https://your-app-name.railway.app/api/$1"
       }
     ]
   }
   ```

2. **Deploy frontend to Vercel**:
   ```bash
   vercel --prod
   ```

## Environment Variables on Railway

Railway automatically provides these PostgreSQL environment variables:
- `PGHOST` - Database host
- `PGPORT` - Database port
- `PGUSER` - Database user
- `PGPASSWORD` - Database password
- `PGDATABASE` - Default database name
- `DATABASE_URL` - Complete connection string

You need to set:
- `JWT_SECRET` - Secret key for JWT tokens (generate a secure random string)
- `PORT` - Application port (Railway sets this automatically)

## Database Migrations

The deployment script automatically:
1. Waits for PostgreSQL to be ready
2. Creates the required databases (`auth`, `wizard`, `mana`)
3. Runs all database migrations
4. Starts all microservices

## Monitoring and Maintenance

### View Logs
```bash
railway logs --follow
```

### Connect to Database
```bash
railway connect postgresql
```

### Restart Services
```bash
railway up --detach
```

### Scale Services
Railway automatically scales based on traffic. You can configure scaling in the dashboard.

## Custom Domain (Optional)

1. In Railway dashboard, go to your project
2. Click on "Settings" → "Domains"
3. Add your custom domain
4. Update DNS records as instructed
5. Railway automatically handles SSL certificates

## Troubleshooting

### Common Issues

1. **Services not starting**:
   - Check logs: `railway logs`
   - Verify environment variables are set
   - Ensure database is ready

2. **Database connection errors**:
   - Verify PostgreSQL service is running
   - Check network configuration
   - Ensure migrations completed successfully

3. **Frontend not connecting to backend**:
   - Update API URLs in frontend
   - Check CORS configuration
   - Verify Railway domain is accessible

### Debug Commands

```bash
# Check environment variables
railway variables

# View service status
railway status

# Connect to database directly
railway connect postgresql

# View detailed logs
railway logs --tail 100
```

## Cost Optimization

Railway offers:
- **Free tier**: $5/month credit, perfect for development
- **Usage-based pricing**: Pay only for what you use
- **Sleep mode**: Inactive services automatically sleep to save costs

For production:
- Monitor usage in Railway dashboard
- Set up billing alerts
- Consider upgrading to Pro plan for better performance

## Security Considerations

1. **Environment Variables**: Never commit secrets to git
2. **JWT Secret**: Use a long, random string
3. **Database**: Railway handles security and backups
4. **HTTPS**: Automatically enabled on Railway domains
5. **CORS**: Configure properly for your frontend domains

## Next Steps

After successful deployment:

1. **Set up monitoring**: Use Railway's built-in metrics
2. **Configure custom domain**: For production use
3. **Set up CI/CD**: Automatic deployments on git push
4. **Database backups**: Railway handles automatic backups
5. **Performance optimization**: Monitor and optimize as needed

## Support

- **Railway Documentation**: [docs.railway.app](https://docs.railway.app)
- **Railway Discord**: Get help from the community
- **Railway Support**: Available for Pro plan users

Your MysticFunds backend should now be running on Railway with full database support!