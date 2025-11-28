package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListGenres handles GET /api/genres with show counts.
func (h *Handler) ListGenres(c *gin.Context) {
	ctx := c.Request.Context()

	rows, err := h.queries.ListGenresWithShowCount(ctx)
	if err != nil {
		slog.Error("failed to list genres", "error", err)
		respondInternalError(c)
		return
	}

	genres := make([]GenreListItem, len(rows))
	for i, r := range rows {
		genres[i] = GenreListItem{
			ID:          r.ID,
			Name:        r.Name,
			Slug:        r.Slug,
			Description: r.Description,
			ShowCount:   r.ShowCount,
		}
	}

	respondJSON(c, http.StatusOK, genres)
}
