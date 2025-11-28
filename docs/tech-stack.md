# Technology Stack

## Overview

This document details every technology choice for The Asheville Setlist, including rationale, alternatives considered, and decision criteria.

---

## Backend Technologies

### Language: Go 1.22+

**Decision**: Go as the primary backend language

**Rationale**:
1. **Concurrency**: Goroutines perfect for scraping multiple venues simultaneously
2. **Performance**: Fast API response times for complex filtering queries
3. **Developer Experience**: Team already experienced with Go
4. **Type Safety**: Strong static typing prevents runtime errors
5. **Deployment**: Compiles to single binary, easy containerization
6. **Ecosystem**: Excellent web frameworks and database libraries

**Alternatives Considered**:
- **Python**: Better scraping libraries (BeautifulSoup, Scrapy) but slower performance, less suited for concurrent API
- **TypeScript/Node.js**: Full-stack JavaScript benefit, but not as performant for concurrent scraping
- **Java**: Verbose, heavier runtime, overkill for project size

**Why Go Won**:
- Best balance of performance, concurrency, and developer experience
- Single binary deployment simplifies Cloud Run containers
- Excellent for both API and scraper workloads

---

### API Framework: Gin (or Fiber/Echo)

**Decision**: Gin as the HTTP framework

**Candidates**:

| Framework | Performance | Features | Learning Curve | Community |
|-----------|------------|----------|----------------|-----------|
| **Gin** | ⭐⭐⭐⭐⭐ Fast | ⭐⭐⭐⭐ Good | ⭐⭐⭐⭐ Easy | ⭐⭐⭐⭐⭐ Large |
| **Fiber** | ⭐⭐⭐⭐⭐ Fastest | ⭐⭐⭐⭐⭐ Excellent | ⭐⭐⭐⭐ Easy | ⭐⭐⭐⭐ Growing |
| **Echo** | ⭐⭐⭐⭐ Fast | ⭐⭐⭐⭐ Good | ⭐⭐⭐⭐ Easy | ⭐⭐⭐⭐ Good |
| **net/http** | ⭐⭐⭐ Adequate | ⭐⭐ Basic | ⭐⭐⭐⭐⭐ Simple | ⭐⭐⭐⭐⭐ Standard lib |

**Recommendation: Gin**
- Most popular (44k GitHub stars)
- Excellent documentation
- Middleware ecosystem (CORS, logging, recovery)
- Fast enough (benchmarks show minimal difference vs Fiber for our scale)

**Example Usage**:
```go
package main

import (
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    r.GET("/api/shows", getShows)
    r.GET("/api/shows/:id", getShow)
    r.POST("/api/shows", createShow)

    r.Run(":8080")
}
```

**Alternative**: **Fiber** if you want Express.js-like API and fastest performance

---

### Database Driver: pgx

**Decision**: pgx as PostgreSQL driver

**Why pgx**:
- Fastest PostgreSQL driver for Go
- Native Go implementation (no CGO)
- Advanced features (connection pooling, prepared statements, COPY)
- Type-safe scanning
- Compatible with database/sql interface

**Alternative**: lib/pq (older, slower, maintenance mode)

**Installation**:
```go
import (
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)
```

---

### ORM/Query Builder: sqlc (Recommended) or GORM

**Decision**: sqlc for type-safe SQL

**Comparison**:

| Feature | sqlc | GORM |
|---------|------|------|
| **Type Safety** | ⭐⭐⭐⭐⭐ Compile-time | ⭐⭐⭐ Runtime |
| **Performance** | ⭐⭐⭐⭐⭐ Raw SQL | ⭐⭐⭐⭐ Good |
| **Learning Curve** | ⭐⭐⭐⭐ Easy (SQL) | ⭐⭐⭐ Moderate |
| **Flexibility** | ⭐⭐⭐⭐⭐ Full SQL | ⭐⭐⭐⭐ Query builder |
| **Migrations** | Manual | Built-in |

**sqlc Example**:
```sql
-- queries/shows.sql
-- name: GetShow :one
SELECT * FROM shows
WHERE id = $1;

-- name: ListShows :many
SELECT * FROM shows
WHERE date >= $1
ORDER BY date
LIMIT $2;
```

```go
// Generated code
show, err := queries.GetShow(ctx, showID)
shows, err := queries.ListShows(ctx, ListShowsParams{
    Date: startDate,
    Limit: 20,
})
```

**GORM Example**:
```go
// Manual struct mapping
type Show struct {
    ID      int       `gorm:"primaryKey"`
    Title   string
    Date    time.Time
    VenueID int
}

db.Where("date >= ?", startDate).Limit(20).Find(&shows)
```

**Recommendation**: **sqlc** for this project
- Write SQL, get type-safe Go code
- No magic, explicit queries
- Better for complex queries (our use case)
- Easier to optimize

**Alternative**: **GORM** if you prefer ORM patterns and want built-in migrations

---

### Web Scraping: Colly (primary) + chromedp (fallback)

**Decision**: Colly for most scraping, chromedp for JavaScript-heavy sites

**Colly** (Fast, static HTML):
```go
import "github.com/gocolly/colly/v2"

c := colly.NewCollector()

c.OnHTML(".event-item", func(e *colly.HTMLElement) {
    show := Show{
        Title: e.ChildText(".event-title"),
        Date:  parseDate(e.ChildText(".event-date")),
    }
    saveShow(show)
})

c.Visit("https://venue.com/events")
```

**Pros**:
- Fast (no browser overhead)
- Low memory usage
- Rate limiting built-in
- Concurrent scraping support

**chromedp** (JavaScript-rendered sites):
```go
import "github.com/chromedp/chromedp"

ctx, cancel := chromedp.NewContext(context.Background())
defer cancel()

var html string
chromedp.Run(ctx,
    chromedp.Navigate("https://venue.com/events"),
    chromedp.WaitVisible(".event-item"),
    chromedp.OuterHTML("body", &html),
)
```

**Pros**:
- Handles JavaScript-rendered content
- Can interact with pages (click, scroll)
- Real browser environment

**Cons**:
- Slower (runs headless Chrome)
- Higher memory usage

**Strategy**: Try Colly first, fallback to chromedp if needed

---

### Task Scheduling: go-co-op/gocron

**Decision**: gocron for in-process scheduling (if needed)

**Note**: Cloud Scheduler triggers Cloud Run Jobs, so local scheduling may not be needed. But if we want to schedule tasks within the API server:

```go
import "github.com/go-co-op/gocron"

s := gocron.NewScheduler(time.UTC)
s.Every(1).Hour().Do(func() {
    cleanupExpiredShows()
})
s.StartAsync()
```

**Alternative**: Rely solely on Cloud Scheduler (recommended for simplicity)

---

### Migrations: golang-migrate

**Decision**: golang-migrate for database migrations

```bash
# Install
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Create migration
migrate create -ext sql -dir migrations -seq init_schema

# Run migrations
migrate -path migrations -database "postgres://localhost/asheville" up
```

**Migration Files**:
```sql
-- migrations/000001_init_schema.up.sql
CREATE TABLE venues (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    region TEXT
);

-- migrations/000001_init_schema.down.sql
DROP TABLE venues;
```

**Alternative**: GORM AutoMigrate (less control, not recommended for production)

---

## Frontend Technologies

### Framework: Next.js 15

**Decision**: Next.js 15 with App Router

**Why Next.js**:
1. **SEO**: Server-side rendering crucial for concert discovery (Google indexing)
2. **Performance**: Automatic code splitting, image optimization, edge caching
3. **Developer Experience**: File-based routing, TypeScript support, hot reload
4. **Ecosystem**: Largest React ecosystem, solutions for everything
5. **Flexibility**: SSR, SSG, ISR, Client-side - choose per page

**Key Features for Our Project**:
- **Server Components**: Fetch data on server, reduce client JS
- **Streaming**: Progressive page rendering
- **Middleware**: Edge-based redirects, auth checks
- **API Routes**: Optional backend endpoints (though we use Go API)

**Alternatives Considered**:
- **SvelteKit**: Less boilerplate, easier to learn, but smaller ecosystem
- **Astro**: Best for content-heavy sites, but less interactive
- **Remix**: Great DX, but smaller community than Next.js

**Why Next.js Won**: SEO + ecosystem + industry adoption

---

### Language: TypeScript 5+

**Decision**: TypeScript (strict mode)

**Configuration** (`tsconfig.json`):
```json
{
  "compilerOptions": {
    "strict": true,
    "target": "ES2022",
    "lib": ["ES2022", "DOM"],
    "jsx": "preserve",
    "module": "ESNext",
    "moduleResolution": "bundler",
    "paths": {
      "@/*": ["./src/*"]
    }
  }
}
```

**Why TypeScript**:
- Matches Go's type safety approach
- Catch errors at compile-time
- Better IDE support (autocomplete, refactoring)
- Interfaces with backend API (type-safe fetch)

**Alternative**: JavaScript (faster to write, but error-prone)

---

### Styling: Tailwind CSS

**Decision**: Tailwind CSS utility-first framework

**Why Tailwind**:
- **Rapid development**: No context switching (HTML + CSS in one place)
- **Consistency**: Design system built-in (spacing, colors)
- **Responsive**: Mobile-first breakpoints
- **Performance**: Purges unused CSS (small bundle size)
- **Integration**: Works seamlessly with Next.js

**Example**:
```tsx
<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 p-6">
  <Card className="hover:shadow-lg transition-shadow">
    <h2 className="text-2xl font-bold text-gray-900 dark:text-white">
      Show Title
    </h2>
  </Card>
</div>
```

**Alternatives Considered**:
- **CSS Modules**: More traditional, but more verbose
- **Styled Components**: CSS-in-JS, runtime cost
- **Vanilla CSS**: Too much boilerplate

**Why Tailwind Won**: Best DX for rapid iteration, component-based design

---

### Component Library: shadcn/ui

**Decision**: shadcn/ui for UI components

**What is shadcn/ui**:
- NOT an npm package (copy-paste components into your code)
- Built on Radix UI (accessibility primitives)
- Styled with Tailwind
- Full customization (you own the code)

**Why shadcn/ui**:
- **Ownership**: Components live in your codebase
- **Customization**: Full control, no ejecting
- **Accessibility**: Built on Radix (ARIA, keyboard nav)
- **TypeScript**: Fully typed
- **Design**: Beautiful defaults

**Installation**:
```bash
npx shadcn@latest init
npx shadcn@latest add button card dialog
```

**Usage**:
```tsx
import { Button } from "@/components/ui/button"
import { Card } from "@/components/ui/card"

<Card>
  <Button variant="outline">Click Me</Button>
</Card>
```

**Alternatives Considered**:
- **Material UI**: Heavy, opinionated design
- **Chakra UI**: Great, but less customizable
- **Headless UI**: Good, but need to style everything

**Why shadcn/ui Won**: Best balance of quality, customization, and DX

---

### Forms: React Hook Form + Zod

**Decision**: React Hook Form for form state, Zod for validation

**Why React Hook Form**:
- Best performance (uncontrolled components)
- Less re-renders
- Built-in validation
- TypeScript support

**Why Zod**:
- Runtime + TypeScript validation
- Composable schemas
- Excellent error messages

**Example**:
```tsx
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { z } from "zod"

const showSchema = z.object({
  bandName: z.string().min(1, "Required"),
  showDate: z.date(),
  venueId: z.number(),
  ticketUrl: z.string().url().optional(),
})

type ShowForm = z.infer<typeof showSchema>

function SubmitShowForm() {
  const form = useForm<ShowForm>({
    resolver: zodResolver(showSchema),
  })

  const onSubmit = (data: ShowForm) => {
    // data is type-safe!
    submitShow(data)
  }

  return <form onSubmit={form.handleSubmit(onSubmit)}>...</form>
}
```

**Alternatives**:
- **Formik**: Older, more boilerplate
- **React Final Form**: Good, but smaller community

---

### Data Fetching: TanStack Query (React Query)

**Decision**: TanStack Query for server state management

**Why TanStack Query**:
- **Caching**: Automatic caching, no duplicate requests
- **Automatic Refetching**: Stale-while-revalidate
- **Optimistic Updates**: Update UI before server confirms
- **DevTools**: Inspect cache, queries, mutations
- **TypeScript**: Full type inference

**Example**:
```tsx
import { useQuery } from "@tanstack/react-query"

function ShowsList() {
  const { data, isLoading, error } = useQuery({
    queryKey: ['shows', filters],
    queryFn: () => fetch('/api/shows').then(r => r.json()),
    staleTime: 5 * 60 * 1000, // 5 minutes
  })

  if (isLoading) return <Skeleton />
  if (error) return <Error />

  return <div>{data.shows.map(show => <ShowCard {...show} />)}</div>
}
```

**Features We'll Use**:
- Query caching (reduce API calls)
- Infinite queries (infinite scroll shows)
- Mutations (band submissions)
- Optimistic updates (instant UI feedback)

**Alternatives**:
- **SWR**: Similar, but less features
- **Apollo Client**: GraphQL-focused (we're using REST)
- **RTK Query**: Redux-based (too heavy)

---

### Date Handling: date-fns

**Decision**: date-fns for date manipulation

**Why date-fns**:
- Functional, immutable
- Tree-shakeable (only import what you use)
- Great TypeScript support
- Intuitive API

**Example**:
```tsx
import { format, parseISO, addDays, isBefore } from 'date-fns'

const showDate = parseISO('2025-11-15T19:00:00Z')
format(showDate, 'EEEE, MMMM d, yyyy') // "Friday, November 15, 2025"

const tomorrow = addDays(new Date(), 1)
```

**Alternative**: Day.js (smaller, but less features)

---

## Database Technologies

### Database: PostgreSQL 16

**Decision**: PostgreSQL as primary database

**Why PostgreSQL** (detailed in previous discussion):
- Relational data model fits perfectly (shows ↔ venues ↔ bands ↔ genres)
- Complex JOINs for filtering
- JSONB for flexible scraped data
- Full-text search built-in
- Best Go ecosystem support

**Features We'll Use**:
- Foreign keys (referential integrity)
- Indexes (query performance)
- JSONB (variable scraped data)
- Full-text search (tsvector)
- Transactions (scraper upserts)

---

### Database Hosting: Neon

**Decision**: Neon serverless PostgreSQL

**Why Neon**:
- **Serverless**: Auto-scaling compute
- **Free tier**: 0.5GB storage, always available
- **Branching**: Database branches (like Git) for testing
- **Fast cold starts**: Resume in <500ms
- **Connection pooling**: Built-in PgBouncer

**Alternatives Considered**:
- **Supabase**: Great, but need more than just DB (auth, storage)
- **Railway PostgreSQL**: $5/month (no free tier)
- **Cloud SQL**: $15+/month (too expensive)
- **Self-hosted on VPS**: Extra operational burden

**Why Neon Won**: Best free tier + serverless model + PostgreSQL compatibility

---

## Hosting & Infrastructure

### Frontend Hosting: Vercel

**Decision**: Vercel for Next.js hosting

**Why Vercel**:
- Made by Next.js creators (best optimization)
- Automatic deployments (git push → deploy)
- Global CDN (low latency worldwide)
- Edge Functions (serverless at edge)
- Free tier is generous

**Vercel Free Tier**:
- 100GB bandwidth/month
- Unlimited deployments
- Custom domains
- HTTPS automatic
- Preview deployments (for PRs)

**Alternatives**:
- **Netlify**: Similar, good alternative
- **Cloudflare Pages**: Free, but less Next.js optimization
- **AWS Amplify**: More complex, overkill

---

### Backend Hosting: Google Cloud Run

**Decision**: Cloud Run for containerized Go services

**Why Cloud Run** (detailed in architecture):
- Serverless containers (no cluster management)
- Scale to zero (cost efficiency)
- Fast deployments
- Always Free tier (2M requests/month)
- Container-based (meets requirement)

**Cloud Run vs Alternatives**:

| Platform | Ease | Cost | Scalability | Control |
|----------|------|------|-------------|---------|
| **Cloud Run** | ⭐⭐⭐⭐⭐ | $0 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ |
| GKE Autopilot | ⭐⭐⭐ | $10+ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| AWS App Runner | ⭐⭐⭐⭐ | $5+ | ⭐⭐⭐⭐ | ⭐⭐⭐ |
| AWS ECS Fargate | ⭐⭐⭐ | $10+ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| Hetzner VPS | ⭐⭐⭐ | $4 | ⭐⭐ | ⭐⭐⭐⭐⭐ |

---

### Container Orchestration: None (Cloud Run is enough)

**Decision**: Use Cloud Run, NOT Kubernetes

**Why NOT Kubernetes (for now)**:
- Only 2 services (API + scraper)
- No complex inter-service communication
- Developer already knows K8s (not a learning goal)
- Operational overhead not justified
- Cloud Run is K8s-based (easy migration later)

**When to Consider Kubernetes**:
- 10+ microservices
- Need service mesh
- Complex networking requirements
- Want Kubernetes practice

---

## Development Tools

### Version Control: Git + GitHub

**Decision**: Git for version control, GitHub for hosting

**Repository Structure**:
```
- Monorepo OR
- Separate repos (frontend/ + backend/)
```

**Recommendation**: Monorepo for easier coordination

---

### CI/CD: Cloud Build (backend) + Vercel (frontend)

**Decision**: Cloud Build for backend deployments

**cloudbuild.yaml**:
```yaml
steps:
  # Build API
  - name: 'gcr.io/cloud-builders/docker'
    args: ['build', '-t', 'gcr.io/$PROJECT_ID/asheville-api', './backend/cmd/api']

  # Push to registry
  - name: 'gcr.io/cloud-builders/docker'
    args: ['push', 'gcr.io/$PROJECT_ID/asheville-api']

  # Deploy to Cloud Run
  - name: 'gcr.io/google.com/cloudsdktool/cloud-sdk'
    args:
      - 'run'
      - 'deploy'
      - 'asheville-api'
      - '--image=gcr.io/$PROJECT_ID/asheville-api'
      - '--region=us-central1'
```

**Alternative**: GitHub Actions (works for both frontend + backend)

---

### Code Quality

**Backend (Go)**:
- **Linter**: golangci-lint (aggregates 50+ linters)
- **Formatter**: gofmt (standard)
- **Testing**: go test + testify

**Frontend (TypeScript)**:
- **Linter**: ESLint (Next.js config)
- **Formatter**: Prettier
- **Testing**: Vitest + React Testing Library

---

## Infrastructure as Code (Optional)

### Terraform

**Decision**: Optional - use Terraform if you want IaC

**Example** (`terraform/main.tf`):
```hcl
resource "google_cloud_run_service" "api" {
  name     = "asheville-api"
  location = "us-central1"

  template {
    spec {
      containers {
        image = "gcr.io/project/asheville-api"

        resources {
          limits = {
            memory = "512Mi"
            cpu    = "1"
          }
        }
      }
    }
  }
}
```

**Alternative**: Manual setup with gcloud CLI (simpler for MVP)

---

## Monitoring & Observability

### Logging: Cloud Logging

**Decision**: Built-in Cloud Logging for Cloud Run

**Structured Logging** (Go):
```go
import "log/slog"

slog.Info("Show created",
    "show_id", show.ID,
    "venue_id", show.VenueID,
    "date", show.Date)
```

---

### Monitoring: Cloud Monitoring

**Decision**: Built-in Cloud Monitoring (formerly Stackdriver)

**Metrics**:
- Request count, latency, errors (automatic)
- Custom metrics (via OpenTelemetry)

---

## Summary Table

| Category | Technology | Why |
|----------|-----------|-----|
| **Backend Language** | Go 1.22+ | Concurrency, performance, type safety |
| **API Framework** | Gin | Popular, fast, good DX |
| **Database** | PostgreSQL 16 | Relational model, complex queries |
| **Database Driver** | pgx | Fastest Go driver |
| **Query Layer** | sqlc | Type-safe SQL |
| **Web Scraping** | Colly + chromedp | Fast static + JS support |
| **Migrations** | golang-migrate | Industry standard |
| **Frontend Framework** | Next.js 15 | SEO, ecosystem, performance |
| **Language** | TypeScript 5 | Type safety, IDE support |
| **Styling** | Tailwind CSS | Rapid development, consistent design |
| **Components** | shadcn/ui | Customizable, accessible |
| **Forms** | React Hook Form + Zod | Performance, validation |
| **Data Fetching** | TanStack Query | Caching, optimistic updates |
| **Date Library** | date-fns | Functional, tree-shakeable |
| **Frontend Host** | Vercel | Next.js optimized, free tier |
| **Backend Host** | Cloud Run | Serverless containers, $0 |
| **Database Host** | Neon | Serverless PostgreSQL, free tier |
| **CI/CD** | Cloud Build + Vercel | Integrated with platforms |

---

## Decision Framework for Future Choices

When evaluating new technologies:

1. **Does it solve a real problem?** (avoid shiny object syndrome)
2. **Is it mature?** (check community size, maintenance)
3. **Does it fit our stack?** (Go + TypeScript ecosystem)
4. **What's the learning curve?** (balance new skills vs. velocity)
5. **What's the cost?** (stay within budget)
6. **Can we migrate away if needed?** (avoid lock-in)

---

This tech stack prioritizes:
- ✅ Type safety (Go + TypeScript + sqlc + Zod)
- ✅ Developer experience (modern tools, great docs)
- ✅ Performance (Go concurrency, React optimization)
- ✅ Cost efficiency ($0/month free tiers)
- ✅ Scalability (serverless architecture)
- ✅ Maintainability (popular, well-supported tools)
