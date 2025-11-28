-- ============================================
-- GLOBAL SEARCH QUERIES
-- ============================================

-- name: SearchAll :many
-- Unified search across shows, bands, and venues
-- Returns results with a type discriminator
-- Note: This uses UNION ALL for efficiency (no deduplication needed)
SELECT
    'show' AS entity_type,
    s.id,
    COALESCE(s.title, '') AS name,
    v.name AS extra_info,
    s.date::text AS date_info
FROM shows s
JOIN venues v ON s.venue_id = v.id
WHERE s.date >= NOW()
  AND s.status = 'scheduled'
  AND to_tsvector('english', COALESCE(s.title, '')) @@ plainto_tsquery('english', $1)

UNION ALL

SELECT
    'band' AS entity_type,
    b.id,
    b.name,
    b.slug AS extra_info,
    NULL AS date_info
FROM bands b
WHERE to_tsvector('english', b.name || ' ' || COALESCE(b.bio, '')) @@ plainto_tsquery('english', $1)

UNION ALL

SELECT
    'venue' AS entity_type,
    v.id,
    v.name,
    v.slug AS extra_info,
    NULL AS date_info
FROM venues v
WHERE to_tsvector('english', v.name) @@ plainto_tsquery('english', $1)

LIMIT $2;

-- name: GlobalSearchShows :many
-- Search shows only (for search endpoint's shows section)
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

-- name: GlobalSearchBands :many
-- Search bands only (for search endpoint's bands section)
SELECT
    id,
    name,
    slug
FROM bands
WHERE to_tsvector('english', name || ' ' || COALESCE(bio, '')) @@ plainto_tsquery('english', $1)
ORDER BY name ASC
LIMIT $2;

-- name: GlobalSearchVenues :many
-- Search venues only (for search endpoint's venues section)
SELECT
    id,
    name,
    slug
FROM venues
WHERE to_tsvector('english', name) @@ plainto_tsquery('english', $1)
ORDER BY name ASC
LIMIT $2;

-- name: SearchShowsWithBands :many
-- Search shows including band names in the search
-- This finds shows where either the title OR any band name matches
SELECT DISTINCT ON (s.date, s.id)
    s.id,
    s.title,
    s.date,
    v.name AS venue_name
FROM shows s
JOIN venues v ON s.venue_id = v.id
LEFT JOIN show_bands sb ON s.id = sb.show_id
LEFT JOIN bands b ON sb.band_id = b.id
WHERE s.date >= NOW()
  AND s.status = 'scheduled'
  AND (
    to_tsvector('english', COALESCE(s.title, '')) @@ plainto_tsquery('english', $1)
    OR to_tsvector('english', COALESCE(b.name, '')) @@ plainto_tsquery('english', $1)
  )
ORDER BY s.date ASC, s.id ASC
LIMIT $2;
