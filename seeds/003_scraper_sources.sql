-- Seed file: scraper_sources
-- Data source configuration for The Asheville Setlist
-- Run: psql $DATABASE_URL < seeds/003_scraper_sources.sql

BEGIN;

-- =============================================================================
-- SCRAPER SOURCES TABLE (extends venue_scrapers for non-venue sources)
-- =============================================================================
-- Note: Live Music Asheville is the PRIMARY source - covers 400+ venues
-- Other sources are for enrichment or venues not on LMA

-- For now, we use venue_scrapers table but treat LMA as a special case
-- since it's an aggregator, not a single venue

-- =============================================================================
-- LIVE MUSIC ASHEVILLE (Primary Aggregator)
-- =============================================================================
-- This is handled specially in code since it covers ALL venues
-- API Endpoints:
--   Events: https://livemusicasheville.com/wp-json/tribe/events/v1/events
--   Venues: https://livemusicasheville.com/wp-json/tribe/events/v1/venues
--   Categories: https://livemusicasheville.com/wp-json/tribe/events/v1/categories
--
-- Query Parameters:
--   ?per_page=100 (max)
--   ?page=1,2,3...
--   ?start_date=2025-01-01
--   ?end_date=2025-12-31
--   ?status=publish
--
-- Total: 3421 events, 427 venues, 685 pages

-- Create a special "source" record for tracking LMA scraper status
INSERT INTO venue_scrapers (venue_id, url, scraper_type, selectors, date_format, is_active)
SELECT 
  v.id,
  'https://livemusicasheville.com/wp-json/tribe/events/v1/events',
  'api',
  '{
    "type": "the_events_calendar_api",
    "base_url": "https://livemusicasheville.com/wp-json/tribe/events/v1",
    "endpoints": {
      "events": "/events",
      "venues": "/venues",
      "categories": "/categories"
    },
    "pagination": {
      "per_page": 100,
      "max_pages": 100
    }
  }'::jsonb,
  'ISO8601',
  true
FROM venues v
WHERE v.slug = 'the-orange-peel'  -- Placeholder: LMA covers all venues
ON CONFLICT DO NOTHING;

-- =============================================================================
-- INDIVIDUAL VENUE SCRAPERS (Fallback / Enrichment)
-- =============================================================================
-- These are for venues that might have exclusive events not on LMA
-- or for getting additional data (prices, age restrictions)

-- The Orange Peel (Rockhouse/ETIX platform)
INSERT INTO venue_scrapers (venue_id, url, scraper_type, selectors, date_format, is_active)
SELECT 
  v.id,
  'https://theorangepeel.net/events/',
  'static',
  '{
    "platform": "rockhouse",
    "container": ".event-item",
    "title": ".event-title",
    "date": ".event-date",
    "time": ".event-time",
    "price": ".ticket-price",
    "ticket_url": ".buy-tickets a",
    "age_restriction": ".age-restriction"
  }'::jsonb,
  'F j, Y',
  false  -- Disabled: use LMA as primary
FROM venues v
WHERE v.slug = 'the-orange-peel';

-- Asheville Music Hall (custom WordPress)
INSERT INTO venue_scrapers (venue_id, url, scraper_type, selectors, date_format, is_active)
SELECT 
  v.id,
  'https://ashevillemusichall.com/all-shows/',
  'static',
  '{
    "platform": "wordpress_wfea",
    "venue_ids": ["43374215", "64643809", "43374223"],
    "container": ".event-card",
    "title": ".event-title",
    "date": ".event-date",
    "time": ".event-time"
  }'::jsonb,
  'F j, Y',
  false  -- Disabled: use LMA as primary
FROM venues v
WHERE v.slug = 'asheville-music-hall';

COMMIT;

-- =============================================================================
-- SCRAPER PRIORITY DOCUMENTATION
-- =============================================================================
-- 
-- Priority 1: Live Music Asheville REST API
--   - Covers 427 venues, 3400+ events
--   - Structured JSON data
--   - Run: Every 6 hours
--   - No bot blocking (tested)
--
-- Priority 2: Bandsintown API (for band enrichment)
--   - Get band images, Spotify links, social media
--   - Query by artist name after scraping LMA
--   - Free API with app_id
--
-- Priority 3: Direct venue scrapers (disabled by default)
--   - Enable only for venues with exclusive events
--   - Or for additional data not in LMA (exact prices, age restrictions)
--
-- =============================================================================
