-- ============================================
-- BANDS QUERIES
-- ============================================

-- name: GetBand :one
-- Get band by ID
SELECT * FROM bands
WHERE id = $1 LIMIT 1;

-- name: GetBandBySlug :one
-- Get band by slug for detail page
SELECT
    id,
    name,
    slug,
    bio,
    hometown,
    image_url,
    website,
    spotify_url,
    instagram,
    facebook,
    bandcamp_url,
    metadata,
    created_at,
    updated_at
FROM bands
WHERE slug = $1 LIMIT 1;

-- name: GetBandByName :one
-- Get band by name (case-insensitive) for matching during submissions
SELECT * FROM bands
WHERE LOWER(name) = LOWER($1) LIMIT 1;

-- name: ListBands :many
-- List bands with pagination
SELECT
    b.id,
    b.name,
    b.slug,
    b.bio,
    b.hometown,
    b.image_url,
    COUNT(*) OVER() AS total_count
FROM bands b
ORDER BY b.name
LIMIT $1 OFFSET $2;

-- name: ListBandsByGenre :many
-- List bands filtered by genre slug(s) with pagination
SELECT DISTINCT ON (b.name, b.id)
    b.id,
    b.name,
    b.slug,
    b.bio,
    b.hometown,
    b.image_url
FROM bands b
JOIN band_genres bg ON b.id = bg.band_id
JOIN genres g ON bg.genre_id = g.id
WHERE g.slug = ANY($1::text[])
ORDER BY b.name, b.id
LIMIT $2 OFFSET $3;

-- name: CountBandsByGenre :one
-- Count bands by genre (for pagination)
SELECT COUNT(DISTINCT b.id)
FROM bands b
JOIN band_genres bg ON b.id = bg.band_id
JOIN genres g ON bg.genre_id = g.id
WHERE g.slug = ANY($1::text[]);

-- name: GetBandGenres :many
-- Get genres for a band
SELECT
    g.id,
    g.name,
    g.slug
FROM genres g
JOIN band_genres bg ON g.id = bg.genre_id
WHERE bg.band_id = $1
ORDER BY g.name;

-- name: GetBandGenresBatch :many
-- Get genres for multiple bands (batch load)
SELECT
    bg.band_id,
    g.id,
    g.name,
    g.slug
FROM genres g
JOIN band_genres bg ON g.id = bg.genre_id
WHERE bg.band_id = ANY($1::int[])
ORDER BY bg.band_id, g.name;

-- name: GetBandUpcomingShows :many
-- Get upcoming shows for a band
SELECT
    s.id,
    s.date,
    s.title,
    v.id AS venue_id,
    v.name AS venue_name,
    v.slug AS venue_slug,
    sb.is_headliner
FROM shows s
JOIN show_bands sb ON s.id = sb.show_id
JOIN venues v ON s.venue_id = v.id
WHERE sb.band_id = $1
  AND s.status = 'scheduled'
  AND s.date >= NOW()
ORDER BY s.date ASC;

-- name: GetSimilarBands :many
-- Find bands with shared genres (similar bands)
-- Excludes the source band and orders by number of shared genres
SELECT
    b.id,
    b.name,
    b.slug,
    b.image_url,
    COUNT(bg1.genre_id) AS shared_genre_count
FROM bands b
JOIN band_genres bg1 ON b.id = bg1.band_id
JOIN band_genres bg2 ON bg1.genre_id = bg2.genre_id
WHERE bg2.band_id = $1
  AND b.id != $1
GROUP BY b.id, b.name, b.slug, b.image_url
ORDER BY shared_genre_count DESC, b.name ASC
LIMIT $2;

-- name: GetSimilarBandsWithGenres :many
-- Find similar bands with their shared genre names
SELECT
    b.id,
    b.name,
    b.slug,
    b.image_url,
    COUNT(bg1.genre_id) AS shared_genre_count,
    ARRAY_AGG(g.name ORDER BY g.name) AS shared_genres,
    ARRAY_AGG(g.id ORDER BY g.name) AS shared_genre_ids,
    ARRAY_AGG(g.slug ORDER BY g.name) AS shared_genre_slugs
FROM bands b
JOIN band_genres bg1 ON b.id = bg1.band_id
JOIN band_genres bg2 ON bg1.genre_id = bg2.genre_id
JOIN genres g ON bg1.genre_id = g.id
WHERE bg2.band_id = $1
  AND b.id != $1
GROUP BY b.id, b.name, b.slug, b.image_url
ORDER BY shared_genre_count DESC, b.name ASC
LIMIT $2;

-- name: SearchBands :many
-- Full-text search on band name and bio
SELECT
    id,
    name,
    slug,
    bio,
    hometown,
    image_url,
    COUNT(*) OVER() AS total_count
FROM bands
WHERE to_tsvector('english', name || ' ' || COALESCE(bio, '')) @@ plainto_tsquery('english', $1)
ORDER BY name
LIMIT $2 OFFSET $3;

-- name: SearchBandsSimple :many
-- Simple search returning minimal fields (for global search)
SELECT
    id,
    name,
    slug
FROM bands
WHERE to_tsvector('english', name || ' ' || COALESCE(bio, '')) @@ plainto_tsquery('english', $1)
ORDER BY name
LIMIT $2;

-- name: CreateBand :one
-- Create a new band
INSERT INTO bands (name, slug)
VALUES ($1, $2)
RETURNING id, name, slug, created_at;

-- name: CreateBandFull :one
-- Create a new band with all fields
INSERT INTO bands (
    name,
    slug,
    bio,
    hometown,
    image_url,
    website,
    spotify_url,
    instagram,
    facebook,
    bandcamp_url
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: AddBandGenre :exec
-- Add a genre to a band
INSERT INTO band_genres (band_id, genre_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: BandExists :one
-- Check if band exists by slug
SELECT EXISTS(SELECT 1 FROM bands WHERE slug = $1);
