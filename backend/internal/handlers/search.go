package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/paulsena/asheville-setlist/internal/db"
)

// Search handles GET /api/search for global search across shows, bands, and venues.
func (h *Handler) Search(c *gin.Context) {
	ctx := c.Request.Context()

	query := c.Query("q")
	if query == "" {
		respondMissingParam(c, "q")
		return
	}

	if len(query) < 2 {
		respondInvalidParam(c, "q", "must be at least 2 characters")
		return
	}

	limit := DefaultSearchLimit
	if l := c.Query("limit"); l != "" {
		parsed, err := strconv.Atoi(l)
		if err != nil || parsed < 1 {
			respondInvalidParam(c, "limit", "must be a positive integer")
			return
		}
		if parsed > MaxSearchLimit {
			parsed = MaxSearchLimit
		}
		limit = parsed
	}

	// Search shows
	showRows, err := h.queries.GlobalSearchShows(ctx, db.GlobalSearchShowsParams{
		PlaintoTsquery: query,
		Limit:          int32(limit),
	})
	if err != nil {
		slog.Error("failed to search shows", "error", err)
		showRows = []db.GlobalSearchShowsRow{}
	}

	shows := make([]SearchShowItem, len(showRows))
	for i, r := range showRows {
		shows[i] = SearchShowItem{
			ID:        r.ID,
			Title:     r.Title,
			Date:      formatTimestamp(r.Date),
			VenueName: r.VenueName,
		}
	}

	// Search bands
	bandRows, err := h.queries.GlobalSearchBands(ctx, db.GlobalSearchBandsParams{
		PlaintoTsquery: query,
		Limit:          int32(limit),
	})
	if err != nil {
		slog.Error("failed to search bands", "error", err)
		bandRows = []db.GlobalSearchBandsRow{}
	}

	bands := make([]SearchBandItem, len(bandRows))
	for i, r := range bandRows {
		bands[i] = SearchBandItem{
			ID:   r.ID,
			Name: r.Name,
			Slug: r.Slug,
		}
	}

	// Search venues
	venueRows, err := h.queries.GlobalSearchVenues(ctx, db.GlobalSearchVenuesParams{
		PlaintoTsquery: query,
		Limit:          int32(limit),
	})
	if err != nil {
		slog.Error("failed to search venues", "error", err)
		venueRows = []db.GlobalSearchVenuesRow{}
	}

	venues := make([]SearchVenueItem, len(venueRows))
	for i, r := range venueRows {
		venues[i] = SearchVenueItem{
			ID:   r.ID,
			Name: r.Name,
			Slug: r.Slug,
		}
	}

	result := SearchResult{
		Shows:  shows,
		Bands:  bands,
		Venues: venues,
	}

	respondJSON(c, http.StatusOK, result)
}
