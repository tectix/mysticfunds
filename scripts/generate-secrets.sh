#!/bin/bash

# MysticFunds Secret Generation Script
echo "Generating secrets for MysticFunds deployment..."
echo ""

# Generate JWT Secret (64 characters)
JWT_SECRET=$(openssl rand -hex 32)
echo "Generated JWT Secret:"
echo "JWT_SECRET=${JWT_SECRET}"
echo ""

# Generate additional secrets if needed
API_KEY=$(openssl rand -hex 16)
echo "Generated API Key (for future use):"
echo "API_KEY=${API_KEY}"
echo ""

# Database encryption key (for future use)
DB_ENCRYPTION_KEY=$(openssl rand -hex 24)
echo "Generated Database Encryption Key (for future use):"
echo "DB_ENCRYPTION_KEY=${DB_ENCRYPTION_KEY}"
echo ""

echo "Railway Deployment Commands:"
echo "railway variables set JWT_SECRET=\"${JWT_SECRET}\""
echo "railway variables set API_KEY=\"${API_KEY}\""
echo "railway variables set DB_ENCRYPTION_KEY=\"${DB_ENCRYPTION_KEY}\""
echo ""

echo "Add these to your GitHub Secrets:"
echo "1. Go to GitHub repo → Settings → Secrets → Actions"
echo "2. Add the following secrets:"
echo "   - JWT_SECRET: ${JWT_SECRET}"
echo "   - RAILWAY_TOKEN: (get from Railway dashboard)"
echo "   - VERCEL_TOKEN: (get from Vercel dashboard)"
echo "   - VERCEL_ORG_ID: (get from Vercel)"
echo "   - VERCEL_PROJECT_ID: (get from Vercel)"
echo ""

echo "Save these secrets securely - they won't be shown again!"
echo ""

# Optionally save to a secure file (not committed to git)
read -p "Save secrets to .env.local file? (y/n): " save_secrets
if [[ $save_secrets == "y" || $save_secrets == "Y" ]]; then
    cat > .env.local << EOF
# Generated secrets for MysticFunds - DO NOT COMMIT
JWT_SECRET=${JWT_SECRET}
API_KEY=${API_KEY}
DB_ENCRYPTION_KEY=${DB_ENCRYPTION_KEY}

# Add your Railway and Vercel tokens here
# RAILWAY_TOKEN=
# VERCEL_TOKEN=
# VERCEL_ORG_ID=
# VERCEL_PROJECT_ID=
EOF
    echo "Secrets saved to .env.local"
    echo "Make sure .env.local is in your .gitignore!"
fi