package handlers

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/paulsena/asheville-setlist/internal/db"
)

// showRowData holds common fields from all show list query results.
// This eliminates duplication across 8+ convert functions.
type showRowData struct {
	ID             int32
	Title          *string
	ImageUrl       *string
	Date           pgtype.Timestamptz
	DoorsTime      pgtype.Time
	ShowTime       pgtype.Time
	PriceMin       pgtype.Numeric
	PriceMax       pgtype.Numeric
	TicketUrl      *string
	AgeRestriction *string
	Status         *string
	VenueID        int32
	VenueName      string
	VenueSlug      string
	VenueRegion    *string
	VenueAddress   *string
	VenueImageUrl  *string
}

// convertShowRowToListItem converts common show row data to ShowListItem.
func convertShowRowToListItem(r showRowData) ShowListItem {
	return ShowListItem{
		ID:             r.ID,
		Title:          r.Title,
		ImageURL:       r.ImageUrl,
		Date:           formatTimestamp(r.Date),
		DoorsTime:      formatTime(r.DoorsTime),
		ShowTime:       formatTime(r.ShowTime),
		PriceMin:       numericToFloat(r.PriceMin),
		PriceMax:       numericToFloat(r.PriceMax),
		TicketURL:      r.TicketUrl,
		AgeRestriction: r.AgeRestriction,
		Status:         stringValue(r.Status),
		Venue: VenueBasic{
			ID:       r.VenueID,
			Name:     r.VenueName,
			Slug:     r.VenueSlug,
			Region:   r.VenueRegion,
			Address:  r.VenueAddress,
			ImageURL: r.VenueImageUrl,
		},
		Bands: []BandBasic{},
	}
}

// Show list conversion functions - extract common data and delegate to single converter.

func convertUpcomingShowsToListItems(rows []db.ListUpcomingShowsRow) ([]ShowListItem, int) {
	if len(rows) == 0 {
		return []ShowListItem{}, 0
	}
	items := make([]ShowListItem, len(rows))
	for i, r := range rows {
		items[i] = convertShowRowToListItem(showRowData{
			ID: r.ID, Title: r.Title, ImageUrl: r.ImageUrl, Date: r.Date,
			DoorsTime: r.DoorsTime, ShowTime: r.ShowTime, PriceMin: r.PriceMin,
			PriceMax: r.PriceMax, TicketUrl: r.TicketUrl, AgeRestriction: r.AgeRestriction,
			Status: r.Status, VenueID: r.VenueID, VenueName: r.VenueName,
			VenueSlug: r.VenueSlug, VenueRegion: r.VenueRegion,
			VenueAddress: r.VenueAddress, VenueImageUrl: r.VenueImageUrl,
		})
	}
	return items, int(rows[0].TotalCount)
}

func convertTonightShowsToListItems(rows []db.ListShowsTonightRow) []ShowListItem {
	items := make([]ShowListItem, len(rows))
	for i, r := range rows {
		items[i] = convertShowRowToListItem(showRowData{
			ID: r.ID, Title: r.Title, ImageUrl: r.ImageUrl, Date: r.Date,
			DoorsTime: r.DoorsTime, ShowTime: r.ShowTime, PriceMin: r.PriceMin,
			PriceMax: r.PriceMax, TicketUrl: r.TicketUrl, AgeRestriction: r.AgeRestriction,
			Status: r.Status, VenueID: r.VenueID, VenueName: r.VenueName,
			VenueSlug: r.VenueSlug, VenueRegion: r.VenueRegion,
			VenueAddress: r.VenueAddress, VenueImageUrl: r.VenueImageUrl,
		})
	}
	return items
}

func convertWeekendShowsToListItems(rows []db.ListShowsThisWeekendRow) []ShowListItem {
	items := make([]ShowListItem, len(rows))
	for i, r := range rows {
		items[i] = convertShowRowToListItem(showRowData{
			ID: r.ID, Title: r.Title, ImageUrl: r.ImageUrl, Date: r.Date,
			DoorsTime: r.DoorsTime, ShowTime: r.ShowTime, PriceMin: r.PriceMin,
			PriceMax: r.PriceMax, TicketUrl: r.TicketUrl, AgeRestriction: r.AgeRestriction,
			Status: r.Status, VenueID: r.VenueID, VenueName: r.VenueName,
			VenueSlug: r.VenueSlug, VenueRegion: r.VenueRegion,
			VenueAddress: r.VenueAddress, VenueImageUrl: r.VenueImageUrl,
		})
	}
	return items
}

func convertFreeShowsToListItems(rows []db.ListFreeShowsRow) ([]ShowListItem, int) {
	if len(rows) == 0 {
		return []ShowListItem{}, 0
	}
	items := make([]ShowListItem, len(rows))
	for i, r := range rows {
		items[i] = convertShowRowToListItem(showRowData{
			ID: r.ID, Title: r.Title, ImageUrl: r.ImageUrl, Date: r.Date,
			DoorsTime: r.DoorsTime, ShowTime: r.ShowTime, PriceMin: r.PriceMin,
			PriceMax: r.PriceMax, TicketUrl: r.TicketUrl, AgeRestriction: r.AgeRestriction,
			Status: r.Status, VenueID: r.VenueID, VenueName: r.VenueName,
			VenueSlug: r.VenueSlug, VenueRegion: r.VenueRegion,
			VenueAddress: r.VenueAddress, VenueImageUrl: r.VenueImageUrl,
		})
	}
	return items, int(rows[0].TotalCount)
}

func convertVenueShowsToListItems(rows []db.ListShowsByVenueRow) ([]ShowListItem, int) {
	if len(rows) == 0 {
		return []ShowListItem{}, 0
	}
	items := make([]ShowListItem, len(rows))
	for i, r := range rows {
		items[i] = convertShowRowToListItem(showRowData{
			ID: r.ID, Title: r.Title, ImageUrl: r.ImageUrl, Date: r.Date,
			DoorsTime: r.DoorsTime, ShowTime: r.ShowTime, PriceMin: r.PriceMin,
			PriceMax: r.PriceMax, TicketUrl: r.TicketUrl, AgeRestriction: r.AgeRestriction,
			Status: r.Status, VenueID: r.VenueID, VenueName: r.VenueName,
			VenueSlug: r.VenueSlug, VenueRegion: r.VenueRegion,
			VenueAddress: r.VenueAddress, VenueImageUrl: r.VenueImageUrl,
		})
	}
	return items, int(rows[0].TotalCount)
}

func convertRegionShowsToListItems(rows []db.ListShowsByRegionRow) ([]ShowListItem, int) {
	if len(rows) == 0 {
		return []ShowListItem{}, 0
	}
	items := make([]ShowListItem, len(rows))
	for i, r := range rows {
		items[i] = convertShowRowToListItem(showRowData{
			ID: r.ID, Title: r.Title, ImageUrl: r.ImageUrl, Date: r.Date,
			DoorsTime: r.DoorsTime, ShowTime: r.ShowTime, PriceMin: r.PriceMin,
			PriceMax: r.PriceMax, TicketUrl: r.TicketUrl, AgeRestriction: r.AgeRestriction,
			Status: r.Status, VenueID: r.VenueID, VenueName: r.VenueName,
			VenueSlug: r.VenueSlug, VenueRegion: r.VenueRegion,
			VenueAddress: r.VenueAddress, VenueImageUrl: r.VenueImageUrl,
		})
	}
	return items, int(rows[0].TotalCount)
}

func convertGenreShowsToListItems(rows []db.ListShowsByGenreRow) []ShowListItem {
	items := make([]ShowListItem, len(rows))
	for i, r := range rows {
		items[i] = convertShowRowToListItem(showRowData{
			ID: r.ID, Title: r.Title, ImageUrl: r.ImageUrl, Date: r.Date,
			DoorsTime: r.DoorsTime, ShowTime: r.ShowTime, PriceMin: r.PriceMin,
			PriceMax: r.PriceMax, TicketUrl: r.TicketUrl, AgeRestriction: r.AgeRestriction,
			Status: r.Status, VenueID: r.VenueID, VenueName: r.VenueName,
			VenueSlug: r.VenueSlug, VenueRegion: r.VenueRegion,
			VenueAddress: r.VenueAddress, VenueImageUrl: r.VenueImageUrl,
		})
	}
	return items
}

func convertDateRangeShowsToListItems(rows []db.ListShowsByDateRangeRow) ([]ShowListItem, int) {
	if len(rows) == 0 {
		return []ShowListItem{}, 0
	}
	items := make([]ShowListItem, len(rows))
	for i, r := range rows {
		items[i] = convertShowRowToListItem(showRowData{
			ID: r.ID, Title: r.Title, ImageUrl: r.ImageUrl, Date: r.Date,
			DoorsTime: r.DoorsTime, ShowTime: r.ShowTime, PriceMin: r.PriceMin,
			PriceMax: r.PriceMax, TicketUrl: r.TicketUrl, AgeRestriction: r.AgeRestriction,
			Status: r.Status, VenueID: r.VenueID, VenueName: r.VenueName,
			VenueSlug: r.VenueSlug, VenueRegion: r.VenueRegion,
			VenueAddress: r.VenueAddress, VenueImageUrl: r.VenueImageUrl,
		})
	}
	return items, int(rows[0].TotalCount)
}

// Band list conversion functions.

func convertBandsToListItems(rows []db.ListBandsRow) ([]BandListItem, int) {
	if len(rows) == 0 {
		return []BandListItem{}, 0
	}
	items := make([]BandListItem, len(rows))
	for i, r := range rows {
		items[i] = BandListItem{
			ID:       r.ID,
			Name:     r.Name,
			Slug:     r.Slug,
			Bio:      r.Bio,
			Hometown: r.Hometown,
			ImageURL: r.ImageUrl,
			Genres:   []GenreBasic{},
		}
	}
	return items, int(rows[0].TotalCount)
}

func convertSearchBandsToListItems(rows []db.SearchBandsRow) ([]BandListItem, int) {
	if len(rows) == 0 {
		return []BandListItem{}, 0
	}
	items := make([]BandListItem, len(rows))
	for i, r := range rows {
		items[i] = BandListItem{
			ID:       r.ID,
			Name:     r.Name,
			Slug:     r.Slug,
			Bio:      r.Bio,
			Hometown: r.Hometown,
			ImageURL: r.ImageUrl,
			Genres:   []GenreBasic{},
		}
	}
	return items, int(rows[0].TotalCount)
}

func convertGenreBandsToListItems(rows []db.ListBandsByGenreRow) []BandListItem {
	items := make([]BandListItem, len(rows))
	for i, r := range rows {
		items[i] = BandListItem{
			ID:       r.ID,
			Name:     r.Name,
			Slug:     r.Slug,
			Bio:      r.Bio,
			Hometown: r.Hometown,
			ImageURL: r.ImageUrl,
			Genres:   []GenreBasic{},
		}
	}
	return items
}

// Venue list conversion functions.

func convertVenuesToListItems(rows []db.ListVenuesWithShowCountRow) []VenueListItem {
	items := make([]VenueListItem, len(rows))
	for i, r := range rows {
		items[i] = VenueListItem{
			ID:                r.ID,
			Name:              r.Name,
			Slug:              r.Slug,
			Address:           r.Address,
			Region:            r.Region,
			Capacity:          r.Capacity,
			Website:           r.Website,
			ImageURL:          r.ImageUrl,
			UpcomingShowCount: r.UpcomingShowCount,
		}
	}
	return items
}

func convertRegionVenuesToListItems(rows []db.ListVenuesByRegionRow) []VenueListItem {
	items := make([]VenueListItem, len(rows))
	for i, r := range rows {
		items[i] = VenueListItem{
			ID:                r.ID,
			Name:              r.Name,
			Slug:              r.Slug,
			Address:           r.Address,
			Region:            r.Region,
			Capacity:          r.Capacity,
			Website:           r.Website,
			ImageURL:          r.ImageUrl,
			UpcomingShowCount: r.UpcomingShowCount,
		}
	}
	return items
}

// Genre conversion functions.

func convertGenreRows(rows []db.GetBandGenresForShowRow) []GenreBasic {
	genres := make([]GenreBasic, len(rows))
	for i, r := range rows {
		genres[i] = GenreBasic{
			ID:   r.ID,
			Name: r.Name,
			Slug: r.Slug,
		}
	}
	return genres
}
