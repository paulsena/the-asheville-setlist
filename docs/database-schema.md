# Database Schema

## Overview

PostgreSQL 16 relational database schema designed for concert discovery with rich filtering, similar artist recommendations, and flexible scraped data storage.

---

## Entity Relationship Diagram

```
┌──────────────┐         ┌──────────────┐         ┌──────────────┐
│   venues     │         │    shows     │         │    bands     │
├──────────────┤         ├──────────────┤         ├──────────────┤
│ id (PK)      │────────<│ venue_id(FK) │>────────│ id (PK)      │
│ name         │         │ id (PK)      │         │ name         │
│ address      │         │ title        │         │ bio          │
│ region       │         │ image_url    │         │ image_url    │
│ capacity     │         │ date         │         │ website      │
│ image_url    │         │ doors_time   │         │ spotify_url  │
│ website      │         │ age_restrict │         │ instagram    │
│ metadata     │         │ price_min    │         │ metadata     │
│ created_at   │         │ price_max    │         │ created_at   │
│ updated_at   │         │ ticket_url   │         │ updated_at   │
└──────────────┘         │ status       │         └──────────────┘
                         │ source       │                │
                         │ scraped_data │                │
                         │ created_at   │                │
                         │ updated_at   │                │
                         └──────────────┘                │
                                │                        │
                                │                        │
                         ┌──────▼────────┐        ┌──────▼───────┐
                         │  show_bands   │        │ band_genres  │
                         ├───────────────┤        ├──────────────┤
                         │ show_id (FK)  │        │ band_id (FK) │
                         │ band_id (FK)  │        │ genre_id(FK) │
                         │ is_headliner  │        │ created_at   │
                         │ created_at    │        └──────────────┘
                         └───────────────┘               │
                                                         │
                                                  ┌──────▼──────┐
                                                  │   genres    │
                                                  ├─────────────┤
                                                  │ id (PK)     │
                                                  │ name        │
                                                  │ slug        │
                                                  │ description │
                                                  │ created_at  │
                                                  └─────────────┘

┌──────────────┐         ┌──────────────┐
│   articles   │         │    users     │
├──────────────┤         ├──────────────┤
│ id (PK)      │         │ id (PK)      │  (Future phase)
│ title        │         │ email        │
│ slug         │         │ name         │
│ content      │         │ password_hash│
│ author       │         │ role         │
│ published_at │         │ created_at   │
│ created_at   │         │ updated_at   │
│ updated_at   │         └──────────────┘
└──────────────┘
```

---

## Tables

### 1. venues

Stores venue information (clubs, bars, concert halls).

```sql
CREATE TABLE venues (
    id SERIAL PRIMARY KEY,

    -- Basic info
    name TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL, -- URL-friendly: "the-orange-peel"

    -- Location
    address TEXT,
    city TEXT DEFAULT 'Asheville',
    state TEXT DEFAULT 'NC',
    zip_code TEXT,
    region TEXT, -- "downtown", "west asheville", "south asheville", etc.
    latitude DECIMAL(9,6),
    longitude DECIMAL(9,6),

    -- Venue details
    capacity INTEGER,
    website TEXT,
    phone TEXT,
    image_url TEXT,

    -- Flexible metadata (JSONB for varying details)
    metadata JSONB,
    -- Example: {"parking": "street", "food": true, "bar": true, "outdoor": false}

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_venues_region ON venues(region);
CREATE INDEX idx_venues_slug ON venues(slug);
CREATE INDEX idx_venues_metadata ON venues USING GIN(metadata);

-- Full-text search
CREATE INDEX idx_venues_search ON venues USING GIN(to_tsvector('english', name));
```

**Example Data**:
```sql
INSERT INTO venues (name, slug, address, region, capacity, website, metadata) VALUES
('The Orange Peel', 'the-orange-peel', '101 Biltmore Ave', 'downtown', 1050, 'https://theorangepeel.net',
 '{"parking": "street", "food": false, "bar": true, "outdoor": false}'::jsonb),
('The Grey Eagle', 'the-grey-eagle', '185 Clingman Ave', 'south asheville', 700, 'https://thegreyeagle.com',
 '{"parking": "lot", "food": true, "bar": true, "outdoor": true}'::jsonb);
```

---

### 2. bands

Stores band/artist information.

```sql
CREATE TABLE bands (
    id SERIAL PRIMARY KEY,

    -- Basic info
    name TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL, -- "moon-taxi"

    -- Details
    bio TEXT,
    hometown TEXT,
    image_url TEXT,

    -- Links
    website TEXT,
    spotify_url TEXT,
    instagram TEXT,
    facebook TEXT,
    bandcamp_url TEXT,

    -- Flexible metadata
    metadata JSONB,
    -- Example: {"verified": true, "local": true, "tour_dates_url": "..."}

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE UNIQUE INDEX idx_bands_slug ON bands(slug);
CREATE INDEX idx_bands_name ON bands(name);
CREATE INDEX idx_bands_search ON bands USING GIN(to_tsvector('english', name || ' ' || COALESCE(bio, '')));
```

**Example Data**:
```sql
INSERT INTO bands (name, slug, bio, hometown, spotify_url) VALUES
('Moon Taxi', 'moon-taxi', 'Nashville-based indie rock band', 'Nashville, TN',
 'https://open.spotify.com/artist/...'),
('The String Cheese Incident', 'string-cheese-incident', 'Jam band from Colorado', 'Boulder, CO',
 'https://open.spotify.com/artist/...');
```

---

### 3. genres

Music genres for categorization and similar artist matching.

```sql
CREATE TABLE genres (
    id SERIAL PRIMARY KEY,

    name TEXT NOT NULL UNIQUE,
    slug TEXT NOT NULL UNIQUE, -- "indie-rock"
    description TEXT,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Index
CREATE INDEX idx_genres_slug ON genres(slug);
```

**Example Data**:
```sql
INSERT INTO genres (name, slug, description) VALUES
('Rock', 'rock', 'Rock music'),
('Indie', 'indie', 'Independent music'),
('Electronic', 'electronic', 'Electronic music'),
('Hip Hop', 'hip-hop', 'Hip hop and rap'),
('Jazz', 'jazz', 'Jazz music'),
('Bluegrass', 'bluegrass', 'Bluegrass and folk'),
('Jam Band', 'jam-band', 'Improvisational jam bands');
```

---

### 4. band_genres

Many-to-many relationship between bands and genres.

```sql
CREATE TABLE band_genres (
    band_id INTEGER NOT NULL REFERENCES bands(id) ON DELETE CASCADE,
    genre_id INTEGER NOT NULL REFERENCES genres(id) ON DELETE CASCADE,

    -- Timestamp
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    PRIMARY KEY (band_id, genre_id)
);

-- Indexes
CREATE INDEX idx_band_genres_band ON band_genres(band_id);
CREATE INDEX idx_band_genres_genre ON band_genres(genre_id);
```

**Example Data**:
```sql
-- Moon Taxi is indie + rock
INSERT INTO band_genres (band_id, genre_id) VALUES
(1, 1), -- rock
(1, 2); -- indie
```

---

### 5. shows

Concert/show events.

```sql
CREATE TABLE shows (
    id SERIAL PRIMARY KEY,

    -- Relationships
    venue_id INTEGER NOT NULL REFERENCES venues(id) ON DELETE CASCADE,

    -- Basic info
    title TEXT, -- "Moon Taxi with special guests" (nullable, can derive from bands)
    description TEXT,
    image_url TEXT, -- Event poster/flyer image (optional, UI falls back to band/venue images)

    -- Date/time
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    doors_time TIME, -- Door time (e.g., 19:00)
    show_time TIME,  -- Show start time (e.g., 20:00)

    -- Ticketing
    price_min NUMERIC(10,2), -- Minimum price in USD (NULL = free/TBD)
    price_max NUMERIC(10,2), -- Maximum price in USD (NULL if single price)
    ticket_url TEXT,
    age_restriction TEXT, -- "21+", "All Ages", "18+"

    -- Status & Source
    status TEXT DEFAULT 'scheduled',
    -- Values: 'scheduled', 'cancelled', 'postponed', 'completed'
    source TEXT DEFAULT 'scraped',
    -- Values: 'scraped', 'band_submitted', 'manual'

    -- Scraped data (varying structure per venue)
    scraped_data JSONB,
    -- Example: {"source": "venue_website", "raw_html": "...", "scraper_version": "1.0"}

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_shows_venue ON shows(venue_id);
CREATE INDEX idx_shows_date ON shows(date);
CREATE INDEX idx_shows_status ON shows(status);
CREATE INDEX idx_shows_source ON shows(source);
CREATE INDEX idx_shows_date_venue ON shows(date, venue_id); -- Composite for common queries
CREATE INDEX idx_shows_scraped_data ON shows USING GIN(scraped_data);

-- Partial index for upcoming shows (most common query)
CREATE INDEX idx_shows_upcoming ON shows(date) WHERE date >= NOW() AND status = 'scheduled';
```

**Example Data**:
```sql
INSERT INTO shows (venue_id, title, image_url, date, doors_time, show_time, price_min, price_max, ticket_url, status, source) VALUES
(1, 'Moon Taxi', 'https://example.com/posters/moon-taxi-asheville.jpg', '2025-11-15 20:00:00-05', '19:00', '20:00', 25.00, 35.00,
 'https://theorangepeel.net/tickets/moon-taxi', 'scheduled', 'scraped');
```

---

### 6. show_bands

Many-to-many relationship between shows and bands (lineup).

```sql
CREATE TABLE show_bands (
    show_id INTEGER NOT NULL REFERENCES shows(id) ON DELETE CASCADE,
    band_id INTEGER NOT NULL REFERENCES bands(id) ON DELETE CASCADE,

    -- Lineup details
    is_headliner BOOLEAN DEFAULT FALSE,
    performance_order INTEGER, -- 1 = opener, 2 = support, 3 = headliner

    -- Timestamp
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    PRIMARY KEY (show_id, band_id)
);

-- Indexes
CREATE INDEX idx_show_bands_show ON show_bands(show_id);
CREATE INDEX idx_show_bands_band ON show_bands(band_id);
```

**Example Data**:
```sql
-- Moon Taxi headlining, opener TBD
INSERT INTO show_bands (show_id, band_id, is_headliner, performance_order) VALUES
(1, 1, TRUE, 2); -- Moon Taxi is headliner
```

---

### 7. articles

Blog posts and editorial content (minor feature).

```sql
CREATE TABLE articles (
    id SERIAL PRIMARY KEY,

    -- Content
    title TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL, -- "asheville-music-scene-2025"
    content TEXT NOT NULL, -- Markdown or HTML
    excerpt TEXT, -- Short summary

    -- Metadata
    author TEXT, -- Future: FK to users table
    cover_image_url TEXT,

    -- Publishing
    published_at TIMESTAMP WITH TIME ZONE,
    is_published BOOLEAN DEFAULT FALSE,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_articles_slug ON articles(slug);
CREATE INDEX idx_articles_published ON articles(published_at) WHERE is_published = TRUE;
CREATE INDEX idx_articles_search ON articles USING GIN(to_tsvector('english', title || ' ' || excerpt));
```

---

### 8. users (Future Phase)

User accounts for favorites, following, etc.

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,

    -- Authentication
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,

    -- Profile
    name TEXT,
    avatar_url TEXT,

    -- Authorization
    role TEXT DEFAULT 'user', -- 'user', 'admin', 'band_manager'

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_login_at TIMESTAMP WITH TIME ZONE
);

-- Indexes
CREATE UNIQUE INDEX idx_users_email ON users(email);
```

---

### 9. venue_scrapers (Configuration Table)

Stores scraping configuration for each venue.

```sql
CREATE TABLE venue_scrapers (
    id SERIAL PRIMARY KEY,
    venue_id INTEGER NOT NULL REFERENCES venues(id) ON DELETE CASCADE,

    -- Scraper config
    url TEXT NOT NULL, -- URL to scrape
    scraper_type TEXT NOT NULL, -- 'static' or 'javascript'

    -- CSS selectors (JSONB for flexibility)
    selectors JSONB NOT NULL,
    -- Example: {"container": ".event", "title": ".event-title", "date": ".event-date"}

    -- Parsing rules
    date_format TEXT, -- "2006-01-02" (Go time format)

    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    last_scraped_at TIMESTAMP WITH TIME ZONE,
    last_success_at TIMESTAMP WITH TIME ZONE,
    error_count INTEGER DEFAULT 0,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_venue_scrapers_venue ON venue_scrapers(venue_id);
CREATE INDEX idx_venue_scrapers_active ON venue_scrapers(is_active) WHERE is_active = TRUE;
```

**Example Data**:
```sql
INSERT INTO venue_scrapers (venue_id, url, scraper_type, selectors, date_format) VALUES
(1, 'https://theorangepeel.net/events', 'static',
 '{"container": ".event-item", "title": ".event-title", "date": ".event-date", "bands": ".lineup a"}'::jsonb,
 'January 2, 2006');
```

---

## Common Queries

### 1. Get Upcoming Shows with Venue and Bands

```sql
SELECT
    s.id,
    s.title,
    s.date,
    s.price_min,
    s.price_max,
    v.name AS venue_name,
    v.region AS venue_region,
    json_agg(
        json_build_object(
            'name', b.name,
            'is_headliner', sb.is_headliner
        ) ORDER BY sb.performance_order DESC
    ) AS bands
FROM shows s
JOIN venues v ON s.venue_id = v.id
JOIN show_bands sb ON s.id = sb.show_id
JOIN bands b ON sb.band_id = b.id
WHERE s.date >= NOW()
  AND s.status = 'scheduled'
GROUP BY s.id, v.name, v.region
ORDER BY s.date
LIMIT 20;
```

---

### 2. Filter Shows by Genre

```sql
SELECT DISTINCT s.*
FROM shows s
JOIN show_bands sb ON s.id = sb.show_id
JOIN band_genres bg ON sb.band_id = bg.band_id
JOIN genres g ON bg.genre_id = g.id
WHERE g.slug IN ('rock', 'indie')
  AND s.date >= NOW()
  AND s.status = 'scheduled'
ORDER BY s.date;
```

---

### 3. Find Similar Bands (by shared genres)

```sql
-- Find bands similar to Moon Taxi (id=1)
SELECT
    b.id,
    b.name,
    COUNT(*) AS shared_genres,
    array_agg(g.name) AS genres
FROM bands b
JOIN band_genres bg1 ON b.id = bg1.band_id
JOIN band_genres bg2 ON bg1.genre_id = bg2.genre_id
JOIN genres g ON bg1.genre_id = g.id
WHERE bg2.band_id = 1  -- Moon Taxi
  AND b.id != 1        -- Exclude Moon Taxi itself
GROUP BY b.id, b.name
ORDER BY shared_genres DESC, b.name
LIMIT 10;
```

---

### 4. Get Shows at Venue with Filters

```sql
SELECT s.*, v.name AS venue_name
FROM shows s
JOIN venues v ON s.venue_id = v.id
WHERE v.slug = 'the-orange-peel'
  AND s.date BETWEEN '2025-11-01' AND '2025-11-30'
  AND s.status = 'scheduled'
ORDER BY s.date;
```

---

### 5. Search Shows/Bands/Venues (Full-Text Search)

```sql
-- Search for "moon"
SELECT 'show' AS type, id, title AS name, date
FROM shows
WHERE to_tsvector('english', title) @@ to_tsquery('english', 'moon')

UNION ALL

SELECT 'band' AS type, id, name, NULL AS date
FROM bands
WHERE to_tsvector('english', name) @@ to_tsquery('english', 'moon')

UNION ALL

SELECT 'venue' AS type, id, name, NULL AS date
FROM venues
WHERE to_tsvector('english', name) @@ to_tsquery('english', 'moon')

ORDER BY type, name;
```

---

### 6. Count Shows by Venue (Analytics)

```sql
SELECT
    v.name,
    v.region,
    COUNT(s.id) AS show_count
FROM venues v
LEFT JOIN shows s ON v.id = s.venue_id AND s.date >= NOW() - INTERVAL '1 year'
GROUP BY v.id, v.name, v.region
ORDER BY show_count DESC
LIMIT 10;
```

---

### 7. Band Performance Count

```sql
SELECT
    b.name,
    COUNT(sb.show_id) AS performance_count
FROM bands b
JOIN show_bands sb ON b.id = sb.band_id
JOIN shows s ON sb.show_id = s.id
WHERE s.date >= NOW() - INTERVAL '1 year'
GROUP BY b.id, b.name
ORDER BY performance_count DESC
LIMIT 20;
```

---

## Data Integrity & Constraints

### Foreign Keys
- Prevent orphaned records (e.g., can't delete venue with shows)
- Cascade deletes where appropriate

### Unique Constraints
- `venues.slug`, `bands.slug`, `articles.slug` (SEO-friendly URLs)
- `genres.name` (prevent duplicates)

### Check Constraints

```sql
ALTER TABLE shows ADD CONSTRAINT check_price_min_positive
    CHECK (price_min IS NULL OR price_min >= 0);

ALTER TABLE shows ADD CONSTRAINT check_price_max_positive
    CHECK (price_max IS NULL OR price_max >= 0);

ALTER TABLE shows ADD CONSTRAINT check_price_range_valid
    CHECK (price_max IS NULL OR price_min IS NULL OR price_max >= price_min);

ALTER TABLE shows ADD CONSTRAINT check_status_valid
    CHECK (status IN ('scheduled', 'cancelled', 'postponed', 'completed'));

ALTER TABLE shows ADD CONSTRAINT check_source_valid
    CHECK (source IN ('scraped', 'band_submitted', 'manual'));

ALTER TABLE shows ADD CONSTRAINT check_age_restriction_valid
    CHECK (age_restriction IS NULL OR age_restriction IN ('All Ages', '18+', '21+'));

ALTER TABLE venues ADD CONSTRAINT check_capacity_positive
    CHECK (capacity IS NULL OR capacity > 0);
```

---

## Indexes Strategy

### High Priority (Query Performance)
- `shows.date` - Most queries filter by date
- `shows.venue_id` - Common JOIN
- `show_bands.show_id`, `show_bands.band_id` - Many-to-many lookups
- `band_genres.band_id`, `band_genres.genre_id` - Similar artists query

### Medium Priority
- Full-text search indexes (GIN indexes on tsvector)
- JSONB indexes (for metadata queries)

### Low Priority (Future Optimization)
- Materialized views for complex aggregations
- Partial indexes for specific query patterns

---

## Sample Data Size Estimates

| Table | Initial | Year 1 | Year 2 |
|-------|---------|--------|--------|
| venues | 20 | 30 | 50 |
| bands | 500 | 2,000 | 5,000 |
| genres | 20 | 30 | 30 |
| shows | 1,000 | 5,000 | 20,000 |
| band_genres | 1,000 | 4,000 | 10,000 |
| show_bands | 2,000 | 10,000 | 40,000 |
| articles | 10 | 50 | 100 |

**Total Storage Estimate**: ~50MB (Year 1), ~200MB (Year 2)
**Well within Neon free tier (0.5GB)**

---

## Migration Strategy

### Initial Schema

```bash
# migrations/000001_init_schema.up.sql
# Contains all CREATE TABLE statements

# migrations/000001_init_schema.down.sql
DROP TABLE IF EXISTS show_bands CASCADE;
DROP TABLE IF EXISTS band_genres CASCADE;
DROP TABLE IF EXISTS shows CASCADE;
DROP TABLE IF EXISTS articles CASCADE;
DROP TABLE IF EXISTS bands CASCADE;
DROP TABLE IF EXISTS genres CASCADE;
DROP TABLE IF EXISTS venues CASCADE;
DROP TABLE IF EXISTS venue_scrapers CASCADE;
```

### Seed Data

```bash
# migrations/000002_seed_genres.up.sql
INSERT INTO genres (name, slug) VALUES
('Rock', 'rock'),
('Indie', 'indie'),
...
```

### Running Migrations

```bash
migrate -path migrations -database "$DATABASE_URL" up
```

---

## Future Enhancements

### Phase 2: User Features

```sql
-- User favorites (shows, bands, venues)
CREATE TABLE user_favorites (
    user_id INTEGER REFERENCES users(id),
    favoritable_type TEXT, -- 'show', 'band', 'venue'
    favoritable_id INTEGER,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, favoritable_type, favoritable_id)
);

-- User follows (bands, venues)
CREATE TABLE user_follows (
    user_id INTEGER REFERENCES users(id),
    followable_type TEXT, -- 'band', 'venue'
    followable_id INTEGER,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, followable_type, followable_id)
);
```

### Phase 3: Advanced Features

```sql
-- Reviews/ratings
CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    show_id INTEGER REFERENCES shows(id),
    rating INTEGER CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Notifications
CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    type TEXT, -- 'new_show', 'show_reminder', 'band_announced'
    data JSONB,
    read_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);
```

---

## Database Maintenance

### Backups
- Neon: Automatic daily snapshots
- Manual: Weekly `pg_dump` to Cloud Storage

### Vacuum & Analyze
```sql
-- Run weekly (or let Neon auto-vacuum handle it)
VACUUM ANALYZE shows;
VACUUM ANALYZE bands;
```

### Monitoring Queries
```sql
-- Slow queries
SELECT query, mean_exec_time, calls
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;

-- Table sizes
SELECT
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

---

This schema supports all MVP features and scales to future enhancements without major rewrites.
