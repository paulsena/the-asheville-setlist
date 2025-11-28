# The Asheville Setlist

> Your guide to live music in Asheville, NC - discover concerts, explore venues, and find your next show.

## Overview

The Asheville Setlist is a concert discovery platform that aggregates show information from local venues, making it easy to find live music in Asheville's vibrant music scene.

**Features:**
- Automated venue scraping for up-to-date show listings
- Advanced filtering by date, venue, genre, and price
- Band discovery with genre-based recommendations
- Mobile-friendly responsive design
- SEO-optimized for discoverability

## Tech Stack

### Backend
- **Language:** Go 1.22+
- **Framework:** Gin web framework
- **Database:** PostgreSQL 16 with JSONB support
- **Query Layer:** sqlc for type-safe SQL
- **Migrations:** golang-migrate
- **Scraping:** Colly + chromedp

### Frontend
- **Framework:** Next.js 15 (App Router)
- **Language:** TypeScript (strict mode)
- **Styling:** Tailwind CSS
- **Components:** shadcn/ui
- **State:** TanStack Query
- **Forms:** React Hook Form + Zod

### Infrastructure
- **Frontend:** Vercel
- **Backend:** Google Cloud Run
- **Database:** Neon (PostgreSQL)
- **Scheduler:** Cloud Scheduler
- **CI/CD:** Cloud Build + Vercel

## Project Structure

```
asheville-setlist/
├── backend/           # Go API server and scraper
├── frontend/          # Next.js web application
├── docs/              # Project documentation
├── seeds/             # Database seed data
├── docker-compose.yml # Local development environment
├── Makefile           # Common development commands
└── .env.example       # Environment variables template
```

## Quick Start

### Prerequisites

- **Go:** 1.22 or higher ([install](https://go.dev/doc/install))
- **Node.js:** 20 or higher ([install](https://nodejs.org/))
- **Docker:** For local PostgreSQL ([install](https://docs.docker.com/get-docker/))
- **Make:** For convenience commands (usually pre-installed on macOS/Linux)

### Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/[username]/asheville-setlist.git
   cd asheville-setlist
   ```

2. **Copy environment file**
   ```bash
   cp .env.example .env
   ```
   Edit `.env` with your configuration (local development works with defaults)

3. **Start local services**
   ```bash
   make dev
   ```
   This starts PostgreSQL via Docker Compose

4. **Run database migrations**
   ```bash
   make migrate
   ```

5. **Load seed data**
   ```bash
   make seed
   ```

6. **Start the API server** (in a new terminal)
   ```bash
   make api
   ```

7. **Start the frontend** (in another terminal)
   ```bash
   make frontend
   ```

8. **Access the application**
   - **Frontend:** http://localhost:3000
   - **API:** http://localhost:8080
   - **API Health:** http://localhost:8080/health

## Development

### Common Commands

```bash
# Start all services
make dev              # PostgreSQL via Docker Compose

# Run specific services
make api              # Start API server (port 8080)
make frontend         # Start Next.js dev server (port 3000)
make scraper          # Run scraper once

# Database management
make migrate          # Run pending migrations
make migrate-down     # Rollback last migration
make seed             # Load seed data
make db-reset         # Drop, recreate, migrate, and seed

# Code generation
make sqlc             # Generate Go code from SQL queries

# Testing
make test             # Run all tests
make test-api         # Backend tests only
make test-frontend    # Frontend tests only

# Linting
make lint             # Run all linters
```

### Backend Development

```bash
cd backend

# Install dependencies
go mod download

# Run API server
go run cmd/api/main.go

# Run scraper
go run cmd/scraper/main.go

# Run tests
go test ./...

# Generate sqlc code (after modifying queries)
sqlc generate
```

### Frontend Development

```bash
cd frontend

# Install dependencies
npm install

# Run dev server
npm run dev

# Build for production
npm run build

# Run production build locally
npm start

# Run linter
npm run lint

# Run type checking
npm run type-check
```

## Database

The application uses PostgreSQL 16 with the following schema:

- **venues** - Concert venues with location and capacity info
- **bands** - Artists with genre associations
- **shows** - Concert events linking venues and bands
- **genres** - Music genre taxonomy
- **show_bands** - Many-to-many relationship (shows ↔ bands)
- **band_genres** - Many-to-many relationship (bands ↔ genres)
- **venue_scrapers** - Scraping configuration per venue

See `docs/database-schema.md` for detailed schema documentation.

## API Documentation

RESTful API endpoints:

- `GET /api/shows` - List shows with filtering
- `GET /api/shows/:id` - Show details
- `GET /api/venues` - List venues
- `GET /api/venues/:slug` - Venue details with shows
- `GET /api/bands` - List bands
- `GET /api/bands/:slug` - Band details
- `GET /api/bands/:slug/similar` - Similar bands by genre
- `GET /api/genres` - List all genres
- `GET /api/search?q=` - Global search

See `docs/api-spec.md` for complete API documentation.

## Deployment

### Production Deployment

The application is deployed to:
- **Frontend:** Vercel (automatic deploys from `main` branch)
- **Backend:** Google Cloud Run
- **Database:** Neon PostgreSQL
- **Scraper:** Cloud Run Job (scheduled via Cloud Scheduler)

See `docs/deployment.md` for deployment instructions.

### Environment Variables

**Backend (.env):**
```bash
DATABASE_URL=postgres://user:pass@host:5432/db?sslmode=require
PORT=8080
GIN_MODE=release
LOG_LEVEL=info
```

**Frontend (.env.local):**
```bash
NEXT_PUBLIC_API_URL=https://api.ashevillesetlist.com
```

## Project Status

See `docs/TASKS.md` for current development status and task tracking.

**Current Phase:** Phase 1 - Foundation
**Completed:** Phase 0 - Design Documentation

## Contributing

This is currently a personal project. Feature requests and bug reports are welcome via GitHub issues.

## Documentation

- `docs/architecture.md` - System architecture and design decisions
- `docs/database-schema.md` - Database schema reference
- `docs/tech-stack.md` - Technology choices and rationale
- `docs/api-spec.md` - API endpoint documentation
- `docs/frontend-structure.md` - Frontend component structure
- `docs/scraper-config.md` - Venue scraper configurations
- `docs/deployment.md` - Deployment guide
- `docs/TASKS.md` - Development task tracking

## License

MIT License - See LICENSE file for details

## Contact

Questions or feedback? Open an issue on GitHub.

---

**Built with ❤️ for Asheville's music community**
