package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/paulsena/asheville-setlist/internal/db"
)

// ListVenues handles GET /api/venues with optional region filter.
func (h *Handler) ListVenues(c *gin.Context) {
	ctx := c.Request.Context()

	regions := c.QueryArray("region")

	var venues []VenueListItem

	if len(regions) > 0 {
		rows, err := h.queries.ListVenuesByRegion(ctx, regions)
		if err != nil {
			slog.Error("failed to list venues by region", "error", err)
			respondInternalError(c)
			return
		}
		venues = convertRegionVenuesToListItems(rows)
	} else {
		rows, err := h.queries.ListVenuesWithShowCount(ctx)
		if err != nil {
			slog.Error("failed to list venues", "error", err)
			respondInternalError(c)
			return
		}
		venues = convertVenuesToListItems(rows)
	}

	respondJSON(c, http.StatusOK, venues)
}

// GetVenue handles GET /api/venues/:slug
func (h *Handler) GetVenue(c *gin.Context) {
	ctx := c.Request.Context()

	slug := c.Param("slug")
	if slug == "" {
		respondInvalidParam(c, "slug", "venue slug is required")
		return
	}

	venue, err := h.queries.GetVenueBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondNotFound(c, "Venue")
			return
		}
		slog.Error("failed to get venue", "slug", slug, "error", err)
		respondInternalError(c)
		return
	}

	// Get upcoming shows for this venue
	showRows, err := h.queries.GetVenueUpcomingShows(ctx, db.GetVenueUpcomingShowsParams{
		VenueID: venue.ID,
		Limit:   VenueUpcomingShowsLimit,
	})
	if err != nil {
		slog.Error("failed to get venue shows", "venue_id", venue.ID, "error", err)
		respondInternalError(c)
		return
	}

	// Batch load bands for all shows
	showIDs := make([]int32, len(showRows))
	for i, s := range showRows {
		showIDs[i] = s.ID
	}

	bandsMap := make(map[int32][]BandBasic)
	if len(showIDs) > 0 {
		bandRows, err := h.queries.GetShowBandsForVenue(ctx, showIDs)
		if err != nil {
			slog.Error("failed to get bands for venue shows", "error", err)
		} else {
			for _, b := range bandRows {
				bandsMap[b.ShowID] = append(bandsMap[b.ShowID], BandBasic{
					ID:          b.ID,
					Name:        b.Name,
					Slug:        b.Slug,
					IsHeadliner: boolValue(b.IsHeadliner),
				})
			}
		}
	}

	// Build upcoming shows list
	upcomingShows := make([]VenueShowItem, len(showRows))
	for i, s := range showRows {
		bands := bandsMap[s.ID]
		if bands == nil {
			bands = []BandBasic{}
		}
		upcomingShows[i] = VenueShowItem{
			ID:       s.ID,
			Title:    s.Title,
			Date:     formatTimestamp(s.Date),
			PriceMin: numericToFloat(s.PriceMin),
			PriceMax: numericToFloat(s.PriceMax),
			Bands:    bands,
		}
	}

	detail := VenueDetail{
		ID:            venue.ID,
		Name:          venue.Name,
		Slug:          venue.Slug,
		Address:       venue.Address,
		City:          stringValue(venue.City),
		State:         stringValue(venue.State),
		ZipCode:       venue.ZipCode,
		Region:        venue.Region,
		Capacity:      venue.Capacity,
		Website:       venue.Website,
		Phone:         venue.Phone,
		ImageURL:      venue.ImageUrl,
		UpcomingShows: upcomingShows,
	}

	respondJSON(c, http.StatusOK, detail)
}
