// Package testutil provides utilities for integration testing
package testutil

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paulsena/asheville-setlist/internal/db"
)

// TestDB holds the test database connection and queries
type TestDB struct {
	Pool    *pgxpool.Pool
	Queries *db.Queries
}

// SetupTestDB creates a connection to the test database
// Requires DATABASE_URL environment variable to be set
func SetupTestDB(t *testing.T) *TestDB {
	t.Helper()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		t.Fatalf("failed to parse database URL: %v", err)
	}

	// Use smaller pool for tests
	poolConfig.MaxConns = 5
	poolConfig.MinConns = 1

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		t.Fatalf("failed to ping database: %v", err)
	}

	queries := db.New(pool)

	return &TestDB{
		Pool:    pool,
		Queries: queries,
	}
}

// Close closes the test database connection
func (tdb *TestDB) Close() {
	if tdb.Pool != nil {
		tdb.Pool.Close()
	}
}

// TestShowTitlePrefix is used to identify test shows for cleanup
const TestShowTitlePrefix = "[TEST] "

// CleanupTestData removes test data created during tests
// This deletes shows with test title prefix and bands with "Test Band" prefix
func (tdb *TestDB) CleanupTestData(ctx context.Context) error {
	// Delete show_bands for test shows first (due to FK constraints)
	_, err := tdb.Pool.Exec(ctx, `
		DELETE FROM show_bands
		WHERE show_id IN (SELECT id FROM shows WHERE title LIKE '[TEST]%')
	`)
	if err != nil {
		return fmt.Errorf("failed to delete test show_bands: %w", err)
	}

	// Delete test shows (identified by title prefix)
	_, err = tdb.Pool.Exec(ctx, `DELETE FROM shows WHERE title LIKE '[TEST]%'`)
	if err != nil {
		return fmt.Errorf("failed to delete test shows: %w", err)
	}

	// Delete band_genres for test bands
	_, err = tdb.Pool.Exec(ctx, `
		DELETE FROM band_genres
		WHERE band_id IN (SELECT id FROM bands WHERE name LIKE 'Test Band%')
	`)
	if err != nil {
		return fmt.Errorf("failed to delete test band_genres: %w", err)
	}

	// Delete test bands
	_, err = tdb.Pool.Exec(ctx, `DELETE FROM bands WHERE name LIKE 'Test Band%'`)
	if err != nil {
		return fmt.Errorf("failed to delete test bands: %w", err)
	}

	return nil
}

// InsertTestShow inserts a test show and returns its ID
// The title is automatically prefixed with TestShowTitlePrefix for cleanup
func (tdb *TestDB) InsertTestShow(ctx context.Context, venueID int32, date time.Time, title string) (int32, error) {
	source := "manual" // Use valid source per database constraint
	status := "scheduled"
	testTitle := TestShowTitlePrefix + title
	var showID int32
	err := tdb.Pool.QueryRow(ctx, `
		INSERT INTO shows (venue_id, date, title, status, source)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, venueID, date, testTitle, status, source).Scan(&showID)
	return showID, err
}

// InsertTestBand inserts a test band and returns its ID
func (tdb *TestDB) InsertTestBand(ctx context.Context, name, slug string) (int32, error) {
	var bandID int32
	err := tdb.Pool.QueryRow(ctx, `
		INSERT INTO bands (name, slug)
		VALUES ($1, $2)
		RETURNING id
	`, name, slug).Scan(&bandID)
	return bandID, err
}

// LinkBandToShow links a band to a show
func (tdb *TestDB) LinkBandToShow(ctx context.Context, showID, bandID int32, isHeadliner bool, order int) error {
	_, err := tdb.Pool.Exec(ctx, `
		INSERT INTO show_bands (show_id, band_id, is_headliner, performance_order)
		VALUES ($1, $2, $3, $4)
	`, showID, bandID, isHeadliner, order)
	return err
}

// AddGenreToBand adds a genre to a band
func (tdb *TestDB) AddGenreToBand(ctx context.Context, bandID, genreID int32) error {
	_, err := tdb.Pool.Exec(ctx, `
		INSERT INTO band_genres (band_id, genre_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, bandID, genreID)
	return err
}

// GetFirstVenueID returns the ID of the first venue in the database
func (tdb *TestDB) GetFirstVenueID(ctx context.Context) (int32, error) {
	var venueID int32
	err := tdb.Pool.QueryRow(ctx, `SELECT id FROM venues LIMIT 1`).Scan(&venueID)
	return venueID, err
}

// GetFirstGenreID returns the ID of the first genre in the database
func (tdb *TestDB) GetFirstGenreID(ctx context.Context) (int32, error) {
	var genreID int32
	err := tdb.Pool.QueryRow(ctx, `SELECT id FROM genres LIMIT 1`).Scan(&genreID)
	return genreID, err
}

// GetVenueBySlug returns venue info by slug
func (tdb *TestDB) GetVenueBySlug(ctx context.Context, slug string) (int32, string, error) {
	var id int32
	var name string
	err := tdb.Pool.QueryRow(ctx, `SELECT id, name FROM venues WHERE slug = $1`, slug).Scan(&id, &name)
	return id, name, err
}
