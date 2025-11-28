# Project Tasks - The Asheville Setlist

> **Purpose**: Granular task tracking for AI agent (Claude Code) implementation. Each task is atomic, testable, and includes acceptance criteria.

---

## Task Status Legend

- `[ ]` - Not started
- `[~]` - In progress
- `[x]` - Completed
- `[!]` - Blocked (see notes)
- `[?]` - Needs clarification

---

## Phase 0: Design Documentation (Current)

> **Goal**: Complete all design docs so implementation can begin without ambiguity.

### 0.1 API Specification
- [x] **TASK-001**: Create `docs/api-spec.md` with full endpoint documentation
  - Acceptance: All endpoints have request/response JSON schemas ✅
  - Acceptance: Error codes and messages defined ✅
  - Acceptance: Pagination structure documented ✅
  - Acceptance: Query parameter validation rules specified ✅

- [x] **TASK-002**: Define API response envelope structure
  - Acceptance: Standard success response format documented ✅
  - Acceptance: Standard error response format documented ✅
  - Acceptance: Pagination metadata structure defined ✅

- [x] **TASK-003**: Document all query parameters for `/api/shows` endpoint
  - Acceptance: Each param has type, validation rules, default value ✅
  - Acceptance: Filter combination behavior documented ✅
  - Acceptance: Sort options and defaults specified ✅

### 0.2 Frontend Specification
- [x] **TASK-004**: Create `docs/frontend-structure.md` with component hierarchy
  - Acceptance: Page-by-page breakdown with components listed ✅
  - Acceptance: SSR vs SSG vs ISR strategy per page ✅
  - Acceptance: State management approach documented ✅

- [x] **TASK-005**: Define URL structure and routing
  - Acceptance: All routes documented with their rendering strategy ✅
  - Acceptance: Dynamic route parameters specified ✅
  - Acceptance: Query string usage for filters documented ✅

- [x] **TASK-006**: Document filter UI behavior
  - Acceptance: Filter state persistence (URL params vs local state) ✅
  - Acceptance: Filter reset behavior defined ✅
  - Acceptance: Mobile filter UX specified ✅

### 0.3 Scraper Specification
- [x] **TASK-007**: Create `docs/scraper-config.md` with venue examples
  - Acceptance: At least 3 real Asheville venue configs documented ✅
  - Acceptance: Selector patterns for each venue ✅
  - Acceptance: Date parsing rules per venue ✅

- [x] **TASK-008**: Document scraper error handling strategy
  - Acceptance: Failure modes identified (network, parsing, rate limit) ✅
  - Acceptance: Retry logic defined ✅
  - Acceptance: Alerting/logging requirements specified ✅

- [x] **TASK-009**: Define band name matching algorithm
  - Acceptance: Fuzzy matching rules documented ✅
  - Acceptance: New band creation criteria defined ✅
  - Acceptance: Manual review queue behavior specified ✅

### 0.4 Seed Data
- [x] **TASK-010**: Create `seeds/genres.sql` with genre taxonomy
  - Acceptance: 15-25 genres covering local music scene ✅
  - Acceptance: Slugs are URL-friendly ✅
  - Acceptance: Descriptions provided for each genre ✅
  - Note: Implemented as `seeds/001_genres.sql`

- [x] **TASK-011**: Create `seeds/venues.sql` with real Asheville venues
  - Acceptance: At least 10 real venues with accurate data ✅
  - Acceptance: Regions assigned (downtown, west, south, etc.) ✅
  - Acceptance: Website URLs verified ✅
  - Note: Implemented as `seeds/002_venues.sql`

- [x] **TASK-012**: Create `seeds/test-data.sql` for development
  - Acceptance: Sample bands with varied genres ✅
  - Acceptance: Sample shows across date ranges ✅
  - Acceptance: Show-band relationships included ✅
  - Note: Implemented as `seeds/004_test_data.sql`

### 0.5 Development Environment
- [x] **TASK-013**: Create `docker-compose.yml` for local development
  - Acceptance: PostgreSQL 16 container configured ✅
  - Acceptance: Volume persistence for database ✅
  - Acceptance: Network configuration for services ✅

- [x] **TASK-014**: Create `.env.example` with all required variables
  - Acceptance: Backend variables documented ✅
  - Acceptance: Frontend variables documented ✅
  - Acceptance: Comments explaining each variable ✅

- [x] **TASK-015**: Create `Makefile` with common commands
  - Acceptance: `make dev` starts local environment ✅
  - Acceptance: `make migrate` runs database migrations ✅
  - Acceptance: `make test` runs all tests ✅
  - Acceptance: `make seed` loads seed data

---

## Phase 1: Foundation

> **Goal**: Set up project structure, local development environment, and core infrastructure.

### 1.1 Repository Setup
- [x] **TASK-101**: Initialize monorepo structure
  - Acceptance: `/backend` and `/frontend` directories created ✅
  - Acceptance: Root `.gitignore` covers both Go and Node ✅
  - Acceptance: Root `README.md` with setup instructions ✅

- [x] **TASK-102**: Initialize Go backend module
  - Acceptance: `go.mod` with module name `github.com/paulsena/asheville-setlist` ✅
  - Acceptance: Directory structure matches planned layout ✅
  - Acceptance: `go build ./...` succeeds ✅

- [x] **TASK-103**: Initialize Next.js frontend
  - Acceptance: Next.js 16 with App Router (latest version) ✅
  - Acceptance: TypeScript strict mode enabled ✅
  - Acceptance: Tailwind CSS configured ✅
  - Acceptance: `npm run build` succeeds ✅

### 1.2 Database Setup
- [x] **TASK-104**: Create initial migration (000001_init_schema)
  - Acceptance: All tables from database-schema.md created ✅
  - Acceptance: Indexes created ✅
  - Acceptance: Constraints added ✅
  - Acceptance: `migrate up` and `migrate down` both work ✅

- [x] **TASK-105**: Create seed migration (000002_seed_data)
  - Acceptance: Genres seeded (31 genres) ✅
  - Acceptance: Venues seeded (33 venues) ✅
  - Acceptance: Migration is idempotent (can run multiple times) ✅

- [x] **TASK-106**: Set up sqlc configuration
  - Acceptance: `sqlc.yaml` configured for PostgreSQL ✅
  - Acceptance: Query files location specified ✅
  - Acceptance: `sqlc generate` produces Go code ✅

### 1.3 Backend Foundation
- [x] **TASK-107**: Create basic Gin server with health endpoint
  - Acceptance: Server starts on configurable port (PORT env var) ✅
  - Acceptance: `GET /health` returns 200 with JSON ✅
  - Acceptance: Graceful shutdown implemented (SIGTERM/SIGINT) ✅

- [x] **TASK-108**: Set up database connection pool
  - Acceptance: Connection string from environment variable (DATABASE_URL) ✅
  - Acceptance: Pool settings configured (25 max, 5 min, 1hr lifetime, 30min idle) ✅
  - Acceptance: Connection test on startup (ping + log message) ✅

- [x] **TASK-109**: Create base middleware stack
  - Acceptance: CORS middleware configured (allows all origins for dev) ✅
  - Acceptance: Request logging middleware (structured JSON) ✅
  - Acceptance: Panic recovery middleware (returns 500) ✅

- [x] **TASK-110**: Set up configuration management
  - Acceptance: Environment variable loading ✅
  - Acceptance: Config struct with validation ✅
  - Acceptance: Defaults for development ✅

### 1.4 Frontend Foundation
- [x] **TASK-111**: Set up shadcn/ui
  - Acceptance: `components.json` configured ✅
  - Acceptance: Base components installed (Button, Card, Input) ✅
  - Acceptance: Theme configuration (colors, fonts) ✅
  - Note: Installed with neutral color scheme and New York style

- [x] **TASK-112**: Create API client utility
  - Acceptance: Base fetch wrapper with error handling ✅
  - Acceptance: Type-safe response parsing ✅
  - Acceptance: Environment-based API URL ✅
  - Note: Created `lib/api.ts` with APIResponse/APIError types

- [x] **TASK-113**: Set up TanStack Query
  - Acceptance: QueryClient configured ✅
  - Acceptance: Provider wrapped in layout ✅
  - Acceptance: DevTools available in development ✅
  - Note: Created `lib/query-provider.tsx` with 5min stale time

- [x] **TASK-114**: Create basic layout components
  - Acceptance: Root layout with metadata ✅
  - Acceptance: Header component (placeholder) ✅
  - Acceptance: Footer component (placeholder) ✅
  - Note: Added Header/Footer with navigation, SEO metadata configured
  - Note: Installed v0.dev CLI for AI-generated UI components

---

## Phase 2: Core API

> **Goal**: Implement all REST API endpoints for shows, venues, bands, and search.

### 2.1 Database Queries (sqlc)
- [x] **TASK-201**: Write sqlc queries for shows
  - Acceptance: ListShows with filters (date, venue, genre, region, price) ✅
  - Acceptance: GetShowByID with bands and venue ✅
  - Acceptance: GetUpcomingShows (next 30 days) ✅
  - Acceptance: CreateShow for band submissions ✅

- [x] **TASK-202**: Write sqlc queries for venues
  - Acceptance: ListVenues with region filter ✅
  - Acceptance: GetVenueBySlug with upcoming shows ✅
  - Acceptance: GetVenueByID ✅

- [x] **TASK-203**: Write sqlc queries for bands
  - Acceptance: ListBands with pagination ✅
  - Acceptance: GetBandBySlug with genres and shows ✅
  - Acceptance: GetSimilarBands (shared genres) ✅
  - Acceptance: SearchBands (full-text) ✅

- [x] **TASK-204**: Write sqlc queries for genres
  - Acceptance: ListGenres ✅
  - Acceptance: GetGenreBySlug with bands ✅

- [x] **TASK-205**: Write sqlc query for global search
  - Acceptance: Search shows, bands, venues ✅
  - Acceptance: Results unified with type field ✅
  - Acceptance: Relevance ranking ✅

### 2.2 API Handlers
- [x] **TASK-206**: Implement `GET /api/shows` handler
  - Acceptance: All query parameters working ✅
  - Acceptance: Pagination implemented ✅
  - Acceptance: Response matches API spec ✅

- [x] **TASK-207**: Implement `GET /api/shows/:id` handler
  - Acceptance: Returns show with bands and venue ✅
  - Acceptance: 404 for non-existent show ✅
  - Acceptance: Response matches API spec ✅

- [x] **TASK-208**: Implement `GET /api/venues` handler
  - Acceptance: Region filter working ✅
  - Acceptance: Returns venue list with basic info ✅

- [x] **TASK-209**: Implement `GET /api/venues/:slug` handler
  - Acceptance: Returns venue details ✅
  - Acceptance: Includes upcoming shows ✅
  - Acceptance: 404 for non-existent venue ✅

- [x] **TASK-210**: Implement `GET /api/bands` handler
  - Acceptance: Pagination working ✅
  - Acceptance: Genre filter working ✅

- [x] **TASK-211**: Implement `GET /api/bands/:slug` handler
  - Acceptance: Returns band details with genres ✅
  - Acceptance: Includes upcoming shows ✅
  - Acceptance: 404 for non-existent band ✅

- [x] **TASK-212**: Implement `GET /api/bands/:slug/similar` handler
  - Acceptance: Returns bands with shared genres ✅
  - Acceptance: Excludes the queried band ✅
  - Acceptance: Ordered by genre overlap count ✅

- [x] **TASK-213**: Implement `GET /api/genres` handler
  - Acceptance: Returns all genres ✅
  - Acceptance: Includes show/band counts (optional) ✅

- [x] **TASK-214**: Implement `GET /api/search` handler
  - Acceptance: Query parameter `q` required ✅
  - Acceptance: Returns mixed results (shows, bands, venues) ✅
  - Acceptance: Type field in each result ✅

- [x] **TASK-215**: Implement `POST /api/shows` handler (band submission)
  - Acceptance: Request validation (required fields) ✅
  - Acceptance: Creates show with 'scheduled' status ✅ (changed from 'pending' per DB constraint)
  - Acceptance: Links bands (creates if needed) ✅
  - Acceptance: Returns created show ID ✅

### 2.3 API Testing
- [x] **TASK-216**: Write integration tests for shows endpoints
  - Acceptance: Test list with various filters ✅
  - Acceptance: Test pagination ✅
  - Acceptance: Test 404 scenarios ✅
  - Note: Tests in `backend/internal/handlers/shows_test.go`

- [x] **TASK-217**: Write integration tests for venues endpoints
  - Acceptance: Test list and detail endpoints ✅
  - Acceptance: Test region filtering ✅
  - Note: Tests in `backend/internal/handlers/venues_test.go`

- [x] **TASK-218**: Write integration tests for bands endpoints
  - Acceptance: Test list, detail, similar endpoints ✅
  - Acceptance: Test genre filtering ✅
  - Note: Tests in `backend/internal/handlers/bands_test.go`

- [x] **TASK-219**: Write integration tests for search endpoint
  - Acceptance: Test various search queries ✅
  - Acceptance: Test empty results ✅
  - Note: Tests in `backend/internal/handlers/search_test.go`
  - Note: Created `backend/internal/testutil/testutil.go` for test utilities

---

## Phase 3: Scraper Service

> **Goal**: Implement automated venue scraping with concurrent execution.

### 3.1 Scraper Core
- [ ] **TASK-301**: Create scraper configuration loader
  - Acceptance: Reads venue configs from database
  - Acceptance: Validates selector structure
  - Acceptance: Filters active scrapers only

- [ ] **TASK-302**: Implement Colly-based static scraper
  - Acceptance: Fetches HTML from venue URL
  - Acceptance: Extracts shows using CSS selectors
  - Acceptance: Handles pagination if needed

- [ ] **TASK-303**: Implement chromedp-based JS scraper
  - Acceptance: Renders JavaScript content
  - Acceptance: Waits for dynamic elements
  - Acceptance: Falls back gracefully

- [ ] **TASK-304**: Create date parsing utility
  - Acceptance: Handles multiple date formats
  - Acceptance: Handles relative dates ("Tomorrow", "This Saturday")
  - Acceptance: Timezone-aware (Eastern)

- [ ] **TASK-305**: Create band name extraction utility
  - Acceptance: Splits headliner/opener patterns
  - Acceptance: Handles "with", "featuring", "&" patterns
  - Acceptance: Cleans up formatting artifacts

### 3.2 Data Processing
- [ ] **TASK-306**: Implement band matching algorithm
  - Acceptance: Exact match by name/slug
  - Acceptance: Fuzzy match with threshold
  - Acceptance: Creates new band if no match

- [ ] **TASK-307**: Implement show deduplication
  - Acceptance: Match by venue + date + headliner
  - Acceptance: Update if changed
  - Acceptance: Skip if identical

- [ ] **TASK-308**: Implement upsert transaction
  - Acceptance: Atomic insert/update
  - Acceptance: Rollback on error
  - Acceptance: Returns created/updated counts

### 3.3 Scraper Orchestration
- [ ] **TASK-309**: Implement concurrent scraping
  - Acceptance: Goroutines per venue
  - Acceptance: Rate limiting per domain
  - Acceptance: Error isolation (one failure doesn't stop others)

- [ ] **TASK-310**: Implement scraper CLI
  - Acceptance: `./scraper run` executes all scrapers
  - Acceptance: `./scraper run --venue=orange-peel` single venue
  - Acceptance: `./scraper dry-run` shows what would be scraped

- [ ] **TASK-311**: Implement scraper logging and metrics
  - Acceptance: Structured JSON logs
  - Acceptance: Per-venue success/failure counts
  - Acceptance: Timing metrics

### 3.4 Venue Configurations
- [ ] **TASK-312**: Create Orange Peel scraper config
  - Acceptance: Selectors tested and working
  - Acceptance: Date format documented
  - Acceptance: Sample output verified

- [ ] **TASK-313**: Create Grey Eagle scraper config
  - Acceptance: Selectors tested and working
  - Acceptance: Date format documented
  - Acceptance: Sample output verified

- [ ] **TASK-314**: Create Salvage Station scraper config
  - Acceptance: Selectors tested and working
  - Acceptance: Date format documented
  - Acceptance: Sample output verified

---

## Phase 4: Frontend UI

> **Goal**: Build all user-facing pages with responsive design.

### 4.1 Homepage
- [ ] **TASK-401**: Create homepage layout
  - Acceptance: Hero section with site intro
  - Acceptance: Featured/upcoming shows section
  - Acceptance: Quick genre links
  - Acceptance: Responsive design

- [ ] **TASK-402**: Create ShowCard component
  - Acceptance: Shows date, venue, bands, price
  - Acceptance: Links to show detail page
  - Acceptance: Hover/focus states

- [ ] **TASK-403**: Implement homepage data fetching
  - Acceptance: SSR with revalidation
  - Acceptance: Shows next 10 upcoming shows
  - Acceptance: Loading state handled

### 4.2 Shows Listing
- [ ] **TASK-404**: Create shows listing page
  - Acceptance: Grid/list of ShowCards
  - Acceptance: Pagination or infinite scroll
  - Acceptance: SEO metadata

- [ ] **TASK-405**: Create filter sidebar/panel
  - Acceptance: Date range picker
  - Acceptance: Venue multiselect
  - Acceptance: Genre multiselect
  - Acceptance: Price range slider
  - Acceptance: Region filter

- [ ] **TASK-406**: Implement filter state management
  - Acceptance: Filters sync to URL params
  - Acceptance: Browser back/forward works
  - Acceptance: Share filtered URL works

- [ ] **TASK-407**: Create mobile filter drawer
  - Acceptance: Full-screen on mobile
  - Acceptance: Apply/clear buttons
  - Acceptance: Smooth animation

### 4.3 Show Detail
- [ ] **TASK-408**: Create show detail page
  - Acceptance: Full show info displayed
  - Acceptance: Band list with links
  - Acceptance: Venue info with map link
  - Acceptance: Ticket link button

- [ ] **TASK-409**: Create band preview component
  - Acceptance: Spotify embed if available
  - Acceptance: Fallback to Bandcamp/website
  - Acceptance: Graceful degradation

- [ ] **TASK-410**: Implement show detail SEO
  - Acceptance: Dynamic title and description
  - Acceptance: Open Graph tags
  - Acceptance: Structured data (JSON-LD)

### 4.4 Venues
- [ ] **TASK-411**: Create venues listing page
  - Acceptance: Venue cards with images
  - Acceptance: Region filter tabs
  - Acceptance: Show count per venue

- [ ] **TASK-412**: Create venue detail page
  - Acceptance: Venue info and image
  - Acceptance: Upcoming shows list
  - Acceptance: Map embed or link

### 4.5 Bands
- [ ] **TASK-413**: Create band detail page
  - Acceptance: Band info and image
  - Acceptance: Genre tags
  - Acceptance: Upcoming shows
  - Acceptance: Similar bands section

- [ ] **TASK-414**: Create similar bands component
  - Acceptance: Horizontal scroll or grid
  - Acceptance: Links to band pages
  - Acceptance: Genre overlap indicator

### 4.6 Search
- [ ] **TASK-415**: Create search page
  - Acceptance: Search input with debounce
  - Acceptance: Results grouped by type
  - Acceptance: Empty state message

- [ ] **TASK-416**: Create global search component (header)
  - Acceptance: Search icon trigger
  - Acceptance: Modal or dropdown results
  - Acceptance: Keyboard navigation

### 4.7 Band Submission
- [ ] **TASK-417**: Create show submission form
  - Acceptance: Band name input
  - Acceptance: Venue selector
  - Acceptance: Date/time picker
  - Acceptance: Price and ticket URL fields

- [ ] **TASK-418**: Implement form validation
  - Acceptance: Zod schema validation
  - Acceptance: Inline error messages
  - Acceptance: Submit button disabled when invalid

- [ ] **TASK-419**: Implement form submission
  - Acceptance: POST to API
  - Acceptance: Success confirmation
  - Acceptance: Error handling with retry

---

## Phase 5: Deployment

> **Goal**: Deploy to production with CI/CD.

### 5.1 Containerization
- [ ] **TASK-501**: Create API Dockerfile
  - Acceptance: Multi-stage build
  - Acceptance: Small final image (<50MB)
  - Acceptance: Non-root user
  - Acceptance: Health check endpoint

- [ ] **TASK-502**: Create Scraper Dockerfile
  - Acceptance: Includes chromium for JS scraping
  - Acceptance: Multi-stage build
  - Acceptance: Handles long-running process

### 5.2 Google Cloud Setup
- [ ] **TASK-503**: Create GCP project and enable APIs
  - Acceptance: Cloud Run API enabled
  - Acceptance: Cloud Build API enabled
  - Acceptance: Cloud Scheduler API enabled
  - Acceptance: Secret Manager API enabled

- [ ] **TASK-504**: Set up Secret Manager
  - Acceptance: DATABASE_URL stored
  - Acceptance: IAM permissions configured

- [ ] **TASK-505**: Deploy API to Cloud Run
  - Acceptance: Service running
  - Acceptance: Health endpoint accessible
  - Acceptance: Secrets mounted

- [ ] **TASK-506**: Deploy Scraper as Cloud Run Job
  - Acceptance: Job created
  - Acceptance: Manual execution works
  - Acceptance: Logs visible

- [ ] **TASK-507**: Set up Cloud Scheduler
  - Acceptance: Cron job triggers scraper
  - Acceptance: Schedule is every 6 hours
  - Acceptance: Service account configured

### 5.3 Vercel Setup
- [ ] **TASK-508**: Connect GitHub to Vercel
  - Acceptance: Auto-deploy on push to main
  - Acceptance: Preview deploys on PRs

- [ ] **TASK-509**: Configure environment variables
  - Acceptance: NEXT_PUBLIC_API_URL set
  - Acceptance: Production and preview environments configured

- [ ] **TASK-510**: Set up custom domain (optional)
  - Acceptance: DNS configured
  - Acceptance: SSL working

### 5.4 Database Setup
- [ ] **TASK-511**: Create Neon project
  - Acceptance: Database created
  - Acceptance: Connection string obtained
  - Acceptance: Pooler connection available

- [ ] **TASK-512**: Run production migrations
  - Acceptance: All migrations applied
  - Acceptance: Seed data loaded

### 5.5 CI/CD
- [ ] **TASK-513**: Create Cloud Build trigger for API
  - Acceptance: Triggers on push to main
  - Acceptance: Builds and deploys automatically
  - Acceptance: Notifications on failure

- [ ] **TASK-514**: Create GitHub Actions for tests
  - Acceptance: Runs on PR
  - Acceptance: Go tests pass
  - Acceptance: TypeScript builds

---

## Phase 6: Polish & Launch

> **Goal**: Final testing, optimizations, and public launch.

### 6.1 Testing
- [ ] **TASK-601**: End-to-end testing with Playwright
  - Acceptance: Homepage loads
  - Acceptance: Filters work
  - Acceptance: Search works
  - Acceptance: Form submission works

- [ ] **TASK-602**: Performance testing
  - Acceptance: Lighthouse score >90 (mobile)
  - Acceptance: API response <200ms (p95)
  - Acceptance: Core Web Vitals pass

- [ ] **TASK-603**: Accessibility audit
  - Acceptance: WCAG 2.1 AA compliance
  - Acceptance: Keyboard navigation works
  - Acceptance: Screen reader tested

### 6.2 SEO
- [ ] **TASK-604**: Create sitemap.xml
  - Acceptance: All public pages included
  - Acceptance: Auto-generates on build

- [ ] **TASK-605**: Create robots.txt
  - Acceptance: Allows search engines
  - Acceptance: Blocks admin routes (if any)

- [ ] **TASK-606**: Submit to Google Search Console
  - Acceptance: Site verified
  - Acceptance: Sitemap submitted
  - Acceptance: No critical errors

### 6.3 Monitoring
- [ ] **TASK-607**: Set up Cloud Monitoring dashboard
  - Acceptance: API request count
  - Acceptance: Error rate
  - Acceptance: Latency percentiles

- [ ] **TASK-608**: Set up error alerting
  - Acceptance: Alerts on high error rate
  - Acceptance: Alerts on scraper failures
  - Acceptance: Email notifications configured

### 6.4 Documentation
- [ ] **TASK-609**: Write user-facing FAQ
  - Acceptance: Common questions answered
  - Acceptance: Accessible from footer

- [ ] **TASK-610**: Write band submission guide
  - Acceptance: Step-by-step instructions
  - Acceptance: What to expect after submission

---

## Backlog (Future Phases)

> **Not scheduled for MVP. Track ideas here.**

### User Features (Phase 2)
- [ ] User authentication (Auth0/Clerk)
- [ ] User favorites
- [ ] Follow bands/venues
- [ ] Email notifications

### Content (Phase 2)
- [ ] Articles/blog CMS
- [ ] Admin dashboard

### Advanced Features (Phase 3)
- [ ] Mobile app
- [ ] Real-time updates
- [ ] Ticketing integration
- [ ] ML recommendations

---

## Notes for AI Agent (Claude Code)

### How to Use This File

1. **Pick a task**: Start with lowest numbered incomplete task in current phase
2. **Read acceptance criteria**: All must pass for task to be complete
3. **Check dependencies**: Some tasks depend on others (implicit by numbering)
4. **Mark progress**: Change `[ ]` to `[~]` when starting, `[x]` when complete
5. **Add notes**: If blocked or needs clarification, add note below task

### Task Sizing

- Each task should take 15-60 minutes
- If a task feels too large, it should be split
- If a task feels too small, it can be combined with related tasks

### Definition of Done

A task is complete when:
1. All acceptance criteria pass
2. Code compiles/builds without errors
3. Tests pass (if applicable)
4. Code is committed with descriptive message

### Commit Message Format

```
[TASK-XXX] Short description

- Acceptance criteria 1: done
- Acceptance criteria 2: done
- Notes if any
```

### When Stuck

If a task is blocked:
1. Mark with `[!]`
2. Add note explaining blocker
3. Move to next available task
4. Flag for human review

---

## Change Log

| Date | Tasks | Notes |
|------|-------|-------|
| 2025-01-XX | Created | Initial task breakdown |

