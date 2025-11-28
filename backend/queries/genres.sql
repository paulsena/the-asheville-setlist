-- ============================================
-- GENRES QUERIES
-- ============================================

-- name: GetGenre :one
-- Get genre by ID
SELECT * FROM genres
WHERE id = $1 LIMIT 1;

-- name: GetGenreBySlug :one
-- Get genre by slug
SELECT
    id,
    name,
    slug,
    description,
    created_at
FROM genres
WHERE slug = $1 LIMIT 1;

-- name: ListGenres :many
-- List all genres ordered by name
SELECT
    id,
    name,
    slug,
    description,
    created_at
FROM genres
ORDER BY name;

-- name: ListGenresWithShowCount :many
-- List genres with count of upcoming shows
SELECT
    g.id,
    g.name,
    g.slug,
    g.description,
    COUNT(DISTINCT s.id) AS show_count
FROM genres g
LEFT JOIN band_genres bg ON g.id = bg.genre_id
LEFT JOIN show_bands sb ON bg.band_id = sb.band_id
LEFT JOIN shows s ON sb.show_id = s.id
    AND s.status = 'scheduled'
    AND s.date >= NOW()
GROUP BY g.id
ORDER BY g.name;

-- name: ListGenresWithBandCount :many
-- List genres with count of bands
SELECT
    g.id,
    g.name,
    g.slug,
    g.description,
    COUNT(bg.band_id) AS band_count
FROM genres g
LEFT JOIN band_genres bg ON g.id = bg.genre_id
GROUP BY g.id
ORDER BY g.name;

-- name: GetBandsByGenre :many
-- Get bands for a specific genre (for genre detail page)
SELECT
    b.id,
    b.name,
    b.slug,
    b.image_url,
    b.hometown
FROM bands b
JOIN band_genres bg ON b.id = bg.band_id
WHERE bg.genre_id = $1
ORDER BY b.name
LIMIT $2 OFFSET $3;

-- name: CountBandsInGenre :one
-- Count bands in a genre (for pagination)
SELECT COUNT(*)
FROM band_genres
WHERE genre_id = $1;

-- name: GenreExists :one
-- Check if genre exists by ID
SELECT EXISTS(SELECT 1 FROM genres WHERE id = $1);

-- name: GenreExistsBySlug :one
-- Check if genre exists by slug
SELECT EXISTS(SELECT 1 FROM genres WHERE slug = $1);
