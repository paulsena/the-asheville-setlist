package handlers_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/paulsena/asheville-setlist/internal/handlers"
	"github.com/paulsena/asheville-setlist/internal/testutil"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// setupShowsTestRouter creates a test router with the shows handler
func setupShowsTestRouter(tdb *testutil.TestDB) *gin.Engine {
	h := handlers.New(tdb.Queries)
	router := gin.New()
	router.GET("/api/shows", h.ListShows)
	router.GET("/api/shows/:id", h.GetShow)
	router.POST("/api/shows", h.CreateShow)
	return router
}

func TestListShows_Default(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Close()

	router := setupShowsTestRouter(tdb)

	req := httptest.NewRequest(http.MethodGet, "/api/shows", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var resp struct {
		Data []interface{} `json:"data"`
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
}

func TestListShows_Pagination(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Close()

	router := setupShowsTestRouter(tdb)

	tests := []struct {
		name           string
		query          string
		expectedStatus int
		expectedPage   int
		expectedPP     int
	}{
		{
			name:           "custom page",
			query:          "?page=2&per_page=10",
			expectedStatus: http.StatusOK,
			expectedPage:   2,
			expectedPP:     10,
		},
		{
			name:           "invalid page",
			query:          "?page=-1",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid per_page",
			query:          "?per_page=abc",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "per_page exceeds max",
			query:          "?per_page=200",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/shows"+tt.query, nil)
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

func TestListShows_VenueFilter(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Close()

	router := setupShowsTestRouter(tdb)

	// Test with a venue slug that exists (from seed data)
	req := httptest.NewRequest(http.MethodGet, "/api/shows?venue=the-orange-peel", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var resp struct {
		Data []interface{} `json:"data"`
		Meta struct {
			Total int `json:"total"`
		} `json:"meta"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	// Response should be valid (may or may not have shows)
	if resp.Data == nil {
		t.Error("data should not be nil")
	}
}

func TestListShows_RegionFilter(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Close()

	router := setupShowsTestRouter(tdb)

	req := httptest.NewRequest(http.MethodGet, "/api/shows?region=downtown", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}
}

func TestListShows_DateRangeFilter(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Close()

	router := setupShowsTestRouter(tdb)

	tests := []struct {
		name           string
		query          string
		expectedStatus int
	}{
		{
			name:           "valid date range",
			query:          "?date_from=2025-01-01&date_to=2025-12-31",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "only date_from",
			query:          "?date_from=2025-01-01",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid date format",
			query:          "?date_from=invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/shows"+tt.query, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d: %s", tt.expectedStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestListShows_SpecialFilters(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Close()

	router := setupShowsTestRouter(tdb)

	filters := []string{"tonight", "this-weekend", "free"}

	for _, filter := range filters {
		t.Run(filter, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/shows?filter="+filter, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
			}

			var resp struct {
				Data []interface{} `json:"data"`
			}
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("failed to parse response: %v", err)
			}

			if resp.Data == nil {
				t.Error("data should not be nil")
			}
		})
	}
}

func TestGetShow_Success(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Close()

	ctx := context.Background()

	// Clean up before and after
	tdb.CleanupTestData(ctx)
	defer tdb.CleanupTestData(ctx)

	// Get a venue to use
	venueID, err := tdb.GetFirstVenueID(ctx)
	if err != nil {
		t.Skipf("no venues in database: %v", err)
	}

	// Insert a test show
	showDate := time.Now().AddDate(0, 0, 7) // 1 week from now
	showID, err := tdb.InsertTestShow(ctx, venueID, showDate, "Test Show Title")
	if err != nil {
		t.Fatalf("failed to insert test show: %v", err)
	}

	router := setupShowsTestRouter(tdb)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/shows/%d", showID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var resp struct {
		Data struct {
			ID    int32  `json:"id"`
			Title string `json:"title"`
			Venue struct {
				ID int32 `json:"id"`
			} `json:"venue"`
			Bands []interface{} `json:"bands"`
		} `json:"data"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Data.ID != showID {
		t.Errorf("expected show ID %d, got %d", showID, resp.Data.ID)
	}

	if resp.Data.Venue.ID != venueID {
		t.Errorf("expected venue ID %d, got %d", venueID, resp.Data.Venue.ID)
	}
}

func TestGetShow_NotFound(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Close()

	router := setupShowsTestRouter(tdb)

	req := httptest.NewRequest(http.MethodGet, "/api/shows/999999", nil)
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

func TestGetShow_InvalidID(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Close()

	router := setupShowsTestRouter(tdb)

	req := httptest.NewRequest(http.MethodGet, "/api/shows/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d: %s", http.StatusBadRequest, w.Code, w.Body.String())
	}
}

func TestCreateShow_Success(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Close()

	ctx := context.Background()

	// Clean up before and after
	tdb.CleanupTestData(ctx)
	defer tdb.CleanupTestData(ctx)

	// Get a venue to use
	venueID, err := tdb.GetFirstVenueID(ctx)
	if err != nil {
		t.Skipf("no venues in database: %v", err)
	}

	router := setupShowsTestRouter(tdb)

	futureDate := time.Now().AddDate(0, 1, 0).Format("2006-01-02")
	body := fmt.Sprintf(`{
		"venue_id": %d,
		"date": "%s",
		"bands": [
			{"name": "Test Band Headliner", "is_headliner": true, "performance_order": 1},
			{"name": "Test Band Opener", "is_headliner": false, "performance_order": 2}
		],
		"price_min": 20.00,
		"price_max": 30.00,
		"age_restriction": "21+"
	}`, venueID, futureDate)

	req := httptest.NewRequest(http.MethodPost, "/api/shows", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d: %s", http.StatusCreated, w.Code, w.Body.String())
	}

	var resp struct {
		Data struct {
			ID     int32  `json:"id"`
			Status string `json:"status"`
		} `json:"data"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Data.ID == 0 {
		t.Error("expected non-zero show ID")
	}

	if resp.Data.Status != "scheduled" {
		t.Errorf("expected status 'scheduled', got '%s'", resp.Data.Status)
	}
}

func TestCreateShow_ValidationErrors(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Close()

	ctx := context.Background()

	venueID, err := tdb.GetFirstVenueID(ctx)
	if err != nil {
		t.Skipf("no venues in database: %v", err)
	}

	router := setupShowsTestRouter(tdb)

	futureDate := time.Now().AddDate(0, 1, 0).Format("2006-01-02")

	tests := []struct {
		name           string
		body           string
		expectedStatus int
	}{
		{
			name:           "missing venue_id",
			body:           `{"date": "2025-06-01", "bands": [{"name": "Test Band"}]}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing date",
			body:           fmt.Sprintf(`{"venue_id": %d, "bands": [{"name": "Test Band"}]}`, venueID),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing bands",
			body:           fmt.Sprintf(`{"venue_id": %d, "date": "%s"}`, venueID, futureDate),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty bands array",
			body:           fmt.Sprintf(`{"venue_id": %d, "date": "%s", "bands": []}`, venueID, futureDate),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "past date",
			body:           fmt.Sprintf(`{"venue_id": %d, "date": "2020-01-01", "bands": [{"name": "Test Band"}]}`, venueID),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid age restriction",
			body:           fmt.Sprintf(`{"venue_id": %d, "date": "%s", "bands": [{"name": "Test Band"}], "age_restriction": "invalid"}`, venueID, futureDate),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "negative price",
			body:           fmt.Sprintf(`{"venue_id": %d, "date": "%s", "bands": [{"name": "Test Band"}], "price_min": -10}`, venueID, futureDate),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "price_max less than price_min",
			body:           fmt.Sprintf(`{"venue_id": %d, "date": "%s", "bands": [{"name": "Test Band"}], "price_min": 30, "price_max": 20}`, venueID, futureDate),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "non-existent venue",
			body:           fmt.Sprintf(`{"venue_id": 999999, "date": "%s", "bands": [{"name": "Test Band"}]}`, futureDate),
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/shows", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d: %s", tt.expectedStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestGetShow_WithBands(t *testing.T) {
	tdb := testutil.SetupTestDB(t)
	defer tdb.Close()

	ctx := context.Background()

	// Clean up before and after
	tdb.CleanupTestData(ctx)
	defer tdb.CleanupTestData(ctx)

	// Get a venue to use
	venueID, err := tdb.GetFirstVenueID(ctx)
	if err != nil {
		t.Skipf("no venues in database: %v", err)
	}

	// Insert a test show
	showDate := time.Now().AddDate(0, 0, 7)
	showID, err := tdb.InsertTestShow(ctx, venueID, showDate, "Test Show With Bands")
	if err != nil {
		t.Fatalf("failed to insert test show: %v", err)
	}

	// Insert test bands
	bandID1, err := tdb.InsertTestBand(ctx, "Test Band Headliner", "test-band-headliner")
	if err != nil {
		t.Fatalf("failed to insert test band 1: %v", err)
	}

	bandID2, err := tdb.InsertTestBand(ctx, "Test Band Opener", "test-band-opener")
	if err != nil {
		t.Fatalf("failed to insert test band 2: %v", err)
	}

	// Link bands to show
	if err := tdb.LinkBandToShow(ctx, showID, bandID1, true, 1); err != nil {
		t.Fatalf("failed to link band 1: %v", err)
	}
	if err := tdb.LinkBandToShow(ctx, showID, bandID2, false, 2); err != nil {
		t.Fatalf("failed to link band 2: %v", err)
	}

	router := setupShowsTestRouter(tdb)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/shows/%d", showID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var resp struct {
		Data struct {
			ID    int32 `json:"id"`
			Bands []struct {
				ID               int32  `json:"id"`
				Name             string `json:"name"`
				Slug             string `json:"slug"`
				IsHeadliner      bool   `json:"is_headliner"`
				PerformanceOrder int    `json:"performance_order"`
			} `json:"bands"`
		} `json:"data"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(resp.Data.Bands) != 2 {
		t.Errorf("expected 2 bands, got %d", len(resp.Data.Bands))
	}

	// Check that one band is the headliner
	hasHeadliner := false
	for _, band := range resp.Data.Bands {
		if band.IsHeadliner {
			hasHeadliner = true
			break
		}
	}
	if !hasHeadliner {
		t.Error("expected at least one band to be headliner")
	}
}
