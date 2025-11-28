-- The Asheville Setlist - Initial Schema Migration
-- Creates all core tables for concert discovery platform

-- ============================================
-- 1. VENUES
-- ============================================
CREATE TABLE venues (
    id SERIAL PRIMARY KEY,

    -- Basic info
    name TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,

    -- Location
    address TEXT,
    city TEXT DEFAULT 'Asheville',
    state TEXT DEFAULT 'NC',
    zip_code TEXT,
    region TEXT,
    latitude DECIMAL(9,6),
    longitude DECIMAL(9,6),

    -- Venue details
    capacity INTEGER,
    website TEXT,
    phone TEXT,
    image_url TEXT,

    -- Flexible metadata
    metadata JSONB,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Venues indexes
CREATE INDEX idx_venues_region ON venues(region);
CREATE INDEX idx_venues_slug ON venues(slug);
CREATE INDEX idx_venues_metadata ON venues USING GIN(metadata);
CREATE INDEX idx_venues_search ON venues USING GIN(to_tsvector('english', name));

-- Venues constraints
ALTER TABLE venues ADD CONSTRAINT check_capacity_positive
    CHECK (capacity IS NULL OR capacity > 0);

-- ============================================
-- 2. BANDS
-- ============================================
CREATE TABLE bands (
    id SERIAL PRIMARY KEY,

    -- Basic info
    name TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,

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

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Bands indexes
CREATE UNIQUE INDEX idx_bands_slug ON bands(slug);
CREATE INDEX idx_bands_name ON bands(name);
CREATE INDEX idx_bands_search ON bands USING GIN(to_tsvector('english', name || ' ' || COALESCE(bio, '')));

-- ============================================
-- 3. GENRES
-- ============================================
CREATE TABLE genres (
    id SERIAL PRIMARY KEY,

    name TEXT NOT NULL UNIQUE,
    slug TEXT NOT NULL UNIQUE,
    description TEXT,

    -- Timestamp
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Genres indexes
CREATE INDEX idx_genres_slug ON genres(slug);

-- ============================================
-- 4. SHOWS
-- ============================================
CREATE TABLE shows (
    id SERIAL PRIMARY KEY,

    -- Relationships
    venue_id INTEGER NOT NULL REFERENCES venues(id) ON DELETE CASCADE,

    -- Basic info
    title TEXT,
    description TEXT,
    image_url TEXT,

    -- Date/time
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    doors_time TIME,
    show_time TIME,

    -- Ticketing
    price_min NUMERIC(10,2),
    price_max NUMERIC(10,2),
    ticket_url TEXT,
    age_restriction TEXT,

    -- Status & Source
    status TEXT DEFAULT 'scheduled',
    source TEXT DEFAULT 'scraped',

    -- Scraped data
    scraped_data JSONB,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Shows indexes
CREATE INDEX idx_shows_venue ON shows(venue_id);
CREATE INDEX idx_shows_date ON shows(date);
CREATE INDEX idx_shows_status ON shows(status);
CREATE INDEX idx_shows_source ON shows(source);
CREATE INDEX idx_shows_date_venue ON shows(date, venue_id);
CREATE INDEX idx_shows_scraped_data ON shows USING GIN(scraped_data);
CREATE INDEX idx_shows_upcoming ON shows(date) WHERE status = 'scheduled';

-- Shows constraints
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

-- ============================================
-- 5. BAND_GENRES (Junction Table)
-- ============================================
CREATE TABLE band_genres (
    band_id INTEGER NOT NULL REFERENCES bands(id) ON DELETE CASCADE,
    genre_id INTEGER NOT NULL REFERENCES genres(id) ON DELETE CASCADE,

    -- Timestamp
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    PRIMARY KEY (band_id, genre_id)
);

-- Band_genres indexes
CREATE INDEX idx_band_genres_band ON band_genres(band_id);
CREATE INDEX idx_band_genres_genre ON band_genres(genre_id);

-- ============================================
-- 6. SHOW_BANDS (Junction Table)
-- ============================================
CREATE TABLE show_bands (
    show_id INTEGER NOT NULL REFERENCES shows(id) ON DELETE CASCADE,
    band_id INTEGER NOT NULL REFERENCES bands(id) ON DELETE CASCADE,

    -- Lineup details
    is_headliner BOOLEAN DEFAULT FALSE,
    performance_order INTEGER,

    -- Timestamp
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    PRIMARY KEY (show_id, band_id)
);

-- Show_bands indexes
CREATE INDEX idx_show_bands_show ON show_bands(show_id);
CREATE INDEX idx_show_bands_band ON show_bands(band_id);

-- ============================================
-- 7. ARTICLES
-- ============================================
CREATE TABLE articles (
    id SERIAL PRIMARY KEY,

    -- Content
    title TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    content TEXT NOT NULL,
    excerpt TEXT,

    -- Metadata
    author TEXT,
    cover_image_url TEXT,

    -- Publishing
    published_at TIMESTAMP WITH TIME ZONE,
    is_published BOOLEAN DEFAULT FALSE,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Articles indexes
CREATE INDEX idx_articles_slug ON articles(slug);
CREATE INDEX idx_articles_published ON articles(published_at) WHERE is_published = TRUE;
CREATE INDEX idx_articles_search ON articles USING GIN(to_tsvector('english', title || ' ' || COALESCE(excerpt, '')));

-- ============================================
-- 8. VENUE_SCRAPERS
-- ============================================
CREATE TABLE venue_scrapers (
    id SERIAL PRIMARY KEY,
    venue_id INTEGER NOT NULL REFERENCES venues(id) ON DELETE CASCADE,

    -- Scraper config
    url TEXT NOT NULL,
    scraper_type TEXT NOT NULL,

    -- CSS selectors
    selectors JSONB NOT NULL,

    -- Parsing rules
    date_format TEXT,

    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    last_scraped_at TIMESTAMP WITH TIME ZONE,
    last_success_at TIMESTAMP WITH TIME ZONE,
    error_count INTEGER DEFAULT 0,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Venue_scrapers indexes
CREATE INDEX idx_venue_scrapers_venue ON venue_scrapers(venue_id);
CREATE INDEX idx_venue_scrapers_active ON venue_scrapers(is_active) WHERE is_active = TRUE;
