-- The Asheville Setlist - Initial Schema Rollback
-- Drops all tables in reverse order of dependencies

DROP TABLE IF EXISTS venue_scrapers CASCADE;
DROP TABLE IF EXISTS articles CASCADE;
DROP TABLE IF EXISTS show_bands CASCADE;
DROP TABLE IF EXISTS band_genres CASCADE;
DROP TABLE IF EXISTS shows CASCADE;
DROP TABLE IF EXISTS genres CASCADE;
DROP TABLE IF EXISTS bands CASCADE;
DROP TABLE IF EXISTS venues CASCADE;
