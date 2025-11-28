package handlers

import (
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/paulsena/asheville-setlist/internal/db"
)

// CreateShow handles POST /api/shows for band submissions.
func (h *Handler) CreateShow(c *gin.Context) {
	ctx := c.Request.Context()

	var req CreateShowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondValidationError(c, "Invalid request body", map[string]any{
			"error": err.Error(),
		})
		return
	}

	// Validate venue exists
	exists, err := h.queries.VenueExists(ctx, req.VenueID)
	if err != nil {
		slog.Error("failed to check venue exists", "error", err)
		respondInternalError(c)
		return
	}
	if !exists {
		respondNotFound(c, "Venue")
		return
	}

	// Parse and validate date
	showDate, err := parseShowDate(req.Date)
	if err != nil {
		respondValidationError(c, "Invalid date format", map[string]any{
			"date": "must be valid ISO 8601 date",
		})
		return
	}

	if showDate.Time.Before(time.Now()) {
		respondValidationError(c, "Invalid date", map[string]any{
			"date": "must be a future date",
		})
		return
	}

	// Validate price range
	if req.PriceMin != nil && *req.PriceMin < 0 {
		respondValidationError(c, "Invalid price", map[string]any{
			"price_min": "must be >= 0",
		})
		return
	}
	if req.PriceMax != nil && *req.PriceMax < 0 {
		respondValidationError(c, "Invalid price", map[string]any{
			"price_max": "must be >= 0",
		})
		return
	}
	if req.PriceMin != nil && req.PriceMax != nil && *req.PriceMax < *req.PriceMin {
		respondValidationError(c, "Invalid price range", map[string]any{
			"price_max": "must be >= price_min",
		})
		return
	}

	// Validate age restriction
	if req.AgeRestriction != nil {
		validAges := map[string]bool{"All Ages": true, "18+": true, "21+": true}
		if !validAges[*req.AgeRestriction] {
			respondValidationError(c, "Invalid age restriction", map[string]any{
				"age_restriction": "must be one of: All Ages, 18+, 21+",
			})
			return
		}
	}

	// Validate bands
	for i, band := range req.Bands {
		if strings.TrimSpace(band.Name) == "" {
			respondValidationError(c, "Invalid band", map[string]any{
				"bands": map[string]any{
					"index": i,
					"name":  "must not be empty",
				},
			})
			return
		}
	}

	doorsTime := parseTimeString(req.DoorsTime)
	showTime := parseTimeString(req.ShowTime)
	priceMin := floatToNumeric(req.PriceMin)
	priceMax := floatToNumeric(req.PriceMax)

	// Band-submitted shows start as scheduled
	status := "scheduled"
	source := "band_submitted"

	showRow, err := h.queries.CreateShow(ctx, db.CreateShowParams{
		VenueID:        req.VenueID,
		Title:          nil, // Title derived from bands
		ImageUrl:       req.ImageURL,
		Date:           showDate,
		DoorsTime:      doorsTime,
		ShowTime:       showTime,
		PriceMin:       priceMin,
		PriceMax:       priceMax,
		TicketUrl:      req.TicketURL,
		AgeRestriction: req.AgeRestriction,
		Status:         &status,
		Source:         &source,
	})
	if err != nil {
		slog.Error("failed to create show", "error", err)
		respondInternalError(c)
		return
	}

	// Process bands
	for _, bandReq := range req.Bands {
		bandName := strings.TrimSpace(bandReq.Name)

		existingBand, err := h.queries.GetBandByName(ctx, bandName)
		var bandID int32

		if err != nil {
			// Create new band
			slug := generateSlug(bandName)
			newBand, err := h.queries.CreateBand(ctx, db.CreateBandParams{
				Name: bandName,
				Slug: slug,
			})
			if err != nil {
				slog.Error("failed to create band", "name", bandName, "error", err)
				continue
			}
			bandID = newBand.ID
		} else {
			bandID = existingBand.ID
		}

		err = h.queries.CreateShowBand(ctx, db.CreateShowBandParams{
			ShowID:           showRow.ID,
			BandID:           bandID,
			IsHeadliner:      bandReq.IsHeadliner,
			PerformanceOrder: bandReq.PerformanceOrder,
		})
		if err != nil {
			slog.Error("failed to link band to show", "show_id", showRow.ID, "band_id", bandID, "error", err)
		}
	}

	response := CreateShowResponse{
		ID:        showRow.ID,
		Status:    stringValue(showRow.Status),
		CreatedAt: formatTimestamp(showRow.CreatedAt),
	}

	respondJSON(c, http.StatusCreated, response)
}

// parseShowDate parses a date string to pgtype.Timestamptz.
func parseShowDate(dateStr string) (pgtype.Timestamptz, error) {
	var result pgtype.Timestamptz

	// Try ISO 8601 with timezone
	t, err := time.Parse(time.RFC3339, dateStr)
	if err == nil {
		result.Time = t
		result.Valid = true
		return result, nil
	}

	// Try date only (default to 8 PM Eastern)
	t, err = time.Parse("2006-01-02", dateStr)
	if err == nil {
		loc, err := time.LoadLocation("America/New_York")
		if err != nil {
			loc = time.FixedZone("EST", -5*60*60)
		}
		t = time.Date(t.Year(), t.Month(), t.Day(), 20, 0, 0, 0, loc)
		result.Time = t
		result.Valid = true
		return result, nil
	}

	return result, err
}

// parseTimeString parses a time string (HH:MM or HH:MM:SS) to pgtype.Time.
func parseTimeString(timeStr *string) pgtype.Time {
	var result pgtype.Time
	if timeStr == nil || *timeStr == "" {
		return result
	}

	t, err := time.Parse("15:04", *timeStr)
	if err == nil {
		result.Microseconds = int64(t.Hour())*3600000000 + int64(t.Minute())*60000000
		result.Valid = true
		return result
	}

	t, err = time.Parse("15:04:05", *timeStr)
	if err == nil {
		result.Microseconds = int64(t.Hour())*3600000000 + int64(t.Minute())*60000000 + int64(t.Second())*1000000
		result.Valid = true
		return result
	}

	return result
}

// floatToNumeric converts a *float64 to pgtype.Numeric.
func floatToNumeric(f *float64) pgtype.Numeric {
	var result pgtype.Numeric
	if f == nil {
		return result
	}

	result.Valid = true
	if err := result.Scan(*f); err != nil {
		return pgtype.Numeric{}
	}
	return result
}

// generateSlug creates a URL-friendly slug from a string.
func generateSlug(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")

	reg := regexp.MustCompile(`[^a-z0-9-]`)
	s = reg.ReplaceAllString(s, "")

	reg = regexp.MustCompile(`-+`)
	s = reg.ReplaceAllString(s, "-")

	s = strings.Trim(s, "-")

	if s == "" {
		return "band"
	}
	return s
}
