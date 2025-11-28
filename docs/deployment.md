# Deployment Guide

## Overview

This guide covers deploying The Asheville Setlist to production using:
- **Frontend**: Vercel (Next.js)
- **Backend API**: Google Cloud Run
- **Scraper**: Google Cloud Run Jobs
- **Database**: Neon PostgreSQL

**Target Cost**: $0/month (within free tiers)

---

## Prerequisites

### Required Accounts
1. **Google Cloud Platform** - [cloud.google.com](https://cloud.google.com)
2. **Vercel** - [vercel.com](https://vercel.com)
3. **Neon** - [neon.tech](https://neon.tech)
4. **GitHub** - [github.com](https://github.com) (code repository)
5. **Domain Registrar** - Namecheap, Porkbun, etc. (optional)

### Local Development Tools
```bash
# Install gcloud CLI
curl https://sdk.cloud.google.com | bash
exec -l $SHELL

# Install Vercel CLI
npm i -g vercel

# Install migrate (database migrations)
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Install Docker (for local testing)
# https://docs.docker.com/get-docker/
```

---

## Part 1: Database Setup (Neon)

### 1. Create Neon Account & Project

1. Sign up at [neon.tech](https://neon.tech)
2. Create new project: `asheville-setlist`
3. Choose region: `us-east-2` (close to GCP us-central1)
4. Copy connection string

### 2. Connection String Format

Neon provides a connection string:
```
postgres://username:password@ep-xxx.us-east-2.aws.neon.tech/neondb?sslmode=require
```

**Save this as `DATABASE_URL` environment variable**

### 3. Run Migrations

```bash
cd backend

# Set DATABASE_URL
export DATABASE_URL="postgres://username:password@ep-xxx.us-east-2.aws.neon.tech/neondb?sslmode=require"

# Run migrations
migrate -path migrations -database "$DATABASE_URL" up

# Verify
migrate -path migrations -database "$DATABASE_URL" version
```

### 4. Seed Initial Data (Optional)

```bash
# Create seed script: scripts/seed.sql
psql "$DATABASE_URL" < scripts/seed.sql
```

---

## Part 2: Google Cloud Platform Setup

### 1. Create GCP Project

```bash
# Login
gcloud auth login

# Create project
gcloud projects create asheville-setlist --name="The Asheville Setlist"

# Set as default
gcloud config set project asheville-setlist

# Enable required APIs
gcloud services enable run.googleapis.com
gcloud services enable cloudbuild.googleapis.com
gcloud services enable cloudscheduler.googleapis.com
gcloud services enable secretmanager.googleapis.com
```

### 2. Set Up Secret Manager (for DATABASE_URL)

```bash
# Create secret
echo -n "$DATABASE_URL" | gcloud secrets create database-url --data-file=-

# Grant Cloud Run access
gcloud secrets add-iam-policy-binding database-url \
  --member="serviceAccount:PROJECT_NUMBER-compute@developer.gserviceaccount.com" \
  --role="roles/secretmanager.secretAccessor"
```

**Note**: Replace `PROJECT_NUMBER` with your actual project number (find with `gcloud projects describe asheville-setlist`)

### 3. Configure Billing (Required for Cloud Run)

Even though you'll use free tier, you need a billing account:
1. Go to [console.cloud.google.com/billing](https://console.cloud.google.com/billing)
2. Set up billing account (won't be charged within free tier)
3. Link to `asheville-setlist` project

---

## Part 3: Deploy Backend API (Cloud Run)

### 1. Prepare Dockerfile

```dockerfile
# backend/cmd/api/Dockerfile
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /api /api

# Cloud Run sets PORT env var
ENV PORT=8080
EXPOSE 8080

CMD ["/api"]
```

### 2. Deploy API to Cloud Run

```bash
cd backend

# Build and deploy (Cloud Build handles Docker)
gcloud run deploy asheville-api \
  --source ./cmd/api \
  --region us-central1 \
  --platform managed \
  --allow-unauthenticated \
  --set-secrets DATABASE_URL=database-url:latest \
  --memory 512Mi \
  --cpu 1 \
  --timeout 300 \
  --max-instances 10 \
  --min-instances 0

# Output will show:
# Service URL: https://asheville-api-xxx-uc.a.run.app
```

### 3. Test API

```bash
# Get service URL
API_URL=$(gcloud run services describe asheville-api --region us-central1 --format 'value(status.url)')

# Test endpoint
curl $API_URL/api/health
```

### 4. Set Up Custom Domain (Optional)

```bash
# Map custom domain
gcloud run domain-mappings create \
  --service asheville-api \
  --domain api.ashevillesetlist.com \
  --region us-central1

# Follow DNS instructions to add records to your domain
```

---

## Part 4: Deploy Scraper (Cloud Run Job)

### 1. Prepare Scraper Dockerfile

```dockerfile
# backend/cmd/scraper/Dockerfile
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /scraper ./cmd/scraper

FROM alpine:latest

RUN apk --no-cache add ca-certificates chromium chromium-chromedriver

COPY --from=builder /scraper /scraper

CMD ["/scraper"]
```

### 2. Create Cloud Run Job

```bash
cd backend

# Create job
gcloud run jobs create asheville-scraper \
  --source ./cmd/scraper \
  --region us-central1 \
  --set-secrets DATABASE_URL=database-url:latest \
  --memory 1Gi \
  --cpu 1 \
  --max-retries 3 \
  --task-timeout 3600

# Test job manually
gcloud run jobs execute asheville-scraper --region us-central1

# Check logs
gcloud run jobs executions logs tail asheville-scraper --region us-central1
```

### 3. Schedule with Cloud Scheduler

```bash
# Create scheduler job (every 6 hours)
gcloud scheduler jobs create http scrape-venues \
  --location us-central1 \
  --schedule "0 */6 * * *" \
  --uri "https://us-central1-run.googleapis.com/apis/run.googleapis.com/v1/namespaces/asheville-setlist/jobs/asheville-scraper:run" \
  --http-method POST \
  --oauth-service-account-email asheville-setlist@appspot.gserviceaccount.com \
  --time-zone "America/New_York"

# Test scheduler
gcloud scheduler jobs run scrape-venues --location us-central1
```

**Cron Schedule Options**:
- Every 6 hours: `0 */6 * * *`
- Every 4 hours: `0 */4 * * *`
- Daily at 3 AM: `0 3 * * *`
- Twice daily (6 AM, 6 PM): `0 6,18 * * *`

---

## Part 5: Deploy Frontend (Vercel)

### 1. Connect GitHub Repository

1. Push code to GitHub
2. Go to [vercel.com/new](https://vercel.com/new)
3. Import your repository
4. Vercel auto-detects Next.js

### 2. Configure Environment Variables

In Vercel dashboard → Settings → Environment Variables:

```
NEXT_PUBLIC_API_URL=https://asheville-api-xxx-uc.a.run.app
```

### 3. Deploy

```bash
# Option 1: Deploy via Vercel dashboard (recommended)
# Push to GitHub → Auto-deploys

# Option 2: Deploy via CLI
cd frontend
vercel --prod
```

### 4. Custom Domain

1. Vercel dashboard → Settings → Domains
2. Add `ashevillesetlist.com`
3. Follow DNS instructions (add A/CNAME records)
4. Vercel handles SSL automatically

---

## Part 6: CI/CD Setup

### Option 1: Cloud Build (Backend)

Create `cloudbuild.yaml` in backend/:

```yaml
steps:
  # Build API
  - name: 'gcr.io/cloud-builders/docker'
    args:
      - 'build'
      - '-t'
      - 'gcr.io/$PROJECT_ID/asheville-api'
      - './cmd/api'

  # Push to Artifact Registry
  - name: 'gcr.io/cloud-builders/docker'
    args:
      - 'push'
      - 'gcr.io/$PROJECT_ID/asheville-api'

  # Deploy to Cloud Run
  - name: 'gcr.io/google.com/cloudsdktool/cloud-sdk'
    args:
      - 'run'
      - 'deploy'
      - 'asheville-api'
      - '--image=gcr.io/$PROJECT_ID/asheville-api'
      - '--region=us-central1'
      - '--platform=managed'

images:
  - 'gcr.io/$PROJECT_ID/asheville-api'
```

**Connect to GitHub**:
```bash
gcloud builds triggers create github \
  --repo-name=TheAshevilleSetlist \
  --repo-owner=YOUR_GITHUB_USERNAME \
  --branch-pattern=^main$ \
  --build-config=backend/cloudbuild.yaml
```

### Option 2: GitHub Actions (Both Frontend + Backend)

Create `.github/workflows/deploy.yml`:

```yaml
name: Deploy

on:
  push:
    branches: [main]

jobs:
  deploy-backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Setup Cloud SDK
        uses: google-github-actions/setup-gcloud@v1
        with:
          service_account_key: ${{ secrets.GCP_SA_KEY }}
          project_id: asheville-setlist

      - name: Deploy API
        run: |
          cd backend
          gcloud run deploy asheville-api \
            --source ./cmd/api \
            --region us-central1

  deploy-frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: vercel/deploy-action@v1
        with:
          vercel-token: ${{ secrets.VERCEL_TOKEN }}
          vercel-org-id: ${{ secrets.VERCEL_ORG_ID }}
          vercel-project-id: ${{ secrets.VERCEL_PROJECT_ID }}
```

---

## Part 7: Monitoring & Logging

### 1. Cloud Logging (Backend)

**View logs**:
```bash
# API logs
gcloud run logs read asheville-api --region us-central1 --limit 50

# Scraper logs
gcloud run jobs executions logs tail asheville-scraper --region us-central1
```

**Cloud Console**:
1. Go to [console.cloud.google.com/logs](https://console.cloud.google.com/logs)
2. Filter by service: `resource.type="cloud_run_revision"`

### 2. Cloud Monitoring

**Create Dashboard**:
1. Go to Monitoring → Dashboards
2. Add charts:
   - Request count
   - Latency (p50, p95, p99)
   - Error rate
   - Instance count

**Set Up Alerts**:
```bash
# Alert on high error rate
gcloud alpha monitoring policies create \
  --notification-channels=CHANNEL_ID \
  --display-name="API Error Rate High" \
  --condition-display-name="Error rate > 5%" \
  --condition-threshold-value=0.05 \
  --condition-threshold-duration=300s
```

### 3. Vercel Analytics

- Automatic analytics in Vercel dashboard
- View page views, performance metrics
- Optional: Upgrade for detailed analytics

---

## Part 8: Cost Monitoring

### 1. Set Up Budget Alerts

```bash
gcloud billing budgets create \
  --billing-account=BILLING_ACCOUNT_ID \
  --display-name="Asheville Setlist Budget" \
  --budget-amount=20 \
  --threshold-rule=percent=50 \
  --threshold-rule=percent=90 \
  --threshold-rule=percent=100
```

### 2. Monitor Usage

**Cloud Run**:
```bash
# Check request count
gcloud monitoring time-series list \
  --filter='metric.type="run.googleapis.com/request_count"'
```

**Neon**:
- Dashboard shows storage usage
- Free tier: 0.5GB limit

**Vercel**:
- Dashboard shows bandwidth usage
- Free tier: 100GB/month

---

## Environment Variables Summary

### Backend (Cloud Run)

| Variable | Value | Source |
|----------|-------|--------|
| `DATABASE_URL` | postgres://... | Secret Manager |
| `PORT` | 8080 | Cloud Run (auto-set) |
| `GIN_MODE` | release | Manual |

### Frontend (Vercel)

| Variable | Value | Source |
|----------|-------|--------|
| `NEXT_PUBLIC_API_URL` | https://asheville-api-xxx.run.app | Manual |

---

## Deployment Checklist

### Initial Setup
- [ ] Create Neon database
- [ ] Run migrations
- [ ] Create GCP project
- [ ] Enable required APIs
- [ ] Set up Secret Manager
- [ ] Configure billing

### Backend Deployment
- [ ] Build and deploy API to Cloud Run
- [ ] Test API endpoints
- [ ] Create Cloud Run Job for scraper
- [ ] Set up Cloud Scheduler
- [ ] Test scraper job

### Frontend Deployment
- [ ] Connect GitHub to Vercel
- [ ] Set environment variables
- [ ] Deploy to Vercel
- [ ] Test production site

### Post-Deployment
- [ ] Set up custom domains (optional)
- [ ] Configure monitoring
- [ ] Set up alerts
- [ ] Test end-to-end flow
- [ ] Document any issues

---

## Rollback Procedures

### Backend Rollback

```bash
# List revisions
gcloud run revisions list --service asheville-api --region us-central1

# Rollback to previous revision
gcloud run services update-traffic asheville-api \
  --to-revisions=asheville-api-00001-abc=100 \
  --region us-central1
```

### Frontend Rollback

**Vercel dashboard**:
1. Go to Deployments
2. Find previous deployment
3. Click "..." → "Promote to Production"

**CLI**:
```bash
vercel rollback
```

---

## Troubleshooting

### Issue: Cloud Run service won't start

**Check logs**:
```bash
gcloud run logs read asheville-api --region us-central1 --limit 50
```

**Common causes**:
- Database connection failed (check DATABASE_URL secret)
- Port mismatch (ensure listening on $PORT)
- Missing dependencies (check Dockerfile)

### Issue: Database connection timeout

**Verify connection**:
```bash
psql "$DATABASE_URL" -c "SELECT NOW();"
```

**Solutions**:
- Check Neon project is active
- Verify connection string format
- Ensure SSL mode is enabled (`?sslmode=require`)

### Issue: Scraper not running

**Check scheduler**:
```bash
gcloud scheduler jobs describe scrape-venues --location us-central1
```

**Test manually**:
```bash
gcloud run jobs execute asheville-scraper --region us-central1
```

### Issue: Frontend can't reach API

**Check CORS**:
```go
// In Go API
r.Use(cors.New(cors.Config{
    AllowOrigins: []string{"https://ashevillesetlist.com", "https://*.vercel.app"},
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
}))
```

**Verify environment variable**:
- Vercel dashboard → Settings → Environment Variables
- Ensure `NEXT_PUBLIC_API_URL` is correct

---

## Performance Optimization

### Backend

**Enable HTTP/2**:
- Cloud Run enables by default

**Connection pooling**:
```go
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
```

**Caching** (future):
- Add Redis (Cloud Memorystore) for API response caching

### Frontend

**Next.js optimization**:
```typescript
// next.config.js
module.exports = {
  images: {
    domains: ['cdn.example.com'],
    formats: ['image/avif', 'image/webp'],
  },
  experimental: {
    optimizeCss: true,
  },
}
```

**Enable ISR**:
```typescript
// Revalidate every 5 minutes
export const revalidate = 300;
```

---

## Scaling Beyond Free Tier

### When to Upgrade

**Neon**:
- Exceeding 0.5GB storage → Upgrade to Launch ($19/mo for 10GB)

**Cloud Run**:
- Exceeding 2M requests/month → Still free up to usage limits
- Need guaranteed instances → Set min-instances > 0 (~$20/mo)

**Vercel**:
- Exceeding 100GB bandwidth → Pro plan ($20/mo)

### Cost at Scale

**Example: 100K monthly users**

| Service | Usage | Cost |
|---------|-------|------|
| Vercel | 500GB bandwidth | $20/mo (Pro) |
| Cloud Run | 10M requests | $0 (within limits) |
| Neon | 3GB storage | $19/mo (Launch) |
| **Total** | | **$39/month** |

---

## Security Best Practices

### 1. API Security

```go
// Rate limiting (future)
import "golang.org/x/time/rate"

limiter := rate.NewLimiter(10, 20) // 10 req/sec, burst 20

// Middleware
if !limiter.Allow() {
    c.JSON(429, gin.H{"error": "Too many requests"})
    return
}
```

### 2. Database Security

- Use Secret Manager (never hardcode DATABASE_URL)
- Enable SSL (Neon requires by default)
- Parameterized queries (prevent SQL injection)

### 3. Frontend Security

**Next.js Security Headers**:
```typescript
// next.config.js
module.exports = {
  headers: async () => [
    {
      source: '/(.*)',
      headers: [
        {
          key: 'X-Frame-Options',
          value: 'DENY',
        },
        {
          key: 'X-Content-Type-Options',
          value: 'nosniff',
        },
      ],
    },
  ],
}
```

---

## Backup & Disaster Recovery

### Database Backups

**Neon automatic backups**:
- Daily snapshots (point-in-time recovery)
- Retained for 7 days (free tier)

**Manual backups**:
```bash
# Weekly backup to Cloud Storage
pg_dump "$DATABASE_URL" | gzip > backup-$(date +%Y%m%d).sql.gz

# Upload to GCS
gsutil cp backup-*.sql.gz gs://asheville-backups/
```

### Code Backups

- GitHub is source of truth
- Clone repository regularly

### Recovery Test

**Test recovery quarterly**:
1. Restore database from backup
2. Deploy services from git
3. Verify end-to-end functionality
4. Document recovery time (target: <30 minutes)

---

## Production Readiness Checklist

- [ ] All services deployed and tested
- [ ] Custom domains configured
- [ ] HTTPS enabled (automatic)
- [ ] Database migrations run
- [ ] Monitoring dashboards created
- [ ] Alerts configured
- [ ] Budget alerts set up
- [ ] Backup strategy tested
- [ ] Documentation updated
- [ ] Runbook created for common issues

---

## Next Steps

After successful deployment:
1. Monitor logs and metrics for first week
2. Optimize queries based on actual usage
3. Add more venues to scraper config
4. Implement user authentication (Phase 2)
5. Set up staging environment

---

## Support Resources

- **Google Cloud Support**: [cloud.google.com/support](https://cloud.google.com/support)
- **Vercel Support**: [vercel.com/support](https://vercel.com/support)
- **Neon Support**: [neon.tech/docs](https://neon.tech/docs)
- **Community**: Discord, Stack Overflow

---

This deployment architecture supports zero-cost MVP launch with clear upgrade path as the project grows.
