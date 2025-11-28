package handlers

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// Type conversion helpers for database types to API response types.

// formatTimestamp converts a pgtype.Timestamptz to RFC3339 string.
func formatTimestamp(ts pgtype.Timestamptz) string {
	if !ts.Valid {
		return ""
	}
	return ts.Time.Format(time.RFC3339)
}

// formatTime converts a pgtype.Time to HH:MM:SS string.
func formatTime(t pgtype.Time) *string {
	if !t.Valid {
		return nil
	}
	hours := t.Microseconds / 3600000000
	minutes := (t.Microseconds % 3600000000) / 60000000
	seconds := (t.Microseconds % 60000000) / 1000000
	s := time.Date(0, 1, 1, int(hours), int(minutes), int(seconds), 0, time.UTC).Format("15:04:05")
	return &s
}

// numericToFloat converts a pgtype.Numeric to *float64.
func numericToFloat(n pgtype.Numeric) *float64 {
	if !n.Valid {
		return nil
	}
	f, _ := n.Float64Value()
	if !f.Valid {
		return nil
	}
	return &f.Float64
}

// Pointer value helpers for safe dereferencing.

// stringValue safely dereferences a *string, returning empty string if nil.
func stringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// boolValue safely dereferences a *bool, returning false if nil.
func boolValue(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// int32Value safely dereferences a *int32, returning 0 if nil.
func int32Value(i *int32) int {
	if i == nil {
		return 0
	}
	return int(*i)
}

// Date parsing helpers.

// parseDateRange parses date_from and date_to query params into pgtype.Timestamptz.
// If from is empty, defaults to now. If to is empty, defaults to 1 year from now.
func parseDateRange(from, to string) (pgtype.Timestamptz, pgtype.Timestamptz, error) {
	var fromTime, toTime pgtype.Timestamptz

	if from != "" {
		t, err := time.Parse("2006-01-02", from)
		if err != nil {
			t, err = time.Parse(time.RFC3339, from)
			if err != nil {
				return fromTime, toTime, err
			}
		}
		fromTime = pgtype.Timestamptz{Time: t, Valid: true}
	} else {
		fromTime = pgtype.Timestamptz{Time: time.Now(), Valid: true}
	}

	if to != "" {
		t, err := time.Parse("2006-01-02", to)
		if err != nil {
			t, err = time.Parse(time.RFC3339, to)
			if err != nil {
				return fromTime, toTime, err
			}
		}
		// End of day for "to" date
		t = t.Add(24*time.Hour - time.Second)
		toTime = pgtype.Timestamptz{Time: t, Valid: true}
	} else {
		// Default to 1 year from now
		toTime = pgtype.Timestamptz{Time: time.Now().AddDate(1, 0, 0), Valid: true}
	}

	return fromTime, toTime, nil
}
