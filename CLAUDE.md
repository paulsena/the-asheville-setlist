# The Asheville Setlist - Project Context

> **Project Overview**: A concert discovery platform for Asheville, NC's music scene featuring automated venue scraping, rich search/filtering, and community content.

---

## Quick Reference

| Attribute | Value |
|-----------|-------|
| **Project Type** | Full-stack web application |
| **Stage** | Planning & Architecture |
| **Primary Goal** | Help music lovers discover local concerts |
| **Task Tracking** | See `TASKS.md` for implementation checklist |

---

## AI Agent Instructions

> **READ THIS FIRST** if you are Claude Code or another AI agent.

### Before Starting Any Work

1. **Read `TASKS.md`** to understand current progress and next tasks
2. **Check the current phase** - don't skip ahead
3. **Read relevant docs** before implementing:
   - `docs/database-schema.md` for any database work
   - `docs/api-spec.md` for API implementation (create if missing)
   - `docs/architecture.md` for system design decisions
   - `docs/tech-stack.md` for technology choices

### Working on Tasks

1. **One task at a time** - complete fully before moving on
2. **Follow acceptance criteria** - all must pass
3. **Update task status** in `TASKS.md`:
   - `[~]` when starting
   - `[x]` when complete
   - `[!]` if blocked (add note)
4. **Commit after each task** with format: `[TASK-XXX] Description`

### Code Standards

#### Go (Backend)
```
- Use standard Go project layout
- Error handling: return errors, don't panic
- Logging: use structured logging (log/slog)
- Testing: table-driven tests preferred
- Naming: follow Go conventions (MixedCaps, not snake_case)
```

#### TypeScript (Frontend)
```
- Strict mode enabled
- Use interfaces over types where possible
- Prefer named exports
- Components: PascalCase files
- Utilities: camelCase files
- No `any` types (use `unknown` if needed)
```

#### SQL
```
- Use lowercase keywords (select, from, where)
- Table names: plural, snake_case (shows, band_genres)
- Column names: snake_case (created_at, venue_id)
- Always use parameterized queries
```

### What NOT to Do

- ❌ Don't add features not in the task list
- ❌ Don't change tech stack decisions without discussion
- ❌ Don't skip writing tests for API endpoints
- ❌ Don't use GraphQL (we chose REST)
- ❌ Don't add authentication until Phase 2+
- ❌ Don't over-engineer (this is a local concert site)

### When Uncertain

1. Check existing docs first
2. Look at similar completed code in the project
3. If still unclear, mark task with `[?]` and add question
4. Prefer simpler solutions over clever ones

---

## Tech Stack (Approved - Do Not Change)

### Backend
| Component | Technology | Notes |
|-----------|-----------|-------|
| Language | Go 1.22+ | Already experienced |
| Framework | Gin | Popular, fast, good docs |
| Database | PostgreSQL 16 | Relational + JSONB |
| Query Layer | sqlc | Type-safe SQL |
| Migrations | golang-migrate | Standard tool |
| Scraping | Colly + chromedp | Static + JS sites |

### Frontend
| Component | Technology | Notes |
|-----------|-----------|-------|
| Framework | Next.js 15 | App Router, TypeScript |
| Styling | Tailwind CSS | Utility-first |
| Components | shadcn/ui | Copy-paste, customizable |
| Forms | React Hook Form + Zod | Performance + validation |
| Data Fetching | TanStack Query | Caching, mutations |
| Icons | Lucide React | Consistent icon set |

### Infrastructure
| Component | Technology | Notes |
|-----------|-----------|-------|
| Frontend Host | Vercel | Free tier |
| Backend Host | Cloud Run | Free tier, serverless |
| Database Host | Neon | Free tier, serverless |
| Scheduler | Cloud Scheduler | Triggers scraper |
| CI/CD | Cloud Build + Vercel | Auto-deploy |

---

## Project Structure

```
asheville-setlist/
├── CLAUDE.md              # This file - project context
├── TASKS.md               # Task tracking for implementation
├── backend/
│   ├── cmd/
│   │   ├── api/           # API server entry point
│   │   │   ├── main.go
│   │   │   └── Dockerfile
│   │   └── scraper/       # Scraper entry point
│   │       ├── main.go
│   │       └── Dockerfile
│   ├── internal/
│   │   ├── config/        # Configuration loading
│   │   ├── db/            # sqlc generated code
│   │   ├── handlers/      # HTTP handlers
│   │   ├── middleware/    # Gin middleware
│   │   ├── models/        # Domain models
│   │   └── scraper/       # Scraping logic
│   ├── migrations/        # SQL migrations
│   ├── queries/           # sqlc query files
│   ├── go.mod
│   ├── go.sum
│   └── sqlc.yaml
├── frontend/
│   ├── app/               # Next.js App Router pages
│   ├── components/        # React components
│   │   ├── ui/            # shadcn/ui components
│   │   └── ...            # Custom components
│   ├── lib/               # Utilities, API client
│   ├── public/            # Static assets
│   ├── package.json
│   ├── tailwind.config.ts
│   └── tsconfig.json
├── docs/
│   ├── architecture.md
│   ├── database-schema.md
│   ├── tech-stack.md
│   ├── deployment.md
│   ├── api-spec.md        # TODO: Create this
│   ├── frontend-structure.md  # TODO: Create this
│   └── scraper-config.md  # TODO: Create this
├── seeds/
│   ├── genres.sql         # TODO: Create this
│   ├── venues.sql         # TODO: Create this
│   └── test-data.sql      # TODO: Create this
├── docker-compose.yml     # Local development
├── Makefile               # Common commands
└── .env.example           # Environment template
```

---

## Database Schema (Summary)

> Full details in `docs/database-schema.md`

### Core Tables
- **venues** - Concert locations (name, address, region, capacity)
- **bands** - Artists (name, bio, spotify_url, genres via junction)
- **shows** - Events (venue, date, price, ticket_url, bands via junction)
- **genres** - Music categories (name, slug)

### Junction Tables
- **show_bands** - Links shows to bands (is_headliner, performance_order)
- **band_genres** - Links bands to genres

### Support Tables
- **venue_scrapers** - Scraping configuration per venue
- **articles** - Blog content (minor feature)
- **users** - Future: authentication

### Key Patterns
- Slugs for SEO-friendly URLs (`the-orange-peel`, `moon-taxi`)
- JSONB `metadata` columns for flexible data
- `scraped_data` JSONB for raw scraper output
- `status` enum on shows (`scheduled`, `cancelled`, `postponed`)
- `source` tracking (`scraped`, `band_submitted`, `manual`)

---

## API Endpoints (Summary)

> Full details to be added in `docs/api-spec.md`

### Shows
```
GET  /api/shows                 # List with filters
GET  /api/shows/:id             # Show details
POST /api/shows                 # Band submission
```

### Venues
```
GET  /api/venues                # List all
GET  /api/venues/:slug          # Venue details + shows
```

### Bands
```
GET  /api/bands                 # List with pagination
GET  /api/bands/:slug           # Band details
GET  /api/bands/:slug/similar   # Similar bands by genre
```

### Other
```
GET  /api/genres                # List all genres
GET  /api/search?q=             # Global search
GET  /health                    # Health check
```

---

## Key Design Decisions

### Why PostgreSQL over MongoDB?
- Complex filtering requires JOINs (PostgreSQL excels)
- Similar artist feature needs genre relationships
- JSONB provides flexibility where needed
- Better Go tooling (sqlc)

### Why Next.js over Remix/SvelteKit?
- SEO critical (SSR is mature)
- Largest React ecosystem
- TypeScript first-class
- Vercel hosting optimized

### Why Cloud Run over Kubernetes?
- Only 2 services (API + scraper)
- No complex networking needed
- Scale-to-zero saves money
- Already know K8s (not a learning goal)

### Why sqlc over GORM?
- Type-safe generated code
- Write SQL, not ORM abstractions
- Easier to optimize queries
- Better for complex JOINs

---

## Environment Variables

### Backend
```bash
# Required
DATABASE_URL=postgres://user:pass@host:5432/db?sslmode=require
PORT=8080

# Optional
GIN_MODE=release  # or debug
LOG_LEVEL=info    # debug, info, warn, error
```

### Frontend
```bash
# Required
NEXT_PUBLIC_API_URL=https://api.ashevillesetlist.com

# Optional (Vercel auto-sets these)
VERCEL_ENV=production
```

---

## Common Commands

```bash
# Development
make dev              # Start all services (Docker Compose)
make api              # Start API only
make frontend         # Start frontend only
make scraper          # Run scraper once

# Database
make migrate          # Run migrations up
make migrate-down     # Rollback one migration
make seed             # Load seed data
make db-reset         # Drop, migrate, seed

# Testing
make test             # Run all tests
make test-api         # Backend tests only
make test-frontend    # Frontend tests only

# Code Generation
make sqlc             # Generate sqlc code
make lint             # Run linters

# Deployment
make deploy-api       # Deploy API to Cloud Run
make deploy-scraper   # Deploy scraper job
```

---

## Development Workflow

### Starting a New Feature

1. Check `TASKS.md` for next task
2. Create feature branch: `git checkout -b task-xxx-description`
3. Implement with tests
4. Update `TASKS.md` status
5. Commit: `git commit -m "[TASK-XXX] Description"`
6. Push and create PR (if using PRs)

### Local Development Setup

```bash
# 1. Clone repository
git clone https://github.com/[user]/asheville-setlist
cd asheville-setlist

# 2. Copy environment file
cp .env.example .env

# 3. Start services
make dev

# 4. Run migrations
make migrate

# 5. Load seed data
make seed

# 6. Access:
#    - Frontend: http://localhost:3000
#    - API: http://localhost:8080
#    - Database: localhost:5432
```

---

## External Resources

### Documentation
- [Next.js App Router](https://nextjs.org/docs/app)
- [Gin Web Framework](https://gin-gonic.com/docs/)
- [sqlc Documentation](https://docs.sqlc.dev/)
- [shadcn/ui Components](https://ui.shadcn.com/)
- [TanStack Query](https://tanstack.com/query/latest)

### Asheville Venues (for scraper reference)
- [The Orange Peel](https://theorangepeel.net)
- [The Grey Eagle](https://thegreyeagle.com)
- [Salvage Station](https://salvagestation.com)
- [Asheville Music Hall](https://ashevillemusichall.com)
- [The Mothlight](https://themothlight.com)

---

## Project Vision

> Build Asheville's go-to platform for concert discovery - a beautiful, fast, and comprehensive resource that helps music lovers find shows, discover new bands, and stay connected to the local music scene.

### Success Metrics
- All Asheville venues scraped automatically
- <1 second page load times
- Mobile-friendly design
- Top 3 Google result for "Asheville concerts"

### Out of Scope (for MVP)
- User accounts
- Ticket purchasing
- Mobile app
- Multi-city support

---

## Contact

**Developer Focus**: Backend-focused engineer learning modern frontend development while building production-quality full-stack application.

**Feedback**: Use GitHub issues for bugs and feature requests.
- As we make design decisions and and implementation changes, please update the documentation in @docs/ 
As we complete tasks from the TODO list, please update the @docs/TASKS.md file with completed work, and feel free to refactor and add new tasks as decisions are made
- When writing documentation to .MD files try to be concise and not too verbose. Instructions should be well structured for an AI LLM agent like claude code to pickup in a new session