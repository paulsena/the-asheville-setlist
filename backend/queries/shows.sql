-- ============================================
-- SHOWS QUERIES
-- ============================================

-- name: GetShowByID :one
-- Get single show with venue info
SELECT
    s.id,
    s.title,
    s.description,
    s.image_url,
    s.date,
    s.doors_time,
    s.show_time,
    s.price_min,
    s.price_max,
    s.ticket_url,
    s.age_restriction,
    s.status,
    s.source,
    s.created_at,
    s.updated_at,
    v.id AS venue_id,
    v.name AS venue_name,
    v.slug AS venue_slug,
    v.address AS venue_address,
    v.region AS venue_region,
    v.website AS venue_website,
    v.image_url AS venue_image_url
FROM shows s
JOIN venues v ON s.venue_id = v.id
WHERE s.id = $1;

-- name: ListUpcomingShows :many
-- List upcoming scheduled shows with pagination
-- Used for homepage and general show listing
SELECT
    s.id,
    s.title,
    s.image_url,
    s.date,
    s.doors_time,
    s.show_time,
    s.price_min,
    s.price_max,
    s.ticket_url,
    s.age_restriction,
    s.status,
    v.id AS venue_id,
    v.name AS venue_name,
    v.slug AS venue_slug,
    v.region AS venue_region,
    v.address AS venue_address,
    v.image_url AS venue_image_url,
    COUNT(*) OVER() AS total_count
FROM shows s
JOIN venues v ON s.venue_id = v.id
WHERE s.date >= NOW()
  AND s.status = 'scheduled'
ORDER BY s.date ASC, s.id ASC
LIMIT $1 OFFSET $2;

-- name: ListShowsByDateRange :many
-- Filter shows by date range (inclusive)
SELECT
    s.id,
    s.title,
    s.image_url,
    s.date,
    s.doors_time,
    s.show_time,
    s.price_min,
    s.price_max,
    s.ticket_url,
    s.age_restriction,
    s.status,
    v.id AS venue_id,
    v.name AS venue_name,
    v.slug AS venue_slug,
    v.region AS venue_region,
    v.address AS venue_address,
    v.image_url AS venue_image_url,
    COUNT(*) OVER() AS total_count
FROM shows s
JOIN venues v ON s.venue_id = v.id
WHERE s.date >= $1
  AND s.date <= $2
  AND s.status = 'scheduled'
ORDER BY s.date ASC, s.id ASC
LIMIT $3 OFFSET $4;

-- name: ListShowsByVenue :many
-- Filter shows by venue slug(s)
SELECT
    s.id,
    s.title,
    s.image_url,
    s.date,
    s.doors_time,
    s.show_time,
    s.price_min,
    s.price_max,
    s.ticket_url,
    s.age_restriction,
    s.status,
    v.id AS venue_id,
    v.name AS venue_name,
    v.slug AS venue_slug,
    v.region AS venue_region,
    v.address AS venue_address,
    v.image_url AS venue_image_url,
    COUNT(*) OVER() AS total_count
FROM shows s
JOIN venues v ON s.venue_id = v.id
WHERE s.date >= NOW()
  AND s.status = 'scheduled'
  AND v.slug = ANY($1::text[])
ORDER BY s.date ASC, s.id ASC
LIMIT $2 OFFSET $3;

-- name: ListShowsByRegion :many
-- Filter shows by region(s)
SELECT
    s.id,
    s.title,
    s.image_url,
    s.date,
    s.doors_time,
    s.show_time,
    s.price_min,
    s.price_max,
    s.ticket_url,
    s.age_restriction,
    s.status,
    v.id AS venue_id,
    v.name AS venue_name,
    v.slug AS venue_slug,
    v.region AS venue_region,
    v.address AS venue_address,
    v.image_url AS venue_image_url,
    COUNT(*) OVER() AS total_count
FROM shows s
JOIN venues v ON s.venue_id = v.id
WHERE s.date >= NOW()
  AND s.status = 'scheduled'
  AND v.region = ANY($1::text[])
ORDER BY s.date ASC, s.id ASC
LIMIT $2 OFFSET $3;

-- name: ListShowsByGenre :many
-- Filter shows by genre slug(s) - shows with bands matching any of the genres
SELECT DISTINCT ON (s.date, s.id)
    s.id,
    s.title,
    s.image_url,
    s.date,
    s.doors_time,
    s.show_time,
    s.price_min,
    s.price_max,
    s.ticket_url,
    s.age_restriction,
    s.status,
    v.id AS venue_id,
    v.name AS venue_name,
    v.slug AS venue_slug,
    v.region AS venue_region,
    v.address AS venue_address,
    v.image_url AS venue_image_url
FROM shows s
JOIN venues v ON s.venue_id = v.id
JOIN show_bands sb ON s.id = sb.show_id
JOIN band_genres bg ON sb.band_id = bg.band_id
JOIN genres g ON bg.genre_id = g.id
WHERE s.date >= NOW()
  AND s.status = 'scheduled'
  AND g.slug = ANY($1::text[])
ORDER BY s.date ASC, s.id ASC
LIMIT $2 OFFSET $3;

-- name: ListShowsByPriceRange :many
-- Filter shows by price range
SELECT
    s.id,
    s.title,
    s.image_url,
    s.date,
    s.doors_time,
    s.show_time,
    s.price_min,
    s.price_max,
    s.ticket_url,
    s.age_restriction,
    s.status,
    v.id AS venue_id,
    v.name AS venue_name,
    v.slug AS venue_slug,
    v.region AS venue_region,
    v.address AS venue_address,
    v.image_url AS venue_image_url,
    COUNT(*) OVER() AS total_count
FROM shows s
JOIN venues v ON s.venue_id = v.id
WHERE s.date >= NOW()
  AND s.status = 'scheduled'
  AND (s.price_min IS NULL OR s.price_min >= $1)
  AND (s.price_max IS NULL OR s.price_max <= $2)
ORDER BY s.date ASC, s.id ASC
LIMIT $3 OFFSET $4;

-- name: ListFreeShows :many
-- Shows that are free (price_min is NULL or 0)
SELECT
    s.id,
    s.title,
    s.image_url,
    s.date,
    s.doors_time,
    s.show_time,
    s.price_min,
    s.price_max,
    s.ticket_url,
    s.age_restriction,
    s.status,
    v.id AS venue_id,
    v.name AS venue_name,
    v.slug AS venue_slug,
    v.region AS venue_region,
    v.address AS venue_address,
    v.image_url AS venue_image_url,
    COUNT(*) OVER() AS total_count
FROM shows s
JOIN venues v ON s.venue_id = v.id
WHERE s.date >= NOW()
  AND s.status = 'scheduled'
  AND (s.price_min IS NULL OR s.price_min = 0)
ORDER BY s.date ASC, s.id ASC
LIMIT $1 OFFSET $2;

-- name: ListShowsTonight :many
-- Shows happening today
SELECT
    s.id,
    s.title,
    s.image_url,
    s.date,
    s.doors_time,
    s.show_time,
    s.price_min,
    s.price_max,
    s.ticket_url,
    s.age_restriction,
    s.status,
    v.id AS venue_id,
    v.name AS venue_name,
    v.slug AS venue_slug,
    v.region AS venue_region,
    v.address AS venue_address,
    v.image_url AS venue_image_url
FROM shows s
JOIN venues v ON s.venue_id = v.id
WHERE DATE(s.date AT TIME ZONE 'America/New_York') = DATE(NOW() AT TIME ZONE 'America/New_York')
  AND s.status = 'scheduled'
ORDER BY s.date ASC, s.id ASC;

-- name: ListShowsThisWeekend :many
-- Shows happening Friday-Sunday of current or next weekend
SELECT
    s.id,
    s.title,
    s.image_url,
    s.date,
    s.doors_time,
    s.show_time,
    s.price_min,
    s.price_max,
    s.ticket_url,
    s.age_restriction,
    s.status,
    v.id AS venue_id,
    v.name AS venue_name,
    v.slug AS venue_slug,
    v.region AS venue_region,
    v.address AS venue_address,
    v.image_url AS venue_image_url
FROM shows s
JOIN venues v ON s.venue_id = v.id
WHERE s.date >= DATE_TRUNC('week', NOW()) + INTERVAL '4 days'  -- Friday
  AND s.date < DATE_TRUNC('week', NOW()) + INTERVAL '8 days'   -- Monday
  AND s.status = 'scheduled'
ORDER BY s.date ASC, s.id ASC;

-- name: GetShowBands :many
-- Get all bands for a show with their genres
SELECT
    b.id,
    b.name,
    b.slug,
    b.bio,
    b.image_url,
    b.spotify_url,
    b.website,
    sb.is_headliner,
    sb.performance_order
FROM bands b
JOIN show_bands sb ON b.id = sb.band_id
WHERE sb.show_id = $1
ORDER BY sb.performance_order DESC NULLS LAST, sb.is_headliner DESC;

-- name: GetBandGenresForShow :many
-- Get genres for a band (used when fetching show details)
SELECT
    g.id,
    g.name,
    g.slug
FROM genres g
JOIN band_genres bg ON g.id = bg.genre_id
WHERE bg.band_id = $1
ORDER BY g.name;

-- name: CountShowsByGenre :one
-- Count upcoming shows by genre (for genre filter)
SELECT COUNT(DISTINCT s.id)
FROM shows s
JOIN show_bands sb ON s.id = sb.show_id
JOIN band_genres bg ON sb.band_id = bg.band_id
JOIN genres g ON bg.genre_id = g.id
WHERE s.date >= NOW()
  AND s.status = 'scheduled'
  AND g.slug = ANY($1::text[]);

-- name: CreateShow :one
-- Create a new show (band submission)
INSERT INTO shows (
    venue_id,
    title,
    image_url,
    date,
    doors_time,
    show_time,
    price_min,
    price_max,
    ticket_url,
    age_restriction,
    status,
    source
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
) RETURNING id, status, created_at;

-- name: CreateShowBand :exec
-- Link a band to a show
INSERT INTO show_bands (show_id, band_id, is_headliner, performance_order)
VALUES ($1, $2, $3, $4);

-- name: SearchShows :many
-- Full-text search on show titles
SELECT
    s.id,
    s.title,
    s.date,
    v.name AS venue_name
FROM shows s
JOIN venues v ON s.venue_id = v.id
WHERE s.date >= NOW()
  AND s.status = 'scheduled'
  AND to_tsvector('english', COALESCE(s.title, '')) @@ plainto_tsquery('english', $1)
ORDER BY s.date ASC
LIMIT $2;
