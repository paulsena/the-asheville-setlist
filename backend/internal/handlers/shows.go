package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/paulsena/asheville-setlist/internal/db"
)

// ListShows handles GET /api/shows with various filter options.
func (h *Handler) ListShows(c *gin.Context) {
	ctx := c.Request.Context()

	page, perPage, err := parsePagination(c)
	if err != nil {
		if pe, ok := err.(*paramError); ok {
			respondInvalidParam(c, pe.param, pe.message)
			return
		}
		respondInternalError(c)
		return
	}

	offset := calculateOffset(page, perPage)

	// Route to appropriate filter handler
	filter := c.Query("filter")

	var shows []ShowListItem
	var total int

	switch filter {
	case "tonight":
		rows, err := h.queries.ListShowsTonight(ctx)
		if err != nil {
			slog.Error("failed to list shows tonight", "error", err)
			respondInternalError(c)
			return
		}
		shows = convertTonightShowsToListItems(rows)
		total = len(shows)

	case "this-weekend":
		rows, err := h.queries.ListShowsThisWeekend(ctx)
		if err != nil {
			slog.Error("failed to list shows this weekend", "error", err)
			respondInternalError(c)
			return
		}
		shows = convertWeekendShowsToListItems(rows)
		total = len(shows)

	case "free":
		rows, err := h.queries.ListFreeShows(ctx, db.ListFreeShowsParams{
			Limit:  int32(perPage),
			Offset: int32(offset),
		})
		if err != nil {
			slog.Error("failed to list free shows", "error", err)
			respondInternalError(c)
			return
		}
		shows, total = convertFreeShowsToListItems(rows)

	default:
		shows, total, err = h.listShowsWithFilters(ctx, c, perPage, offset)
		if err != nil {
			return // Error response already sent
		}
	}

	// Load bands for all shows in batch
	if len(shows) > 0 {
		h.attachBandsToShows(ctx, shows)
	}

	meta := &Meta{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: calculateTotalPages(total, perPage),
	}

	respondJSONWithMeta(c, http.StatusOK, shows, meta)
}

// listShowsWithFilters handles venue, region, genre, and date range filters.
func (h *Handler) listShowsWithFilters(ctx context.Context, c *gin.Context, perPage, offset int) ([]ShowListItem, int, error) {
	venues := c.QueryArray("venue")
	regions := c.QueryArray("region")
	genres := c.QueryArray("genre")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	if len(venues) > 0 {
		rows, err := h.queries.ListShowsByVenue(ctx, db.ListShowsByVenueParams{
			Column1: venues,
			Limit:   int32(perPage),
			Offset:  int32(offset),
		})
		if err != nil {
			slog.Error("failed to list shows by venue", "error", err)
			respondInternalError(c)
			return nil, 0, err
		}
		shows, total := convertVenueShowsToListItems(rows)
		return shows, total, nil
	}

	if len(regions) > 0 {
		rows, err := h.queries.ListShowsByRegion(ctx, db.ListShowsByRegionParams{
			Column1: regions,
			Limit:   int32(perPage),
			Offset:  int32(offset),
		})
		if err != nil {
			slog.Error("failed to list shows by region", "error", err)
			respondInternalError(c)
			return nil, 0, err
		}
		shows, total := convertRegionShowsToListItems(rows)
		return shows, total, nil
	}

	if len(genres) > 0 {
		rows, err := h.queries.ListShowsByGenre(ctx, db.ListShowsByGenreParams{
			Column1: genres,
			Limit:   int32(perPage),
			Offset:  int32(offset),
		})
		if err != nil {
			slog.Error("failed to list shows by genre", "error", err)
			respondInternalError(c)
			return nil, 0, err
		}
		count, err := h.queries.CountShowsByGenre(ctx, genres)
		if err != nil {
			slog.Error("failed to count shows by genre", "error", err)
			respondInternalError(c)
			return nil, 0, err
		}
		return convertGenreShowsToListItems(rows), int(count), nil
	}

	if dateFrom != "" || dateTo != "" {
		fromTime, toTime, err := parseDateRange(dateFrom, dateTo)
		if err != nil {
			respondInvalidParam(c, "date_from/date_to", "invalid date format, use ISO 8601")
			return nil, 0, err
		}
		rows, err := h.queries.ListShowsByDateRange(ctx, db.ListShowsByDateRangeParams{
			Date:   fromTime,
			Date_2: toTime,
			Limit:  int32(perPage),
			Offset: int32(offset),
		})
		if err != nil {
			slog.Error("failed to list shows by date range", "error", err)
			respondInternalError(c)
			return nil, 0, err
		}
		shows, total := convertDateRangeShowsToListItems(rows)
		return shows, total, nil
	}

	// Default: upcoming shows
	rows, err := h.queries.ListUpcomingShows(ctx, db.ListUpcomingShowsParams{
		Limit:  int32(perPage),
		Offset: int32(offset),
	})
	if err != nil {
		slog.Error("failed to list upcoming shows", "error", err)
		respondInternalError(c)
		return nil, 0, err
	}
	shows, total := convertUpcomingShowsToListItems(rows)
	return shows, total, nil
}

// attachBandsToShows loads and attaches bands to show list items.
func (h *Handler) attachBandsToShows(ctx context.Context, shows []ShowListItem) {
	showIDs := make([]int32, len(shows))
	for i, s := range shows {
		showIDs[i] = s.ID
	}

	bandsMap, err := h.loadBandsForShows(ctx, showIDs)
	if err != nil {
		slog.Error("failed to load bands for shows", "error", err)
		return // Continue without bands rather than failing
	}

	for i := range shows {
		if bands, ok := bandsMap[shows[i].ID]; ok {
			shows[i].Bands = bands
		}
	}
}

// GetShow handles GET /api/shows/:id
func (h *Handler) GetShow(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		respondInvalidParam(c, "id", "must be a valid integer")
		return
	}

	show, err := h.queries.GetShowByID(ctx, int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondNotFound(c, "Show")
			return
		}
		slog.Error("failed to get show", "id", id, "error", err)
		respondInternalError(c)
		return
	}

	bandRows, err := h.queries.GetShowBands(ctx, int32(id))
	if err != nil {
		slog.Error("failed to get show bands", "show_id", id, "error", err)
		respondInternalError(c)
		return
	}

	// Load genres for each band
	bands := make([]BandForShow, len(bandRows))
	for i, b := range bandRows {
		genres, err := h.queries.GetBandGenresForShow(ctx, b.ID)
		if err != nil {
			slog.Error("failed to get band genres", "band_id", b.ID, "error", err)
			genres = []db.GetBandGenresForShowRow{}
		}

		bands[i] = BandForShow{
			ID:               b.ID,
			Name:             b.Name,
			Slug:             b.Slug,
			Bio:              b.Bio,
			ImageURL:         b.ImageUrl,
			SpotifyURL:       b.SpotifyUrl,
			Website:          b.Website,
			IsHeadliner:      boolValue(b.IsHeadliner),
			PerformanceOrder: int32Value(b.PerformanceOrder),
			Genres:           convertGenreRows(genres),
		}
	}

	detail := ShowDetail{
		ID:             show.ID,
		Title:          show.Title,
		Description:    show.Description,
		ImageURL:       show.ImageUrl,
		Date:           formatTimestamp(show.Date),
		DoorsTime:      formatTime(show.DoorsTime),
		ShowTime:       formatTime(show.ShowTime),
		PriceMin:       numericToFloat(show.PriceMin),
		PriceMax:       numericToFloat(show.PriceMax),
		TicketURL:      show.TicketUrl,
		AgeRestriction: show.AgeRestriction,
		Status:         stringValue(show.Status),
		Venue: VenueForShow{
			ID:       show.VenueID,
			Name:     show.VenueName,
			Slug:     show.VenueSlug,
			Address:  show.VenueAddress,
			Region:   show.VenueRegion,
			Website:  show.VenueWebsite,
			ImageURL: show.VenueImageUrl,
		},
		Bands: bands,
	}

	respondJSON(c, http.StatusOK, detail)
}

// loadBandsForShows loads bands for multiple shows.
// TODO: Optimize with batch query using GetShowBandsForVenue once available.
func (h *Handler) loadBandsForShows(ctx context.Context, showIDs []int32) (map[int32][]BandBasic, error) {
	result := make(map[int32][]BandBasic)

	for _, id := range showIDs {
		result[id] = []BandBasic{}
	}

	for _, showID := range showIDs {
		bands, err := h.queries.GetShowBands(ctx, showID)
		if err != nil {
			return nil, err
		}
		for _, b := range bands {
			result[showID] = append(result[showID], BandBasic{
				ID:               b.ID,
				Name:             b.Name,
				Slug:             b.Slug,
				ImageURL:         b.ImageUrl,
				IsHeadliner:      boolValue(b.IsHeadliner),
				PerformanceOrder: int32Value(b.PerformanceOrder),
			})
		}
	}

	return result, nil
}
