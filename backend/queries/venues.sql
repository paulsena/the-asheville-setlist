-- ============================================
-- VENUES QUERIES
-- ============================================

-- name: GetVenue :one
-- Get venue by ID
SELECT * FROM venues
WHERE id = $1 LIMIT 1;

-- name: GetVenueBySlug :one
-- Get venue by slug for detail page
SELECT
    id,
    name,
    slug,
    address,
    city,
    state,
    zip_code,
    region,
    latitude,
    longitude,
    capacity,
    website,
    phone,
    image_url,
    metadata,
    created_at,
    updated_at
FROM venues
WHERE slug = $1 LIMIT 1;

-- name: ListVenues :many
-- List all venues ordered by name
SELECT * FROM venues
ORDER BY name;

-- name: ListVenuesWithShowCount :many
-- List venues with count of upcoming scheduled shows
SELECT
    v.id,
    v.name,
    v.slug,
    v.address,
    v.region,
    v.capacity,
    v.website,
    v.image_url,
    COUNT(s.id) AS upcoming_show_count
FROM venues v
LEFT JOIN shows s ON v.id = s.venue_id
    AND s.status = 'scheduled'
    AND s.date >= NOW()
GROUP BY v.id
ORDER BY v.name;

-- name: ListVenuesByRegion :many
-- List venues filtered by region(s)
SELECT
    v.id,
    v.name,
    v.slug,
    v.address,
    v.region,
    v.capacity,
    v.website,
    v.image_url,
    COUNT(s.id) AS upcoming_show_count
FROM venues v
LEFT JOIN shows s ON v.id = s.venue_id
    AND s.status = 'scheduled'
    AND s.date >= NOW()
WHERE v.region = ANY($1::text[])
GROUP BY v.id
ORDER BY v.name;

-- name: GetVenueUpcomingShows :many
-- Get upcoming shows for a venue (for venue detail page)
SELECT
    s.id,
    s.title,
    s.date,
    s.price_min,
    s.price_max,
    s.ticket_url,
    s.age_restriction,
    s.status
FROM shows s
WHERE s.venue_id = $1
  AND s.status = 'scheduled'
  AND s.date >= NOW()
ORDER BY s.date ASC
LIMIT $2;

-- name: GetShowBandsForVenue :many
-- Get bands for shows at a venue (batch load for venue detail)
SELECT
    sb.show_id,
    b.id,
    b.name,
    b.slug,
    sb.is_headliner
FROM show_bands sb
JOIN bands b ON sb.band_id = b.id
WHERE sb.show_id = ANY($1::int[])
ORDER BY sb.show_id, sb.is_headliner DESC, sb.performance_order DESC NULLS LAST;

-- name: SearchVenues :many
-- Full-text search on venue names
SELECT
    id,
    name,
    slug
FROM venues
WHERE to_tsvector('english', name) @@ plainto_tsquery('english', $1)
ORDER BY name
LIMIT $2;

-- name: VenueExists :one
-- Check if venue exists by ID (for validation)
SELECT EXISTS(SELECT 1 FROM venues WHERE id = $1);
