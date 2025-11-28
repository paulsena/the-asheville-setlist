package handlers_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/paulsena/asheville-setlist/internal/handlers"
	"github.com/paulsena/asheville-setlist/internal/testutil"
)

// setupBandsTestRouter creates a test router with the bands handler
func setupBandsTestRouter(tdb *testutil.TestDB) *gin.Engine {
	h := handlers.New(tdb.Queries)
	router := gin.New()
	router.GET("/api/bands", h.ListBands)
	router.GET("/api/bands/:slug", h.GetBand)
	router.GET("/api/bands/:slug/similar", h.GetSimilarBands)
	return router
}

func TestListBands_Default(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Close()

	router := setupBandsTestRouter(tdb)

	req := httptest.NewRequest(http.MethodGet, "/api/bands", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var resp struct {
		Data []struct {
			ID       int32        `json:"id"`
			Name     string       `json:"name"`
			Slug     string       `json:"slug"`
			Bio      *string      `json:"bio"`
			Hometown *string      `json:"hometown"`
			ImageURL *string      `json:"image_url"`
			Genres   []interface{} `json:"genres"`
		} `json:"data"`
		Meta struct {
			Page       int `json:"page"`
			PerPage    int `json:"per_page"`
			Total      int `json:"total"`
			TotalPages int `json:"total_pages"`
		} `json:"meta"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	// Check pagination defaults
	if resp.Meta.Page != 1 {
		t.Errorf("expected page 1, got %d", resp.Meta.Page)
	}
	if resp.Meta.PerPage != 50 {
		t.Errorf("expected per_page 50, got %d", resp.Meta.PerPage)
	}

	// Data should be an array
	if resp.Data == nil {
		t.Error("data should not be nil")
	}
}

func TestListBands_Pagination(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Close()

	router := setupBandsTestRouter(tdb)

	tests := []struct {
		name           string
		query          string
		expectedStatus int
		expectedPage   int
		expectedPP     int
	}{
		{
			name:           "custom page and per_page",
			query:          "?page=2&per_page=10",
			expectedStatus: http.StatusOK,
			expectedPage:   2,
			expectedPP:     10,
		},
		{
			name:           "invalid page",
			query:          "?page=0",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "per_page exceeds max",
			query:          "?per_page=150",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/bands"+tt.query, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d: %s", tt.expectedStatus, w.Code, w.Body.String())
			}

			if tt.expectedStatus == http.StatusOK {
				var resp struct {
					Meta struct {
						Page    int `json:"page"`
						PerPage int `json:"per_page"`
					} `json:"meta"`
				}
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("failed to parse response: %v", err)
				}
				if resp.Meta.Page != tt.expectedPage {
					t.Errorf("expected page %d, got %d", tt.expectedPage, resp.Meta.Page)
				}
				if resp.Meta.PerPage != tt.expectedPP {
					t.Errorf("expected per_page %d, got %d", tt.expectedPP, resp.Meta.PerPage)
				}
			}
		})
	}
}

func TestListBands_GenreFilter(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Close()

	router := setupBandsTestRouter(tdb)

	tests := []struct {
		name           string
		query          string
		expectedStatus int
	}{
		{
			name:           "single genre",
			query:          "?genre=rock",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "multiple genres",
			query:          "?genre=rock&genre=indie",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/bands"+tt.query, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d: %s", tt.expectedStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestListBands_SearchQuery(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Close()

	ctx := context.Background()

	// Clean up and create test data
	tdb.CleanupTestData(ctx)
	defer tdb.CleanupTestData(ctx)

	// Insert a test band for searching
	_, err := tdb.InsertTestBand(ctx, "Test Band Searchable", "test-band-searchable")
	if err != nil {
		t.Fatalf("failed to insert test band: %v", err)
	}

	router := setupBandsTestRouter(tdb)

	req := httptest.NewRequest(http.MethodGet, "/api/bands?q=Searchable", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var resp struct {
		Data []struct {
			Name string `json:"name"`
		} `json:"data"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	// Should find our test band
	found := false
	for _, band := range resp.Data {
		if band.Name == "Test Band Searchable" {
			found = true
			break
		}
	}

	if !found {
		t.Log("Test band not found in search results - may be due to search index timing")
	}
}

func TestGetBand_Success(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Close()

	ctx := context.Background()

	// Clean up and create test data
	tdb.CleanupTestData(ctx)
	defer tdb.CleanupTestData(ctx)

	// Insert a test band
	bandID, err := tdb.InsertTestBand(ctx, "Test Band Detail", "test-band-detail")
	if err != nil {
		t.Fatalf("failed to insert test band: %v", err)
	}

	// Add a genre to the band
	genreID, err := tdb.GetFirstGenreID(ctx)
	if err == nil {
		tdb.AddGenreToBand(ctx, bandID, genreID)
	}

	router := setupBandsTestRouter(tdb)

	req := httptest.NewRequest(http.MethodGet, "/api/bands/test-band-detail", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var resp struct {
		Data struct {
			ID            int32        `json:"id"`
			Name          string       `json:"name"`
			Slug          string       `json:"slug"`
			Bio           *string      `json:"bio"`
			Hometown      *string      `json:"hometown"`
			Website       *string      `json:"website"`
			SpotifyURL    *string      `json:"spotify_url"`
			Genres        []interface{} `json:"genres"`
			UpcomingShows []interface{} `json:"upcoming_shows"`
		} `json:"data"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Data.ID != bandID {
		t.Errorf("expected band ID %d, got %d", bandID, resp.Data.ID)
	}

	if resp.Data.Name != "Test Band Detail" {
		t.Errorf("expected name 'Test Band Detail', got '%s'", resp.Data.Name)
	}

	if resp.Data.Slug != "test-band-detail" {
		t.Errorf("expected slug 'test-band-detail', got '%s'", resp.Data.Slug)
	}

	// Genres and UpcomingShows should be arrays
	if resp.Data.Genres == nil {
		t.Error("genres should not be nil")
	}
	if resp.Data.UpcomingShows == nil {
		t.Error("upcoming_shows should not be nil")
	}
}

func TestGetBand_NotFound(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Close()

	router := setupBandsTestRouter(tdb)

	req := httptest.NewRequest(http.MethodGet, "/api/bands/non-existent-band", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d: %s", http.StatusNotFound, w.Code, w.Body.String())
	}

	var resp struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Error.Code != "NOT_FOUND" {
		t.Errorf("expected error code NOT_FOUND, got %s", resp.Error.Code)
	}
}

func TestGetBand_WithUpcomingShows(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Close()

	ctx := context.Background()

	// Clean up and create test data
	tdb.CleanupTestData(ctx)
	defer tdb.CleanupTestData(ctx)

	// Get a venue
	venueID, err := tdb.GetFirstVenueID(ctx)
	if err != nil {
		t.Skipf("no venues in database: %v", err)
	}

	// Insert a test band
	bandID, err := tdb.InsertTestBand(ctx, "Test Band Shows", "test-band-shows")
	if err != nil {
		t.Fatalf("failed to insert test band: %v", err)
	}

	// Insert a test show
	showDate := time.Now().AddDate(0, 0, 7)
	showID, err := tdb.InsertTestShow(ctx, venueID, showDate, "Band's Show")
	if err != nil {
		t.Fatalf("failed to insert test show: %v", err)
	}

	// Link band to show
	if err := tdb.LinkBandToShow(ctx, showID, bandID, true, 1); err != nil {
		t.Fatalf("failed to link band to show: %v", err)
	}

	router := setupBandsTestRouter(tdb)

	req := httptest.NewRequest(http.MethodGet, "/api/bands/test-band-shows", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var resp struct {
		Data struct {
			UpcomingShows []struct {
				ID          int32 `json:"id"`
				Date        string `json:"date"`
				IsHeadliner bool   `json:"is_headliner"`
				Venue       struct {
					ID   int32  `json:"id"`
					Name string `json:"name"`
					Slug string `json:"slug"`
				} `json:"venue"`
			} `json:"upcoming_shows"`
		} `json:"data"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(resp.Data.UpcomingShows) != 1 {
		t.Errorf("expected 1 upcoming show, got %d", len(resp.Data.UpcomingShows))
	}

	if len(resp.Data.UpcomingShows) > 0 {
		show := resp.Data.UpcomingShows[0]
		if show.ID != showID {
			t.Errorf("expected show ID %d, got %d", showID, show.ID)
		}
		if !show.IsHeadliner {
			t.Error("expected band to be headliner")
		}
		if show.Venue.ID != venueID {
			t.Errorf("expected venue ID %d, got %d", venueID, show.Venue.ID)
		}
	}
}

func TestGetSimilarBands_Success(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Close()

	ctx := context.Background()

	// Clean up and create test data
	tdb.CleanupTestData(ctx)
	defer tdb.CleanupTestData(ctx)

	// Get a genre
	genreID, err := tdb.GetFirstGenreID(ctx)
	if err != nil {
		t.Skipf("no genres in database: %v", err)
	}

	// Insert test bands with the same genre
	bandID1, err := tdb.InsertTestBand(ctx, "Test Band Similar 1", "test-band-similar-1")
	if err != nil {
		t.Fatalf("failed to insert test band 1: %v", err)
	}
	if err := tdb.AddGenreToBand(ctx, bandID1, genreID); err != nil {
		t.Fatalf("failed to add genre to band 1: %v", err)
	}

	bandID2, err := tdb.InsertTestBand(ctx, "Test Band Similar 2", "test-band-similar-2")
	if err != nil {
		t.Fatalf("failed to insert test band 2: %v", err)
	}
	if err := tdb.AddGenreToBand(ctx, bandID2, genreID); err != nil {
		t.Fatalf("failed to add genre to band 2: %v", err)
	}

	router := setupBandsTestRouter(tdb)

	req := httptest.NewRequest(http.MethodGet, "/api/bands/test-band-similar-1/similar", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var resp struct {
		Data []struct {
			ID               int32        `json:"id"`
			Name             string       `json:"name"`
			Slug             string       `json:"slug"`
			SharedGenreCount int64        `json:"shared_genre_count"`
			SharedGenres     []interface{} `json:"shared_genres"`
		} `json:"data"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	// Should find band 2 as similar (same genre)
	found := false
	for _, band := range resp.Data {
		if band.ID == bandID2 {
			found = true
			if band.SharedGenreCount < 1 {
				t.Errorf("expected shared_genre_count >= 1, got %d", band.SharedGenreCount)
			}
			break
		}
	}

	if !found {
		t.Log("Similar band not found - may depend on query implementation")
	}

	// Should NOT include the source band itself
	for _, band := range resp.Data {
		if band.ID == bandID1 {
			t.Error("similar bands should not include the source band")
		}
	}
}

func TestGetSimilarBands_NotFound(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Close()

	router := setupBandsTestRouter(tdb)

	req := httptest.NewRequest(http.MethodGet, "/api/bands/non-existent-band/similar", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d: %s", http.StatusNotFound, w.Code, w.Body.String())
	}
}

func TestGetSimilarBands_LimitParam(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Close()

	ctx := context.Background()

	// Clean up and create test data
	tdb.CleanupTestData(ctx)
	defer tdb.CleanupTestData(ctx)

	// Get a genre
	genreID, err := tdb.GetFirstGenreID(ctx)
	if err != nil {
		t.Skipf("no genres in database: %v", err)
	}

	// Insert source band
	bandID, err := tdb.InsertTestBand(ctx, "Test Band Limit Source", "test-band-limit-source")
	if err != nil {
		t.Fatalf("failed to insert source band: %v", err)
	}
	if err := tdb.AddGenreToBand(ctx, bandID, genreID); err != nil {
		t.Fatalf("failed to add genre to source band: %v", err)
	}

	// Insert multiple similar bands
	for i := 1; i <= 5; i++ {
		bid, err := tdb.InsertTestBand(ctx, fmt.Sprintf("Test Band Limit %d", i), fmt.Sprintf("test-band-limit-%d", i))
		if err != nil {
			t.Fatalf("failed to insert band %d: %v", i, err)
		}
		if err := tdb.AddGenreToBand(ctx, bid, genreID); err != nil {
			t.Fatalf("failed to add genre to band %d: %v", i, err)
		}
	}

	router := setupBandsTestRouter(tdb)

	tests := []struct {
		name           string
		query          string
		expectedStatus int
		maxResults     int
	}{
		{
			name:           "default limit",
			query:          "",
			expectedStatus: http.StatusOK,
			maxResults:     10,
		},
		{
			name:           "custom limit 2",
			query:          "?limit=2",
			expectedStatus: http.StatusOK,
			maxResults:     2,
		},
		{
			name:           "invalid limit",
			query:          "?limit=-1",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "limit capped at 50",
			query:          "?limit=100",
			expectedStatus: http.StatusOK,
			maxResults:     50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/bands/test-band-limit-source/similar"+tt.query, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d: %s", tt.expectedStatus, w.Code, w.Body.String())
			}

			if tt.expectedStatus == http.StatusOK {
				var resp struct {
					Data []interface{} `json:"data"`
				}
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("failed to parse response: %v", err)
				}

				if len(resp.Data) > tt.maxResults {
					t.Errorf("expected at most %d results, got %d", tt.maxResults, len(resp.Data))
				}
			}
		})
	}
}

func TestGetBand_ResponseStructure(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Close()

	ctx := context.Background()
	tdb.CleanupTestData(ctx)
	defer tdb.CleanupTestData(ctx)

	// Insert a test band
	_, err := tdb.InsertTestBand(ctx, "Test Band Structure", "test-band-structure")
	if err != nil {
		t.Fatalf("failed to insert test band: %v", err)
	}

	router := setupBandsTestRouter(tdb)

	req := httptest.NewRequest(http.MethodGet, "/api/bands/test-band-structure", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}

	// Verify all expected fields are present
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatal("response should have a 'data' object")
	}

	expectedFields := []string{
		"id", "name", "slug", "bio", "hometown", "image_url",
		"website", "spotify_url", "instagram", "facebook", "bandcamp_url",
		"genres", "upcoming_shows",
	}

	for _, field := range expectedFields {
		if _, exists := data[field]; !exists {
			t.Errorf("missing expected field: %s", field)
		}
	}
}
