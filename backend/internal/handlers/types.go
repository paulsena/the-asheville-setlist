package handlers

// Response types for API endpoints.
// These types define the JSON structure returned by handlers.

// ShowListItem represents a show in list responses.
type ShowListItem struct {
	ID             int32       `json:"id"`
	Title          *string     `json:"title"`
	ImageURL       *string     `json:"image_url"`
	Date           string      `json:"date"`
	DoorsTime      *string     `json:"doors_time"`
	ShowTime       *string     `json:"show_time"`
	PriceMin       *float64    `json:"price_min"`
	PriceMax       *float64    `json:"price_max"`
	TicketURL      *string     `json:"ticket_url"`
	AgeRestriction *string     `json:"age_restriction"`
	Status         string      `json:"status"`
	Venue          VenueBasic  `json:"venue"`
	Bands          []BandBasic `json:"bands"`
}

// ShowDetail represents a show in detail response with full information.
type ShowDetail struct {
	ID             int32         `json:"id"`
	Title          *string       `json:"title"`
	Description    *string       `json:"description"`
	ImageURL       *string       `json:"image_url"`
	Date           string        `json:"date"`
	DoorsTime      *string       `json:"doors_time"`
	ShowTime       *string       `json:"show_time"`
	PriceMin       *float64      `json:"price_min"`
	PriceMax       *float64      `json:"price_max"`
	TicketURL      *string       `json:"ticket_url"`
	AgeRestriction *string       `json:"age_restriction"`
	Status         string        `json:"status"`
	Venue          VenueForShow  `json:"venue"`
	Bands          []BandForShow `json:"bands"`
}

// VenueBasic represents minimal venue info embedded in other responses.
type VenueBasic struct {
	ID       int32   `json:"id"`
	Name     string  `json:"name"`
	Slug     string  `json:"slug"`
	Region   *string `json:"region,omitempty"`
	Address  *string `json:"address,omitempty"`
	ImageURL *string `json:"image_url,omitempty"`
}

// VenueForShow represents venue info in show detail response.
type VenueForShow struct {
	ID       int32   `json:"id"`
	Name     string  `json:"name"`
	Slug     string  `json:"slug"`
	Address  *string `json:"address"`
	Region   *string `json:"region"`
	Website  *string `json:"website"`
	ImageURL *string `json:"image_url"`
}

// VenueListItem represents a venue in list responses.
type VenueListItem struct {
	ID                int32   `json:"id"`
	Name              string  `json:"name"`
	Slug              string  `json:"slug"`
	Address           *string `json:"address"`
	Region            *string `json:"region"`
	Capacity          *int32  `json:"capacity"`
	Website           *string `json:"website"`
	ImageURL          *string `json:"image_url"`
	UpcomingShowCount int64   `json:"upcoming_show_count"`
}

// VenueDetail represents a venue in detail response with full information.
type VenueDetail struct {
	ID            int32           `json:"id"`
	Name          string          `json:"name"`
	Slug          string          `json:"slug"`
	Address       *string         `json:"address"`
	City          string          `json:"city"`
	State         string          `json:"state"`
	ZipCode       *string         `json:"zip_code"`
	Region        *string         `json:"region"`
	Capacity      *int32          `json:"capacity"`
	Website       *string         `json:"website"`
	Phone         *string         `json:"phone"`
	ImageURL      *string         `json:"image_url"`
	UpcomingShows []VenueShowItem `json:"upcoming_shows"`
}

// VenueShowItem represents a show in venue detail response.
type VenueShowItem struct {
	ID       int32       `json:"id"`
	Title    *string     `json:"title"`
	Date     string      `json:"date"`
	PriceMin *float64    `json:"price_min"`
	PriceMax *float64    `json:"price_max"`
	Bands    []BandBasic `json:"bands"`
}

// BandBasic represents minimal band info embedded in other responses.
type BandBasic struct {
	ID               int32   `json:"id"`
	Name             string  `json:"name"`
	Slug             string  `json:"slug"`
	ImageURL         *string `json:"image_url,omitempty"`
	IsHeadliner      bool    `json:"is_headliner"`
	PerformanceOrder int     `json:"performance_order"`
}

// BandForShow represents band info in show detail response.
type BandForShow struct {
	ID               int32        `json:"id"`
	Name             string       `json:"name"`
	Slug             string       `json:"slug"`
	Bio              *string      `json:"bio"`
	ImageURL         *string      `json:"image_url"`
	SpotifyURL       *string      `json:"spotify_url"`
	Website          *string      `json:"website"`
	IsHeadliner      bool         `json:"is_headliner"`
	PerformanceOrder int          `json:"performance_order"`
	Genres           []GenreBasic `json:"genres"`
}

// BandListItem represents a band in list responses.
type BandListItem struct {
	ID       int32        `json:"id"`
	Name     string       `json:"name"`
	Slug     string       `json:"slug"`
	Bio      *string      `json:"bio"`
	Hometown *string      `json:"hometown"`
	ImageURL *string      `json:"image_url"`
	Genres   []GenreBasic `json:"genres"`
}

// BandDetail represents a band in detail response with full information.
type BandDetail struct {
	ID            int32          `json:"id"`
	Name          string         `json:"name"`
	Slug          string         `json:"slug"`
	Bio           *string        `json:"bio"`
	Hometown      *string        `json:"hometown"`
	ImageURL      *string        `json:"image_url"`
	Website       *string        `json:"website"`
	SpotifyURL    *string        `json:"spotify_url"`
	Instagram     *string        `json:"instagram"`
	Facebook      *string        `json:"facebook"`
	BandcampURL   *string        `json:"bandcamp_url"`
	Genres        []GenreBasic   `json:"genres"`
	UpcomingShows []BandShowItem `json:"upcoming_shows"`
}

// BandShowItem represents a show in band detail response.
type BandShowItem struct {
	ID          int32      `json:"id"`
	Date        string     `json:"date"`
	Venue       VenueBasic `json:"venue"`
	IsHeadliner bool       `json:"is_headliner"`
}

// SimilarBandItem represents a similar band with shared genre information.
type SimilarBandItem struct {
	ID               int32        `json:"id"`
	Name             string       `json:"name"`
	Slug             string       `json:"slug"`
	ImageURL         *string      `json:"image_url"`
	SharedGenreCount int64        `json:"shared_genre_count"`
	SharedGenres     []GenreBasic `json:"shared_genres"`
}

// GenreBasic represents minimal genre info embedded in other responses.
type GenreBasic struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// GenreListItem represents a genre in list responses.
type GenreListItem struct {
	ID          int32   `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description *string `json:"description"`
	ShowCount   int64   `json:"show_count,omitempty"`
}

// SearchResult represents the global search response with categorized results.
type SearchResult struct {
	Shows  []SearchShowItem  `json:"shows"`
	Bands  []SearchBandItem  `json:"bands"`
	Venues []SearchVenueItem `json:"venues"`
}

// SearchShowItem represents a show in search results.
type SearchShowItem struct {
	ID        int32   `json:"id"`
	Title     *string `json:"title"`
	Date      string  `json:"date"`
	VenueName string  `json:"venue_name"`
}

// SearchBandItem represents a band in search results.
type SearchBandItem struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// SearchVenueItem represents a venue in search results.
type SearchVenueItem struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// CreateShowRequest represents the request body for creating a show submission.
type CreateShowRequest struct {
	VenueID        int32            `json:"venue_id" binding:"required"`
	Date           string           `json:"date" binding:"required"`
	ImageURL       *string          `json:"image_url"`
	DoorsTime      *string          `json:"doors_time"`
	ShowTime       *string          `json:"show_time"`
	PriceMin       *float64         `json:"price_min"`
	PriceMax       *float64         `json:"price_max"`
	TicketURL      *string          `json:"ticket_url"`
	AgeRestriction *string          `json:"age_restriction"`
	Bands          []CreateShowBand `json:"bands" binding:"required,min=1"`
}

// CreateShowBand represents a band in the create show request.
type CreateShowBand struct {
	Name             string `json:"name" binding:"required"`
	IsHeadliner      *bool  `json:"is_headliner"`
	PerformanceOrder *int32 `json:"performance_order"`
}

// CreateShowResponse represents the response for creating a show.
type CreateShowResponse struct {
	ID        int32  `json:"id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}
