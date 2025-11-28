package handlers

import (
	"github.com/paulsena/asheville-setlist/internal/db"
)

// Handler contains all HTTP handlers and their dependencies
type Handler struct {
	queries *db.Queries
}

// New creates a new Handler with the given dependencies
func New(queries *db.Queries) *Handler {
	return &Handler{
		queries: queries,
	}
}
