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

// ListBands handles GET /api/bands with optional genre filter and search.
func (h *Handler) ListBands(c *gin.Context) {
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
	genres := c.QueryArray("genre")
	query := c.Query("q")

	var bands []BandListItem
	var total int

	if query != "" {
		rows, err := h.queries.SearchBands(ctx, db.SearchBandsParams{
			PlaintoTsquery: query,
			Limit:          int32(perPage),
			Offset:         int32(offset),
		})
		if err != nil {
			slog.Error("failed to search bands", "error", err)
			respondInternalError(c)
			return
		}
		bands, total = convertSearchBandsToListItems(rows)

	} else if len(genres) > 0 {
		rows, err := h.queries.ListBandsByGenre(ctx, db.ListBandsByGenreParams{
			Column1: genres,
			Limit:   int32(perPage),
			Offset:  int32(offset),
		})
		if err != nil {
			slog.Error("failed to list bands by genre", "error", err)
			respondInternalError(c)
			return
		}

		count, err := h.queries.CountBandsByGenre(ctx, genres)
		if err != nil {
			slog.Error("failed to count bands by genre", "error", err)
			respondInternalError(c)
			return
		}

		bands = convertGenreBandsToListItems(rows)
		total = int(count)

	} else {
		rows, err := h.queries.ListBands(ctx, db.ListBandsParams{
			Limit:  int32(perPage),
			Offset: int32(offset),
		})
		if err != nil {
			slog.Error("failed to list bands", "error", err)
			respondInternalError(c)
			return
		}
		bands, total = convertBandsToListItems(rows)
	}

	// Load genres for each band
	if len(bands) > 0 {
		h.attachGenresToBands(ctx, bands)
	}

	meta := &Meta{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: calculateTotalPages(total, perPage),
	}

	respondJSONWithMeta(c, http.StatusOK, bands, meta)
}

// attachGenresToBands loads and attaches genres to band list items.
func (h *Handler) attachGenresToBands(ctx context.Context, bands []BandListItem) {
	bandIDs := make([]int32, len(bands))
	for i, b := range bands {
		bandIDs[i] = b.ID
	}

	genresMap, err := h.loadGenresForBands(ctx, bandIDs)
	if err != nil {
		slog.Error("failed to load genres for bands", "error", err)
		return
	}

	for i := range bands {
		if g, ok := genresMap[bands[i].ID]; ok {
			bands[i].Genres = g
		}
	}
}

// GetBand handles GET /api/bands/:slug
func (h *Handler) GetBand(c *gin.Context) {
	ctx := c.Request.Context()

	slug := c.Param("slug")
	if slug == "" {
		respondInvalidParam(c, "slug", "band slug is required")
		return
	}

	band, err := h.queries.GetBandBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondNotFound(c, "Band")
			return
		}
		slog.Error("failed to get band", "slug", slug, "error", err)
		respondInternalError(c)
		return
	}

	// Get genres for this band
	genreRows, err := h.queries.GetBandGenres(ctx, band.ID)
	if err != nil {
		slog.Error("failed to get band genres", "band_id", band.ID, "error", err)
		genreRows = []db.GetBandGenresRow{}
	}

	genres := make([]GenreBasic, len(genreRows))
	for i, g := range genreRows {
		genres[i] = GenreBasic{
			ID:   g.ID,
			Name: g.Name,
			Slug: g.Slug,
		}
	}

	// Get upcoming shows for this band
	showRows, err := h.queries.GetBandUpcomingShows(ctx, band.ID)
	if err != nil {
		slog.Error("failed to get band shows", "band_id", band.ID, "error", err)
		respondInternalError(c)
		return
	}

	upcomingShows := make([]BandShowItem, len(showRows))
	for i, s := range showRows {
		upcomingShows[i] = BandShowItem{
			ID:   s.ID,
			Date: formatTimestamp(s.Date),
			Venue: VenueBasic{
				ID:   s.VenueID,
				Name: s.VenueName,
				Slug: s.VenueSlug,
			},
			IsHeadliner: boolValue(s.IsHeadliner),
		}
	}

	detail := BandDetail{
		ID:            band.ID,
		Name:          band.Name,
		Slug:          band.Slug,
		Bio:           band.Bio,
		Hometown:      band.Hometown,
		ImageURL:      band.ImageUrl,
		Website:       band.Website,
		SpotifyURL:    band.SpotifyUrl,
		Instagram:     band.Instagram,
		Facebook:      band.Facebook,
		BandcampURL:   band.BandcampUrl,
		Genres:        genres,
		UpcomingShows: upcomingShows,
	}

	respondJSON(c, http.StatusOK, detail)
}

// GetSimilarBands handles GET /api/bands/:slug/similar
func (h *Handler) GetSimilarBands(c *gin.Context) {
	ctx := c.Request.Context()

	slug := c.Param("slug")
	if slug == "" {
		respondInvalidParam(c, "slug", "band slug is required")
		return
	}

	limit := DefaultSimilarBandsLimit
	if l := c.Query("limit"); l != "" {
		parsed, err := strconv.Atoi(l)
		if err != nil || parsed < 1 {
			respondInvalidParam(c, "limit", "must be a positive integer")
			return
		}
		if parsed > MaxSimilarBandsLimit {
			parsed = MaxSimilarBandsLimit
		}
		limit = parsed
	}

	band, err := h.queries.GetBandBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondNotFound(c, "Band")
			return
		}
		slog.Error("failed to get band", "slug", slug, "error", err)
		respondInternalError(c)
		return
	}

	// Get source band's genres for comparison
	sourceBandGenres, err := h.queries.GetBandGenres(ctx, band.ID)
	if err != nil {
		slog.Error("failed to get source band genres", "band_id", band.ID, "error", err)
		respondInternalError(c)
		return
	}

	sourceGenreIDs := make(map[int32]db.GetBandGenresRow)
	for _, g := range sourceBandGenres {
		sourceGenreIDs[g.ID] = g
	}

	rows, err := h.queries.GetSimilarBands(ctx, db.GetSimilarBandsParams{
		BandID: band.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		slog.Error("failed to get similar bands", "band_id", band.ID, "error", err)
		respondInternalError(c)
		return
	}

	// Load genres for similar bands
	bandIDs := make([]int32, len(rows))
	for i, r := range rows {
		bandIDs[i] = r.ID
	}

	genresMap := make(map[int32][]GenreBasic)
	if len(bandIDs) > 0 {
		genresMap, err = h.loadGenresForBands(ctx, bandIDs)
		if err != nil {
			slog.Error("failed to load genres for similar bands", "error", err)
		}
	}

	similar := make([]SimilarBandItem, len(rows))
	for i, r := range rows {
		// Find shared genres
		bandGenres := genresMap[r.ID]
		sharedGenres := []GenreBasic{}
		for _, g := range bandGenres {
			if _, exists := sourceGenreIDs[g.ID]; exists {
				sharedGenres = append(sharedGenres, g)
			}
		}

		similar[i] = SimilarBandItem{
			ID:               r.ID,
			Name:             r.Name,
			Slug:             r.Slug,
			ImageURL:         r.ImageUrl,
			SharedGenreCount: r.SharedGenreCount,
			SharedGenres:     sharedGenres,
		}
	}

	respondJSON(c, http.StatusOK, similar)
}

// loadGenresForBands loads genres for multiple bands using batch query.
func (h *Handler) loadGenresForBands(ctx context.Context, bandIDs []int32) (map[int32][]GenreBasic, error) {
	result := make(map[int32][]GenreBasic)

	for _, id := range bandIDs {
		result[id] = []GenreBasic{}
	}

	rows, err := h.queries.GetBandGenresBatch(ctx, bandIDs)
	if err != nil {
		return nil, err
	}

	for _, r := range rows {
		result[r.BandID] = append(result[r.BandID], GenreBasic{
			ID:   r.ID,
			Name: r.Name,
			Slug: r.Slug,
		})
	}

	return result, nil
}
