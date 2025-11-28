# Scraper Configuration

## Overview

Two scraping approaches based on venue website technology:

1. **Static Scraper (Colly)** - Server-side rendered HTML
2. **JavaScript Scraper (chromedp)** - Client-side rendered content

---

## Venue Configurations

### The Orange Peel

**URL**: `https://theorangepeel.net/events`

**Technology**: WordPress + RockHouse Platform (ETIX ticketing)

**Scraper Type**: `static` (server-side HTML)

**Configuration**:
```json
{
  "venue_id": 1,
  "url": "https://theorangepeel.net/events",
  "scraper_type": "static",
  "selectors": {
    "container": ".eventWrapper",
    "title": "#eventTitle h2",
    "title_alt": ".eventSeriesTitle",
    "date_month": ".eventMonth",
    "date_day": "#eventDate",
    "time": ".eventTime",
    "price": ".eventMoreInfo a",
    "ticket_url": ".rhp-event-cta a",
    "age_restriction": ".eventAgeRestriction",
    "bands": "#eventTitle h4 a"
  },
  "date_format": "Mon, Jan 2",
  "time_format": "Show: 3:04 pm | Doors: 3:04 pm",
  "pagination": null
}
```

**Parsing Notes**:
- Date split across two elements: month (`.eventMonth`) + day (`#eventDate`)
- Multiple artists: Headliner in `h2`, supporting acts in `h4` links
- Age restriction in separate element (e.g., "All Ages", "21+")
- Ticket URL from `.rhp-event-cta a[href]`

**Example HTML Structure**:
```html
<div class="eventWrapper">
  <div class="eventMonth">Nov</div>
  <div id="eventDate">28</div>
  <div id="eventTitle">
    <h2>Whitechapel</h2>
    <h4><a href="...">Bodysnatcher</a></h4>
    <h4><a href="...">AngelMaker</a></h4>
  </div>
  <div class="eventTime">Show: 8 pm | Doors: 7 pm</div>
  <div class="eventAgeRestriction">All Ages</div>
  <a class="rhp-event-cta" href="https://tickets.etix.com/...">Buy Tickets</a>
</div>
```

---

### Salvage Station

**URL**: `https://salvagestation.com/events`

**Technology**: WordPress + Events Manager plugin

**Scraper Type**: `javascript` (AJAX-loaded content)

**Configuration**:
```json
{
  "venue_id": 3,
  "url": "https://salvagestation.com/events",
  "scraper_type": "javascript",
  "ajax_endpoint": "https://salvagestation.com/wp-admin/admin-ajax.php",
  "wait_selector": ".em-events-list",
  "wait_timeout": 5000,
  "selectors": {
    "container": ".em-event",
    "title": ".em-event-title",
    "date": ".em-event-date",
    "time": ".em-event-time",
    "price": ".em-event-price",
    "ticket_url": ".em-event-link a",
    "description": ".em-event-excerpt"
  },
  "date_format": "m/d/Y",
  "ajax_action": "em_ajax_get_events"
}
```

**Parsing Notes**:
- Content loads via AJAX after page renders
- Must wait for `.em-events-list` to populate
- Date format: "12/25/2024" (frontend display)
- Backend format: "yy-mm-dd"
- Events Manager plugin structure

**JavaScript Scraping Strategy**:
```go
// Use chromedp to render JavaScript
ctx, cancel := chromedp.NewContext(context.Background())
defer cancel()

var htmlContent string
err := chromedp.Run(ctx,
    chromedp.Navigate(config.URL),
    chromedp.WaitVisible(config.WaitSelector, chromedp.ByQuery),
    chromedp.Sleep(2 * time.Second), // Allow AJAX to complete
    chromedp.OuterHTML("body", &htmlContent),
)
```

---

### General Venue Pattern (Template)

For venues not yet configured:

```json
{
  "venue_id": null,
  "url": "",
  "scraper_type": "static",
  "selectors": {
    "container": "",
    "title": "",
    "date": "",
    "time": "",
    "price": "",
    "ticket_url": "",
    "bands": ""
  },
  "date_format": "",
  "notes": ""
}
```

---

## Scraper Implementation

### Static Scraper (Colly)

```go
package scraper

import (
    "github.com/gocolly/colly/v2"
)

type StaticScraper struct {
    collector *colly.Collector
    config    VenueConfig
}

func NewStaticScraper(config VenueConfig) *StaticScraper {
    c := colly.NewCollector(
        colly.AllowedDomains(getDomain(config.URL)),
        colly.UserAgent("AshevilleSetlist/1.0"),
    )

    // Rate limiting: 1 request per second
    c.Limit(&colly.LimitRule{
        DomainGlob:  "*",
        Delay:       1 * time.Second,
        RandomDelay: 500 * time.Millisecond,
    })

    return &StaticScraper{
        collector: c,
        config:    config,
    }
}

func (s *StaticScraper) Scrape() ([]RawShow, error) {
    var shows []RawShow

    s.collector.OnHTML(s.config.Selectors.Container, func(e *colly.HTMLElement) {
        show := RawShow{
            Title:         e.ChildText(s.config.Selectors.Title),
            Date:          s.parseDate(e),
            TicketURL:     e.ChildAttr(s.config.Selectors.TicketURL, "href"),
            AgeRestriction: e.ChildText(s.config.Selectors.AgeRestriction),
            Bands:         s.extractBands(e),
            Source:        s.config.URL,
        }
        shows = append(shows, show)
    })

    err := s.collector.Visit(s.config.URL)
    return shows, err
}

func (s *StaticScraper) extractBands(e *colly.HTMLElement) []string {
    var bands []string
    e.ForEach(s.config.Selectors.Bands, func(_ int, el *colly.HTMLElement) {
        if name := el.Text; name != "" {
            bands = append(bands, strings.TrimSpace(name))
        }
    })
    return bands
}
```

---

### JavaScript Scraper (chromedp)

```go
package scraper

import (
    "context"
    "github.com/chromedp/chromedp"
    "github.com/PuerkitoBio/goquery"
)

type JavaScriptScraper struct {
    config VenueConfig
}

func NewJavaScriptScraper(config VenueConfig) *JavaScriptScraper {
    return &JavaScriptScraper{config: config}
}

func (s *JavaScriptScraper) Scrape() ([]RawShow, error) {
    ctx, cancel := chromedp.NewContext(context.Background())
    defer cancel()

    ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
    defer cancel()

    var htmlContent string
    err := chromedp.Run(ctx,
        chromedp.Navigate(s.config.URL),
        chromedp.WaitVisible(s.config.WaitSelector, chromedp.ByQuery),
        chromedp.Sleep(time.Duration(s.config.WaitTimeout) * time.Millisecond),
        chromedp.OuterHTML("body", &htmlContent),
    )

    if err != nil {
        return nil, fmt.Errorf("chromedp failed: %w", err)
    }

    return s.parseHTML(htmlContent)
}

func (s *JavaScriptScraper) parseHTML(html string) ([]RawShow, error) {
    doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
    if err != nil {
        return nil, err
    }

    var shows []RawShow
    doc.Find(s.config.Selectors.Container).Each(func(i int, sel *goquery.Selection) {
        show := RawShow{
            Title:     sel.Find(s.config.Selectors.Title).Text(),
            Date:      sel.Find(s.config.Selectors.Date).Text(),
            TicketURL: sel.Find(s.config.Selectors.TicketURL).AttrOr("href", ""),
            Source:    s.config.URL,
        }
        shows = append(shows, show)
    })

    return shows, nil
}
```

---

## Date Parsing

### Common Date Formats

| Format | Example | Go Parse String |
|--------|---------|----------------|
| "Mon, Jan 2" | "Fri, Nov 28" | Parse month+day separately |
| "m/d/Y" | "11/28/2025" | "1/2/2006" |
| "January 2, 2006" | "November 28, 2025" | "January 2, 2006" |
| "2006-01-02" | "2025-11-28" | "2006-01-02" |
| Relative | "Tonight", "Tomorrow" | Custom logic |

### Date Parser Implementation

```go
package scraper

import (
    "time"
    "strings"
)

func ParseDate(dateStr, formatHint string) (time.Time, error) {
    dateStr = strings.TrimSpace(dateStr)

    // Handle relative dates
    switch strings.ToLower(dateStr) {
    case "tonight", "today":
        return time.Now(), nil
    case "tomorrow":
        return time.Now().AddDate(0, 0, 1), nil
    }

    // Handle "This Friday", "Next Saturday" etc
    if strings.HasPrefix(strings.ToLower(dateStr), "this ") ||
       strings.HasPrefix(strings.ToLower(dateStr), "next ") {
        return parseRelativeWeekday(dateStr)
    }

    // Try standard formats
    formats := []string{
        "January 2, 2006",
        "Jan 2, 2006",
        "1/2/2006",
        "2006-01-02",
        "01/02/2006",
    }

    for _, format := range formats {
        if t, err := time.Parse(format, dateStr); err == nil {
            return t, nil
        }
    }

    return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

func parseRelativeWeekday(dateStr string) (time.Time, error) {
    // "This Friday" or "Next Saturday"
    parts := strings.Fields(dateStr)
    if len(parts) < 2 {
        return time.Time{}, fmt.Errorf("invalid relative date: %s", dateStr)
    }

    modifier := strings.ToLower(parts[0]) // "this" or "next"
    weekdayStr := parts[1]                 // "Friday"

    targetWeekday := parseWeekday(weekdayStr)
    today := time.Now()
    daysUntil := int(targetWeekday - today.Weekday())

    if modifier == "next" || daysUntil < 0 {
        daysUntil += 7
    }

    return today.AddDate(0, 0, daysUntil), nil
}
```

---

## Band Name Extraction

### Patterns to Handle

1. **Single headliner**: "Moon Taxi"
2. **With openers**: "Moon Taxi with The Revivalists"
3. **Multiple openers**: "Moon Taxi with The Revivalists and Neighbor"
4. **Featuring**: "Moon Taxi featuring special guest"
5. **Ampersand**: "Moon Taxi & The Revivalists"
6. **Comma separated**: "Moon Taxi, The Revivalists, Neighbor"
7. **Plus sign**: "Moon Taxi + The Revivalists"

### Band Name Extractor

```go
package scraper

import (
    "regexp"
    "strings"
)

type BandExtraction struct {
    Headliner string
    Openers   []string
}

func ExtractBands(title string) BandExtraction {
    title = strings.TrimSpace(title)

    // Patterns that indicate supporting acts
    separators := []string{
        " with ",
        " w/ ",
        " featuring ",
        " feat. ",
        " ft. ",
        " and ",
        " & ",
        " + ",
        ", ",
    }

    extraction := BandExtraction{}

    // Try to split on separators
    for _, sep := range separators {
        if strings.Contains(strings.ToLower(title), sep) {
            parts := splitCaseInsensitive(title, sep)
            extraction.Headliner = strings.TrimSpace(parts[0])

            if len(parts) > 1 {
                // Further split openers by commas/and
                openerStr := strings.Join(parts[1:], ", ")
                extraction.Openers = splitOpeners(openerStr)
            }
            return extraction
        }
    }

    // No separators found - single headliner
    extraction.Headliner = title
    return extraction
}

func splitOpeners(openerStr string) []string {
    // Split by ", " "and" "&"
    re := regexp.MustCompile(`(?:,\s*|\s+and\s+|\s*&\s*|\s*\+\s*)`)
    parts := re.Split(openerStr, -1)

    var openers []string
    for _, p := range parts {
        p = strings.TrimSpace(p)
        if p != "" && !isFillerWord(p) {
            openers = append(openers, p)
        }
    }
    return openers
}

func isFillerWord(s string) bool {
    filler := []string{"special guests", "more", "tba", "tbd"}
    s = strings.ToLower(s)
    for _, f := range filler {
        if s == f {
            return true
        }
    }
    return false
}
```

---

## Band Matching Algorithm

### Matching Strategy

When scraper extracts band name "Moon Taxi", determine if band exists in database.

**Steps**:
1. **Exact match** - Case-insensitive exact match on name or slug
2. **Fuzzy match** - Levenshtein distance < threshold
3. **Create new** - No match found, create new band

### Implementation

```go
package scraper

import (
    "strings"
    "github.com/agnivade/levenshtein"
)

const (
    FuzzyMatchThreshold = 3 // Max edit distance
)

type BandMatcher struct {
    db *sql.DB
}

func (m *BandMatcher) FindOrCreateBand(name string) (int, error) {
    // 1. Try exact match
    if id, err := m.findExactMatch(name); err == nil {
        return id, nil
    }

    // 2. Try fuzzy match
    if id, err := m.findFuzzyMatch(name); err == nil {
        return id, nil
    }

    // 3. Create new band
    return m.createBand(name)
}

func (m *BandMatcher) findExactMatch(name string) (int, error) {
    var id int
    err := m.db.QueryRow(`
        SELECT id FROM bands
        WHERE LOWER(name) = LOWER($1)
           OR slug = $2
        LIMIT 1
    `, name, slugify(name)).Scan(&id)

    return id, err
}

func (m *BandMatcher) findFuzzyMatch(name string) (int, error) {
    // Get all bands (or limit to first 1000 for performance)
    rows, err := m.db.Query(`
        SELECT id, name FROM bands
        ORDER BY created_at DESC
        LIMIT 1000
    `)
    if err != nil {
        return 0, err
    }
    defer rows.Close()

    bestMatch := struct {
        id       int
        distance int
    }{distance: FuzzyMatchThreshold + 1}

    for rows.Next() {
        var id int
        var dbName string
        if err := rows.Scan(&id, &dbName); err != nil {
            continue
        }

        distance := levenshtein.ComputeDistance(
            strings.ToLower(name),
            strings.ToLower(dbName),
        )

        if distance <= FuzzyMatchThreshold && distance < bestMatch.distance {
            bestMatch.id = id
            bestMatch.distance = distance
        }
    }

    if bestMatch.distance <= FuzzyMatchThreshold {
        return bestMatch.id, nil
    }

    return 0, fmt.Errorf("no fuzzy match found")
}

func (m *BandMatcher) createBand(name string) (int, error) {
    var id int
    err := m.db.QueryRow(`
        INSERT INTO bands (name, slug, source)
        VALUES ($1, $2, 'scraped')
        RETURNING id
    `, name, slugify(name)).Scan(&id)

    return id, err
}

func slugify(s string) string {
    s = strings.ToLower(s)
    s = regexp.MustCompile(`[^a-z0-9\s-]`).ReplaceAllString(s, "")
    s = regexp.MustCompile(`\s+`).ReplaceAllString(s, "-")
    s = strings.Trim(s, "-")
    return s
}
```

### Handling Duplicates

If slug collision occurs (e.g., two bands named "The Band"):

```go
func (m *BandMatcher) createBand(name string) (int, error) {
    slug := slugify(name)
    baseSlug := slug

    // Try appending -2, -3, etc. until unique
    for i := 2; i < 100; i++ {
        var id int
        err := m.db.QueryRow(`
            INSERT INTO bands (name, slug, source)
            VALUES ($1, $2, 'scraped')
            RETURNING id
        `, name, slug).Scan(&id)

        if err == nil {
            return id, nil
        }

        // Check if error is unique constraint violation
        if strings.Contains(err.Error(), "duplicate") {
            slug = fmt.Sprintf("%s-%d", baseSlug, i)
            continue
        }

        return 0, err
    }

    return 0, fmt.Errorf("could not create unique slug for: %s", name)
}
```

---

## Error Handling Strategy

### Failure Modes

| Error Type | Example | Handling |
|------------|---------|----------|
| **Network** | DNS failure, timeout | Retry 3x with exponential backoff |
| **Parsing** | Selector not found | Log error, skip event, continue |
| **Rate Limit** | 429 Too Many Requests | Wait + retry with longer delay |
| **JavaScript Timeout** | AJAX didn't load | Increase timeout, retry once |
| **Invalid Data** | Date can't be parsed | Log error, mark as manual review |
| **Database** | Insert fails | Rollback transaction, log error |

### Implementation

```go
package scraper

import (
    "time"
    "log/slog"
)

type ScraperJob struct {
    config  VenueConfig
    scraper Scraper
    logger  *slog.Logger
    metrics *Metrics
}

func (j *ScraperJob) Run() error {
    j.logger.Info("starting scrape", "venue", j.config.VenueID, "url", j.config.URL)

    var shows []RawShow
    var err error

    // Retry logic with exponential backoff
    for attempt := 1; attempt <= 3; attempt++ {
        shows, err = j.scraper.Scrape()

        if err == nil {
            break // Success
        }

        // Check if retryable
        if !isRetryable(err) {
            j.logger.Error("non-retryable error", "error", err)
            j.metrics.RecordFailure(j.config.VenueID, err)
            return err
        }

        // Exponential backoff: 2s, 4s, 8s
        backoff := time.Duration(1<<uint(attempt)) * time.Second
        j.logger.Warn("scrape failed, retrying",
            "attempt", attempt,
            "backoff", backoff,
            "error", err,
        )
        time.Sleep(backoff)
    }

    if err != nil {
        j.logger.Error("scrape failed after retries", "error", err)
        j.metrics.RecordFailure(j.config.VenueID, err)
        return err
    }

    j.logger.Info("scrape succeeded", "shows_found", len(shows))

    // Process shows (parse, match bands, upsert)
    created, updated, skipped, errors := j.processShows(shows)

    j.metrics.RecordSuccess(j.config.VenueID, created, updated, skipped)

    if len(errors) > 0 {
        j.logger.Warn("some shows had errors",
            "total", len(shows),
            "errors", len(errors),
        )
    }

    return nil
}

func isRetryable(err error) bool {
    // Network errors are retryable
    if strings.Contains(err.Error(), "timeout") ||
       strings.Contains(err.Error(), "connection refused") {
        return true
    }

    // 5xx errors are retryable
    if strings.Contains(err.Error(), "500") ||
       strings.Contains(err.Error(), "502") ||
       strings.Contains(err.Error(), "503") {
        return true
    }

    return false
}

func (j *ScraperJob) processShows(rawShows []RawShow) (created, updated, skipped int, errors []error) {
    for _, raw := range rawShows {
        if err := j.processShow(raw); err != nil {
            errors = append(errors, err)
            continue
        }

        // Determine if created/updated/skipped
        // ... (implementation details)
    }

    return
}
```

### Logging & Metrics

```go
type Metrics struct {
    mu sync.Mutex
    stats map[int]*VenueStats
}

type VenueStats struct {
    LastRun      time.Time
    LastSuccess  time.Time
    SuccessCount int
    FailureCount int
    ShowsCreated int
    ShowsUpdated int
    ShowsSkipped int
}

func (m *Metrics) RecordSuccess(venueID, created, updated, skipped int) {
    m.mu.Lock()
    defer m.mu.Unlock()

    if m.stats[venueID] == nil {
        m.stats[venueID] = &VenueStats{}
    }

    s := m.stats[venueID]
    s.LastRun = time.Now()
    s.LastSuccess = time.Now()
    s.SuccessCount++
    s.ShowsCreated += created
    s.ShowsUpdated += updated
    s.ShowsSkipped += skipped
}
```

### Alerting

```go
// Alert if scraper fails 3+ times in a row
func (m *Metrics) ShouldAlert(venueID int) bool {
    m.mu.Lock()
    defer m.mu.Unlock()

    s := m.stats[venueID]
    if s == nil {
        return false
    }

    // Alert if:
    // - 3+ consecutive failures
    // - No success in last 24 hours
    return s.FailureCount >= 3 ||
           time.Since(s.LastSuccess) > 24*time.Hour
}
```

---

## Orchestration

### Concurrent Scraping

```go
package scraper

func RunAllScrapers(configs []VenueConfig) {
    var wg sync.WaitGroup
    semaphore := make(chan struct{}, 5) // Max 5 concurrent scrapers

    for _, config := range configs {
        if !config.IsActive {
            continue
        }

        wg.Add(1)
        go func(cfg VenueConfig) {
            defer wg.Done()

            semaphore <- struct{}{} // Acquire
            defer func() { <-semaphore }() // Release

            job := NewScraperJob(cfg)
            if err := job.Run(); err != nil {
                log.Printf("Scraper failed for venue %d: %v", cfg.VenueID, err)
            }
        }(config)
    }

    wg.Wait()
}
```

### CLI

```go
// cmd/scraper/main.go
package main

import (
    "flag"
    "log"
)

func main() {
    venueSlug := flag.String("venue", "", "Scrape specific venue by slug")
    dryRun := flag.Bool("dry-run", false, "Don't save to database")
    flag.Parse()

    configs, err := loadVenueConfigs()
    if err != nil {
        log.Fatal(err)
    }

    if *venueSlug != "" {
        // Scrape single venue
        config := findConfigBySlug(configs, *venueSlug)
        if config == nil {
            log.Fatalf("Venue not found: %s", *venueSlug)
        }
        runScraper(*config, *dryRun)
    } else {
        // Scrape all venues
        RunAllScrapers(configs)
    }
}
```

**Usage**:
```bash
# Scrape all venues
./scraper run

# Scrape single venue
./scraper run --venue=orange-peel

# Dry run (don't save)
./scraper run --dry-run

# Scrape with verbose logging
LOG_LEVEL=debug ./scraper run
```

---

## Database Storage

### Upsert Logic

```go
func (s *ScraperService) UpsertShow(raw RawShow, venueID int) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // 1. Parse and validate
    date, err := ParseDate(raw.Date, "")
    if err != nil {
        return fmt.Errorf("invalid date: %w", err)
    }

    // 2. Extract bands
    extraction := ExtractBands(raw.Title)

    // 3. Match/create headliner
    headlinerID, err := s.matcher.FindOrCreateBand(extraction.Headliner)
    if err != nil {
        return err
    }

    // 4. Check if show exists (venue + date + headliner)
    var showID int
    err = tx.QueryRow(`
        SELECT s.id FROM shows s
        JOIN show_bands sb ON s.id = sb.show_id
        WHERE s.venue_id = $1
          AND DATE(s.date) = DATE($2)
          AND sb.band_id = $3
          AND sb.is_headliner = true
    `, venueID, date, headlinerID).Scan(&showID)

    if err == sql.ErrNoRows {
        // Create new show
        err = tx.QueryRow(`
            INSERT INTO shows (venue_id, title, date, ticket_url, source, scraped_data)
            VALUES ($1, $2, $3, $4, 'scraped', $5)
            RETURNING id
        `, venueID, raw.Title, date, raw.TicketURL, raw.ToJSON()).Scan(&showID)

        if err != nil {
            return err
        }
    } else if err != nil {
        return err
    } else {
        // Update existing show
        _, err = tx.Exec(`
            UPDATE shows
            SET title = $1,
                ticket_url = $2,
                scraped_data = $3,
                updated_at = NOW()
            WHERE id = $4
        `, raw.Title, raw.TicketURL, raw.ToJSON(), showID)

        if err != nil {
            return err
        }
    }

    // 5. Link bands (delete old, insert new)
    _, err = tx.Exec(`DELETE FROM show_bands WHERE show_id = $1`, showID)
    if err != nil {
        return err
    }

    // Insert headliner
    _, err = tx.Exec(`
        INSERT INTO show_bands (show_id, band_id, is_headliner, performance_order)
        VALUES ($1, $2, true, 1)
    `, showID, headlinerID)
    if err != nil {
        return err
    }

    // Insert openers
    for i, opener := range extraction.Openers {
        openerID, err := s.matcher.FindOrCreateBand(opener)
        if err != nil {
            continue // Skip if can't match
        }

        _, err = tx.Exec(`
            INSERT INTO show_bands (show_id, band_id, is_headliner, performance_order)
            VALUES ($1, $2, false, $3)
        `, showID, openerID, i+2)
        if err != nil {
            continue
        }
    }

    return tx.Commit()
}
```

---

## Testing

### Unit Tests

```go
func TestExtractBands(t *testing.T) {
    tests := []struct {
        input     string
        headliner string
        openers   []string
    }{
        {
            input:     "Moon Taxi",
            headliner: "Moon Taxi",
            openers:   []string{},
        },
        {
            input:     "Moon Taxi with The Revivalists",
            headliner: "Moon Taxi",
            openers:   []string{"The Revivalists"},
        },
        {
            input:     "Moon Taxi with The Revivalists and Neighbor",
            headliner: "Moon Taxi",
            openers:   []string{"The Revivalists", "Neighbor"},
        },
        {
            input:     "Moon Taxi & The Revivalists",
            headliner: "Moon Taxi",
            openers:   []string{"The Revivalists"},
        },
    }

    for _, tt := range tests {
        t.Run(tt.input, func(t *testing.T) {
            result := ExtractBands(tt.input)
            assert.Equal(t, tt.headliner, result.Headliner)
            assert.Equal(t, tt.openers, result.Openers)
        })
    }
}
```

---

## Summary

**Implementation Checklist**:
- [ ] Implement static scraper (Colly)
- [ ] Implement JavaScript scraper (chromedp)
- [ ] Date parser with relative dates
- [ ] Band name extractor
- [ ] Band matching algorithm (exact + fuzzy)
- [ ] Error handling with retries
- [ ] Upsert logic (show + bands)
- [ ] Concurrent orchestration
- [ ] CLI interface
- [ ] Logging and metrics
- [ ] Unit tests
