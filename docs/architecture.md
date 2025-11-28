# System Architecture

## Overview

The Asheville Setlist is a cloud-native, microservices-oriented application built on Google Cloud Platform, leveraging serverless containers for scalability and cost efficiency.

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Users / Browsers                      │
│            (Desktop, Mobile, Tablets)                    │
└────────────────────────┬────────────────────────────────┘
                         │
                         │ HTTPS
                         ▼
┌─────────────────────────────────────────────────────────┐
│                  Vercel Edge Network                     │
│                                                          │
│  ┌────────────────────────────────────────────────┐    │
│  │         Next.js 15 Application                 │    │
│  │                                                 │    │
│  │  • Server-Side Rendering (SSR)                │    │
│  │  • Static Site Generation (SSG)               │    │
│  │  • API Routes (optional)                      │    │
│  │  • Image Optimization                         │    │
│  │  • Edge Caching                               │    │
│  └────────────────────────────────────────────────┘    │
│                                                          │
│  Global CDN: Low latency worldwide                      │
└────────────────────────┬────────────────────────────────┘
                         │
                         │ HTTPS (api.ashevillesetlist.com)
                         ▼
┌─────────────────────────────────────────────────────────┐
│              Google Cloud Platform                       │
│              us-central1 region                          │
│                                                          │
│  ┌──────────────────────────────────────────────────┐  │
│  │  Cloud Run Service: asheville-api                │  │
│  │                                                   │  │
│  │  Container: Go API Server                        │  │
│  │  • REST API endpoints                            │  │
│  │  • Auto-scaling (0-N instances)                  │  │
│  │  • Resources: 512MB RAM, 1 vCPU                  │  │
│  │  • Timeout: 5 minutes                            │  │
│  │  • Concurrency: 80 requests/instance             │  │
│  │                                                   │  │
│  │  Endpoints:                                       │  │
│  │  • GET /api/shows                                │  │
│  │  • GET /api/shows/:id                            │  │
│  │  • GET /api/venues                               │  │
│  │  • GET /api/bands                                │  │
│  │  • GET /api/search                               │  │
│  │  • POST /api/bands/:id/similar                   │  │
│  │  • POST /api/shows (band submission)             │  │
│  └────────────────────┬─────────────────────────────┘  │
│                       │                                 │
│  ┌────────────────────┴─────────────────────────────┐  │
│  │  Cloud Run Job: asheville-scraper                │  │
│  │                                                   │  │
│  │  Container: Go Scraper Service                   │  │
│  │  • Triggered by Cloud Scheduler                  │  │
│  │  • Schedule: "0 */6 * * *" (every 6 hours)       │  │
│  │  • Resources: 1GB RAM, 1 vCPU                    │  │
│  │  • Timeout: 60 minutes                           │  │
│  │  • Max retries: 3                                │  │
│  │                                                   │  │
│  │  Process:                                         │  │
│  │  1. Fetch configured venue URLs                  │  │
│  │  2. Scrape show data concurrently (goroutines)   │  │
│  │  3. Parse and normalize data                     │  │
│  │  4. Upsert to PostgreSQL                         │  │
│  │  5. Log results to Cloud Logging                 │  │
│  └────────────────────┬─────────────────────────────┘  │
│                       │                                 │
│  ┌────────────────────┴─────────────────────────────┐  │
│  │  Cloud Scheduler                                  │  │
│  │                                                   │  │
│  │  • Cron job: scrape-venues                       │  │
│  │  • Schedule: Every 6 hours                       │  │
│  │  • Invokes: asheville-scraper                    │  │
│  └───────────────────────────────────────────────────┘  │
│                                                          │
│  ┌───────────────────────────────────────────────────┐  │
│  │  Cloud Logging & Monitoring                       │  │
│  │                                                   │  │
│  │  • Centralized logs (API + Scraper)              │  │
│  │  • Metrics & dashboards                          │  │
│  │  • Alerts & notifications                        │  │
│  │  • Request tracing                               │  │
│  └───────────────────────────────────────────────────┘  │
└────────────────────────┬────────────────────────────────┘
                         │
                         │ PostgreSQL wire protocol (TLS)
                         ▼
┌─────────────────────────────────────────────────────────┐
│                  Neon PostgreSQL                         │
│                  Serverless Database                     │
│                                                          │
│  • Region: us-east-2 (low latency to GCP)               │
│  • Auto-scaling compute                                 │
│  • Always-available storage                             │
│  • Automatic backups                                    │
│  • Connection pooling (PgBouncer)                       │
│  • Free tier: 0.5GB storage                             │
│                                                          │
│  Database: asheville_setlist                            │
│  Tables: shows, bands, venues, genres,                  │
│          band_genres, show_bands, articles              │
└─────────────────────────────────────────────────────────┘
```

---

## Component Details

### Frontend: Vercel + Next.js

**Deployment Model**: Edge-optimized static + server-rendered pages

**Key Features**:
- **SSR for concert listings**: Fresh data on every request for SEO
- **SSG for static pages**: About, FAQ, etc. (build-time generation)
- **ISR (Incremental Static Regeneration)**: Cached pages with revalidation
- **Image Optimization**: Automatic resizing, WebP conversion
- **Edge Functions**: Run serverless functions at edge locations

**Routing**:
```
/ (homepage)                    → SSR (show featured concerts)
/shows                          → SSR (list all shows, filters)
/shows/[id]                     → SSR (show details, SEO critical)
/venues                         → SSR (venue listings)
/venues/[id]                    → SSR (venue page with shows)
/bands/[id]                     → SSR (band page, similar artists)
/search                         → SSR (search results)
/articles                       → ISR (blog listing, revalidate: 3600s)
/articles/[slug]                → ISR (article content)
```

**Data Fetching Pattern**:
```typescript
// Server Component (default in Next.js 15)
export default async function ShowsPage() {
  const shows = await fetch('https://api.ashevillesetlist.com/api/shows', {
    next: { revalidate: 300 } // Cache for 5 minutes
  }).then(r => r.json());

  return <ShowsList shows={shows} />;
}
```

**Client-Side Interactivity**:
- Filters: React state + URL params
- Search: Debounced input → API call
- Infinite scroll: TanStack Query
- Favorites: Local storage (future: user accounts)

---

### Backend: Cloud Run API

**Language**: Go 1.22+
**Framework**: Gin (or Fiber/Echo)

**Service Configuration**:
```yaml
Service: asheville-api
Region: us-central1
CPU: 1
Memory: 512Mi
Min instances: 0 (scale to zero)
Max instances: 10
Timeout: 5 minutes
Concurrency: 80
Port: 8080
```

**API Design**:

**RESTful endpoints**:

```
GET    /api/shows                     # List shows with filters
GET    /api/shows/:id                 # Show details
GET    /api/venues                    # List venues
GET    /api/venues/:id                # Venue details
GET    /api/venues/:id/shows          # Shows at venue
GET    /api/bands                     # List bands
GET    /api/bands/:id                 # Band details
GET    /api/bands/:id/similar         # Similar bands (genre-based)
GET    /api/genres                    # List genres
GET    /api/search?q=moon+taxi        # Search shows/bands/venues
POST   /api/shows                     # Submit show (band submission)
GET    /api/articles                  # Blog articles
GET    /api/articles/:slug            # Article content
```

**Query Parameters** (for /api/shows):
```
?date_start=2025-11-10         # Date range filtering
?date_end=2025-11-20
?venue_id=5                    # Filter by venue
?genre=rock,indie              # Filter by genres (comma-separated)
?region=downtown               # Filter by venue region
?price_max=25                  # Price range
?page=1                        # Pagination
?limit=20                      # Results per page
?sort=date                     # Sort: date, price, venue
```

**Response Format**:
```json
{
  "data": [...],
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "total_pages": 8
  }
}
```

**Error Handling**:
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid date format",
    "details": { "field": "date_start" }
  }
}
```

**Middleware Stack**:
1. CORS (allow Vercel frontend origin)
2. Request logging (structured JSON logs)
3. Authentication (future: JWT verification)
4. Rate limiting (future: Redis-based)
5. Panic recovery

---

### Scraper: Cloud Run Job

**Language**: Go 1.22+
**Libraries**: Colly (fast scraper) or chromedp (JS-heavy sites)

**Service Configuration**:
```yaml
Job: asheville-scraper
Region: us-central1
CPU: 1
Memory: 1Gi
Timeout: 60 minutes
Max retries: 3
Execution environment: Second generation
```

**Scraping Strategy**:

1. **Configuration-driven**:
   - Store venue scraping configs in database (table: venue_scrapers)
   - Each config: URL, CSS selectors, parsing rules

2. **Concurrent scraping**:
   ```go
   // Pseudo-code
   venues := getActiveVenues()
   var wg sync.WaitGroup

   for _, venue := range venues {
       wg.Add(1)
       go func(v Venue) {
           defer wg.Done()
           scrapeVenue(v)
       }(venue)
   }
   wg.Wait()
   ```

3. **Resilience**:
   - Retry failed scrapes (exponential backoff)
   - Log errors to Cloud Logging
   - Continue on individual failures (don't fail entire job)

4. **Data normalization**:
   - Parse dates (various formats → ISO 8601)
   - Extract band names (regex patterns)
   - Match existing bands (fuzzy matching)
   - Create new bands if not found

5. **Deduplication**:
   - Check if show already exists (venue + date + bands)
   - Update if changed, skip if identical

**Example Venue Config**:
```json
{
  "venue_id": 1,
  "name": "The Orange Peel",
  "url": "https://theorangepeel.net/events",
  "scraper_type": "static", // or "javascript"
  "selectors": {
    "container": ".event-item",
    "title": ".event-title",
    "date": ".event-date",
    "bands": ".lineup .band",
    "price": ".ticket-price"
  },
  "active": true
}
```

**Scheduled Execution**:
- Cloud Scheduler triggers job every 6 hours
- Alternative: Webhook-based (venue posts update)

---

### Database: Neon PostgreSQL

**Plan**: Free tier (0.5GB storage, always available)
**Region**: us-east-2 (close to GCP us-central1)
**Connection**: PgBouncer pooling (connection limit friendly)

**Schema Design**: See `database-schema.md`

**Connection Pattern**:
```go
// API server: Connection pool
db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)

// Scraper: Single connection (job-based, no pooling needed)
db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
defer db.Close()
```

**Performance Optimizations**:
- Indexes on frequently queried columns (date, venue_id, genre_id)
- GIN index on JSONB columns (scraped_data)
- Materialized views for complex queries (future)
- Full-text search index (tsvector) for search

---

## Data Flow

### User Browse Flow

```
User visits /shows
  ↓
Next.js SSR fetches data from Cloud Run API
  ↓
API queries PostgreSQL (with filters)
  ↓
PostgreSQL returns results
  ↓
API formats response (JSON)
  ↓
Next.js renders HTML
  ↓
Browser displays page (with client-side hydration)
  ↓
User applies filters → Client-side fetch
  ↓
TanStack Query caches results
```

### Scraper Flow

```
Cloud Scheduler triggers (every 6 hours)
  ↓
Cloud Run Job starts scraper container
  ↓
Scraper reads venue configs from DB
  ↓
Concurrent scraping (goroutines)
  ↓
For each venue:
  - Fetch HTML
  - Parse with selectors
  - Extract show data
  - Normalize data
  ↓
Match/create bands
  ↓
Upsert shows to PostgreSQL
  ↓
Log results to Cloud Logging
  ↓
Container exits (success/failure)
```

### Band Submission Flow

```
User fills form on frontend
  ↓
Form validation (React Hook Form + Zod)
  ↓
POST /api/shows
  ↓
API validates data (Go validation)
  ↓
Check for duplicates
  ↓
Insert to PostgreSQL (pending approval status)
  ↓
Return success
  ↓
Frontend shows confirmation
  ↓
(Future: Admin reviews and approves)
```

---

## Deployment Architecture

### CI/CD Pipeline

**Frontend (Vercel)**:
```
git push to main
  ↓
Vercel webhook triggered
  ↓
Vercel builds Next.js app
  ↓
Deploy to edge network
  ↓
Automatic HTTPS, CDN distribution
```

**Backend (Cloud Run)**:
```
git push to main
  ↓
GitHub webhook → Cloud Build
  ↓
Cloud Build triggers (cloudbuild.yaml)
  ↓
Build Go container (multi-stage Dockerfile)
  ↓
Push to Artifact Registry
  ↓
Deploy to Cloud Run (rolling update)
  ↓
Health check passes
  ↓
Route traffic to new revision
```

**Infrastructure**:
- Option 1: Manual setup (gcloud CLI)
- Option 2: Terraform (Infrastructure as Code)

---

## Security Architecture

### Authentication & Authorization

**Current (MVP)**: Public read-only API
**Future**: JWT-based authentication

```
User login → JWT issued
  ↓
Frontend stores JWT (httpOnly cookie)
  ↓
API requests include JWT
  ↓
API validates JWT (middleware)
  ↓
Authorized requests proceed
```

### Network Security

- **HTTPS only**: Enforced on Vercel + Cloud Run
- **CORS**: Whitelist frontend origin
- **Rate limiting**: (Future) Cloud Armor or application-level
- **Secrets management**: Google Secret Manager
- **Database**: TLS-encrypted connections

### Data Security

- **Input validation**: Both frontend (Zod) and backend (Go validators)
- **SQL injection prevention**: Parameterized queries (sqlc/GORM)
- **XSS prevention**: Next.js auto-escaping, CSP headers
- **CSRF protection**: (Future) CSRF tokens for mutations

---

## Monitoring & Observability

### Logging

**Cloud Logging** (formerly Stackdriver):
- All Cloud Run logs automatically collected
- Structured JSON logging
- Log levels: INFO, WARN, ERROR, DEBUG
- Queryable with filters

**Example log query**:
```
resource.type="cloud_run_revision"
resource.labels.service_name="asheville-api"
severity="ERROR"
timestamp>="2025-11-08T00:00:00Z"
```

### Metrics

**Built-in Cloud Run metrics**:
- Request count
- Request latency (p50, p95, p99)
- Error rate
- Instance count
- CPU/Memory utilization

**Custom metrics** (future):
- Shows scraped per job
- API endpoint usage
- Search query patterns

### Alerting

**Cloud Monitoring alerts**:
- API error rate > 5%
- Scraper job failures
- Database connection errors
- High response latency (p95 > 1s)

### Tracing

**Cloud Trace** (future):
- Distributed tracing across services
- Request flow visualization
- Performance bottleneck identification

---

## Scalability Considerations

### Current Scale (MVP)

- **Shows**: ~5,000 concerts/year
- **Bands**: ~2,000 unique artists
- **Venues**: ~20-30 venues
- **Users**: ~10,000 monthly visitors
- **API Requests**: ~50,000/month

**Cloud Run handles this easily within free tier.**

### Growth Scale (Year 2+)

- **Shows**: ~20,000 concerts/year
- **Bands**: ~10,000 unique artists
- **Venues**: ~100 venues (expand to nearby cities)
- **Users**: ~100,000 monthly visitors
- **API Requests**: ~1M/month

**Scaling strategies**:
1. Increase Cloud Run max instances
2. Add Redis caching layer (reduce DB load)
3. Database read replicas (Neon supports)
4. CDN caching for API responses (Cloudflare)
5. Consider moving to paid Neon tier (3GB storage)

### Potential Bottlenecks

1. **Database**: Neon free tier 0.5GB limit
   - **Solution**: Upgrade to paid tier ($19/mo for 3GB)

2. **API rate limits**: Free tier 2M Cloud Run requests
   - **Solution**: Implement caching, stays in free tier

3. **Scraper**: Long-running jobs (>60 min timeout)
   - **Solution**: Break into smaller jobs per venue

---

## Disaster Recovery

### Backup Strategy

**Database**:
- Neon automatic backups (daily snapshots)
- Manual exports: `pg_dump` weekly to Cloud Storage

**Code**:
- Git repository (GitHub) is source of truth

**Infrastructure**:
- Terraform state (if using IaC)
- Documented manual setup steps

### Recovery Procedures

**Database failure**:
1. Restore from Neon snapshot (point-in-time recovery)
2. Or restore from Cloud Storage backup

**API service failure**:
1. Cloud Run auto-restarts failed instances
2. If persistent: Rollback to previous revision

**Total platform failure**:
1. Restore database from backup
2. Redeploy services (git + Cloud Build)
3. Update DNS if needed
4. Total recovery time: ~30 minutes

---

## Cost Projections

### MVP (Within Free Tiers)

| Service | Usage | Cost |
|---------|-------|------|
| Vercel | <100GB bandwidth | $0 |
| Cloud Run (API) | 30K requests/mo | $0 |
| Cloud Run (Scraper) | 120 executions/mo | $0 |
| Cloud Scheduler | 1 job | $0 |
| Neon PostgreSQL | 0.5GB storage | $0 |
| **Total** | | **$0/month** |

### Growth (Year 2)

| Service | Usage | Cost |
|---------|-------|------|
| Vercel | 500GB bandwidth | $0 (still free) |
| Cloud Run (API) | 500K requests/mo | $0 (within free tier) |
| Cloud Run (Scraper) | 120 executions/mo | $0 |
| Cloud Scheduler | 1 job | $0 |
| Neon PostgreSQL | 3GB storage (paid plan) | $19/mo |
| Domain | ashevillesetlist.com | $1/mo |
| **Total** | | **$20/month** |

---

## Technology Choices Recap

| Component | Technology | Why |
|-----------|-----------|-----|
| Frontend Framework | Next.js 15 | SEO (SSR), React ecosystem, TypeScript |
| Backend Language | Go | Concurrency, performance, dev experience |
| Database | PostgreSQL | Relational data, complex queries, JSONB flexibility |
| Frontend Hosting | Vercel | Next.js optimized, free tier, global CDN |
| Backend Hosting | Cloud Run | Serverless containers, free tier, no cluster mgmt |
| Scraper | Cloud Run Jobs | Batch processing, scheduled execution |
| Styling | Tailwind CSS | Utility-first, rapid development |
| Components | shadcn/ui | Customizable, accessible, TypeScript |
| ORM | sqlc or GORM | Type-safe queries, Go-first |

---

## Future Enhancements

### Phase 2 Features
- User authentication (Auth0 or Clerk)
- User profiles (favorites, followed bands)
- Email notifications (new shows for followed bands)
- Band/venue claim system (verified accounts)

### Phase 3 Features
- Mobile app (React Native or Flutter)
- Real-time updates (WebSockets or SSE)
- Ticketing integration (direct purchase links)
- Recommendations (ML-based, not just genre matching)

### Advanced Infrastructure
- Multi-region deployment (if expanding beyond Asheville)
- GraphQL API (Apollo Server)
- Elasticsearch for advanced search
- Redis caching layer
- Kubernetes migration (GKE Autopilot) if microservices grow

---

## Conclusion

This architecture balances:
- ✅ **Simplicity**: Managed services, minimal DevOps
- ✅ **Cost**: $0/month for MVP, <$20/month at scale
- ✅ **Scalability**: Cloud Run auto-scales, Neon serverless
- ✅ **Developer Experience**: Fast iteration, modern tools
- ✅ **Production Quality**: Industry-standard technologies

The architecture supports growth from local project to regional platform without major rewrites.
