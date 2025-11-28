# API Specification

## Overview

RESTful JSON API for concert discovery. All dates in ISO 8601. All endpoints return JSON.

Base URL: `/api`

---

## Response Standards

### Success Response Envelope

All successful responses use this structure:

```typescript
{
  data: T | T[];           // Single resource or array
  meta?: {                 // Present on paginated lists only
    page: number;
    per_page: number;
    total: number;
    total_pages: number;
  };
}
```

### Error Response

```typescript
{
  error: {
    code: string;          // SCREAMING_SNAKE_CASE
    message: string;       // Human-readable
    details?: object;      // Optional field-specific errors
    stack?: string;        // Only in development mode
  };
}
```

**Standard Error Codes:**
- `VALIDATION_ERROR` - Invalid input data
- `INVALID_PARAMETER` - Invalid query parameter
- `MISSING_PARAMETER` - Required parameter missing
- `NOT_FOUND` - Resource doesn't exist
- `INTERNAL_ERROR` - Unexpected server error
- `DATABASE_ERROR` - Database operation failed

### Pagination

**Parameters:**
- `page` - Default: 1
- `per_page` - Default: 50, Max: 100

**Validation:**
- Reject `per_page > 100` with `INVALID_PARAMETER`
- Reject `page < 1` with `INVALID_PARAMETER`

### Query Parameter Conventions

- **snake_case** for all parameter names
- **Repeated parameters** for multiple values: `?genre=rock&genre=indie`
- **ISO 8601** for dates: `2025-11-15` or `2025-11-15T20:00:00-05:00`
- **Boolean** as strings: `true` or `false`
- **Sort** with `-` prefix for descending: `sort=-date`

---

## Shows Endpoints

### `GET /api/shows`

List shows with filtering and pagination.

**Query Parameters:**

```typescript
{
  // Pagination
  page?: number;              // Default: 1
  per_page?: number;          // Default: 50, Max: 100

  // Date filters
  date?: string;              // Exact date match (ISO 8601)
  date_from?: string;         // Shows on/after (inclusive)
  date_to?: string;           // Shows on/before (inclusive)

  // Location
  venue?: string[];           // Venue slug(s), repeatable
  region?: string[];          // Region(s), repeatable

  // Genre
  genre?: string[];           // Genre slug(s), repeatable

  // Price
  price_min?: number;         // Minimum price filter
  price_max?: number;         // Maximum price filter

  // Status
  status?: string;            // Default: "scheduled", or "all"

  // Special filters
  filter?: string;            // "popular" | "trending" | "tonight" | "this-weekend" | "free" | "featured"

  // Search
  q?: string;                 // Full-text search in title/bands

  // Sorting
  sort?: string;              // "date" | "-date" | "price" | "-price", Default: "date"
}
```

**Validation Rules:**
- `date`, `date_from`, `date_to` must be valid ISO 8601 dates
- `price_min`, `price_max` must be >= 0
- `status` must be "scheduled" or "all"
- `filter` must be one of: popular, trending, tonight, this-weekend, free, featured
- `sort` must be: date, -date, price, -price

**Filter Logic:**
- `date` - Exact match on date (ignores time)
- `date_from` - `show.date >= date_from`
- `date_to` - `show.date <= date_to`
- `venue` - Match any of the provided venue slugs (OR logic)
- `region` - Match any of the provided regions (OR logic)
- `genre` - Shows with bands matching any genre (OR logic)
- `price_min` - `show.price_min >= price_min`
- `price_max` - `show.price_max <= price_max`
- `status=scheduled` - Only shows with status='scheduled'
- `status=all` - Include cancelled/postponed
- `filter=tonight` - Shows where DATE(date) = TODAY
- `filter=this-weekend` - Shows where date between next Friday-Sunday
- `filter=free` - Shows where price_min IS NULL or price_min = 0
- `q` - Full-text search on show.title and band names

**Response:**

```typescript
{
  data: {
    id: number;
    title: string | null;
    image_url: string | null;
    date: string;                    // ISO 8601 with timezone
    doors_time: string | null;       // "HH:MM:SS"
    show_time: string | null;        // "HH:MM:SS"
    price_min: number | null;
    price_max: number | null;
    ticket_url: string | null;
    age_restriction: string | null;
    status: string;

    venue: {
      id: number;
      name: string;
      slug: string;
      region: string | null;
      address: string | null;
      image_url: string | null;
    };

    bands: {
      id: number;
      name: string;
      slug: string;
      image_url: string | null;
      is_headliner: boolean;
      performance_order: number;
    }[];
  }[];

  meta: {
    page: number;
    per_page: number;
    total: number;
    total_pages: number;
  };
}
```

**SQL Notes:**
- Default ORDER BY: `date ASC, id ASC`
- Join shows → venues (required)
- Join shows → show_bands → bands (required)
- Join show_bands → band_genres → genres (for genre filter)
- Default WHERE: `status = 'scheduled' AND date >= NOW()`

---

### `GET /api/shows/:id`

Get single show with full details.

**Path Parameters:**
- `id` - Show ID (integer)

**Response:**

```typescript
{
  data: {
    id: number;
    title: string | null;
    description: string | null;
    image_url: string | null;
    date: string;
    doors_time: string | null;
    show_time: string | null;
    price_min: number | null;
    price_max: number | null;
    ticket_url: string | null;
    age_restriction: string | null;
    status: string;

    venue: {
      id: number;
      name: string;
      slug: string;
      address: string | null;
      region: string | null;
      website: string | null;
      image_url: string | null;
    };

    bands: {
      id: number;
      name: string;
      slug: string;
      bio: string | null;
      image_url: string | null;
      spotify_url: string | null;
      website: string | null;
      is_headliner: boolean;
      performance_order: number;
      genres: {
        id: number;
        name: string;
        slug: string;
      }[];
    }[];
  };
}
```

**Errors:**
- `404 NOT_FOUND` - Show ID doesn't exist

**SQL Notes:**
- Join shows → venues
- Join shows → show_bands → bands
- Join bands → band_genres → genres
- Order bands by `performance_order DESC`

---

### `POST /api/shows`

Band submission - create show with pending status.

**Request Body:**

```typescript
{
  venue_id: number;              // Required
  date: string;                  // Required, ISO 8601
  image_url?: string;
  doors_time?: string;           // "HH:MM" or "HH:MM:SS"
  show_time?: string;            // "HH:MM" or "HH:MM:SS"
  price_min?: number;
  price_max?: number;
  ticket_url?: string;
  age_restriction?: string;

  bands: {                       // At least 1 required
    name: string;                // Required
    is_headliner?: boolean;      // Default: false
    performance_order?: number;
  }[];
}
```

**Validation Rules:**
- `venue_id` must exist in venues table
- `date` must be valid ISO 8601, must be future date
- `bands` array must have at least 1 element
- `bands[].name` must not be empty
- `price_min` >= 0 if provided
- `price_max` >= 0 if provided
- `price_max` >= `price_min` if both provided
- `age_restriction` must be one of: "All Ages", "18+", "21+" (if provided)

**Response:**

```typescript
{
  data: {
    id: number;
    status: string;              // Always "pending" for submissions
    created_at: string;
  };
}
```

**Implementation Notes:**
- Set `status = 'pending'`
- Set `source = 'band_submitted'`
- For each band in `bands` array:
  - Search for existing band by name (case-insensitive)
  - If not found, create new band with auto-generated slug
  - Create show_bands entry with is_headliner and performance_order

**Errors:**
- `400 VALIDATION_ERROR` - Invalid input, return `details` object with field errors
- `404 NOT_FOUND` - venue_id doesn't exist

---

## Venues Endpoints

### `GET /api/venues`

List all venues.

**Query Parameters:**

```typescript
{
  region?: string[];           // Filter by region(s), repeatable
}
```

**Response:**

```typescript
{
  data: {
    id: number;
    name: string;
    slug: string;
    address: string | null;
    region: string | null;
    capacity: number | null;
    website: string | null;
    image_url: string | null;
    upcoming_show_count: number;  // Count of scheduled future shows
  }[];
}
```

**SQL Notes:**
- LEFT JOIN to shows with WHERE clause: `status='scheduled' AND date >= NOW()`
- COUNT shows and GROUP BY venue
- ORDER BY name ASC

---

### `GET /api/venues/:slug`

Get venue details with upcoming shows.

**Path Parameters:**
- `slug` - Venue slug (string)

**Response:**

```typescript
{
  data: {
    id: number;
    name: string;
    slug: string;
    address: string | null;
    city: string;
    state: string;
    zip_code: string | null;
    region: string | null;
    capacity: number | null;
    website: string | null;
    phone: string | null;
    image_url: string | null;

    upcoming_shows: {
      id: number;
      title: string | null;
      date: string;
      price_min: number | null;
      price_max: number | null;
      bands: {
        id: number;
        name: string;
        slug: string;
        is_headliner: boolean;
      }[];
    }[];
  };
}
```

**Errors:**
- `404 NOT_FOUND` - Venue slug doesn't exist

**SQL Notes:**
- Find venue by slug
- LEFT JOIN shows WHERE `status='scheduled' AND date >= NOW()`
- For each show, JOIN show_bands → bands
- ORDER shows by date ASC

---

## Bands Endpoints

### `GET /api/bands`

List bands with pagination.

**Query Parameters:**

```typescript
{
  page?: number;               // Default: 1
  per_page?: number;           // Default: 50, Max: 100
  genre?: string[];            // Filter by genre slug(s), repeatable
  q?: string;                  // Search by name
}
```

**Response:**

```typescript
{
  data: {
    id: number;
    name: string;
    slug: string;
    bio: string | null;
    hometown: string | null;
    image_url: string | null;
    genres: {
      id: number;
      name: string;
      slug: string;
    }[];
  }[];

  meta: {
    page: number;
    per_page: number;
    total: number;
    total_pages: number;
  };
}
```

**SQL Notes:**
- JOIN band_genres → genres
- If `genre` param: filter by genre.slug IN (...)
- If `q` param: WHERE to_tsvector('english', name) @@ plainto_tsquery('english', q)
- ORDER BY name ASC

---

### `GET /api/bands/:slug`

Get band details with upcoming shows and genres.

**Path Parameters:**
- `slug` - Band slug (string)

**Response:**

```typescript
{
  data: {
    id: number;
    name: string;
    slug: string;
    bio: string | null;
    hometown: string | null;
    image_url: string | null;
    website: string | null;
    spotify_url: string | null;
    instagram: string | null;
    facebook: string | null;
    bandcamp_url: string | null;

    genres: {
      id: number;
      name: string;
      slug: string;
    }[];

    upcoming_shows: {
      id: number;
      date: string;
      venue: {
        id: number;
        name: string;
        slug: string;
      };
      is_headliner: boolean;
    }[];
  };
}
```

**Errors:**
- `404 NOT_FOUND` - Band slug doesn't exist

**SQL Notes:**
- Find band by slug
- JOIN band_genres → genres
- JOIN show_bands → shows → venues WHERE `shows.status='scheduled' AND shows.date >= NOW()`
- ORDER shows by date ASC

---

### `GET /api/bands/:slug/similar`

Get similar bands based on shared genres.

**Path Parameters:**
- `slug` - Band slug (string)

**Query Parameters:**

```typescript
{
  limit?: number;              // Default: 10, Max: 50
}
```

**Response:**

```typescript
{
  data: {
    id: number;
    name: string;
    slug: string;
    image_url: string | null;
    shared_genre_count: number;
    shared_genres: {
      id: number;
      name: string;
      slug: string;
    }[];
  }[];
}
```

**Errors:**
- `404 NOT_FOUND` - Band slug doesn't exist

**SQL Notes:**
- Find bands that share genres with the target band
- Self-join on band_genres: `bg1.genre_id = bg2.genre_id WHERE bg2.band_id = target_band_id`
- COUNT shared genres, aggregate genre list
- ORDER BY shared_genre_count DESC, name ASC
- LIMIT to requested limit

---

## Genres Endpoints

### `GET /api/genres`

List all genres.

**Response:**

```typescript
{
  data: {
    id: number;
    name: string;
    slug: string;
    description: string | null;
    show_count?: number;         // Optional: count of upcoming shows
  }[];
}
```

**SQL Notes:**
- SELECT all genres
- Optional: LEFT JOIN to count shows via band_genres → show_bands → shows WHERE scheduled + future
- ORDER BY name ASC

---

## Search Endpoint

### `GET /api/search`

Global search across shows, bands, and venues.

**Query Parameters:**

```typescript
{
  q: string;                   // Required, minimum 2 characters
  limit?: number;              // Default: 20, Max: 50 (applied per entity type)
}
```

**Validation:**
- `q` required, minimum length 2 characters
- Return `400 MISSING_PARAMETER` if q not provided
- Return `400 INVALID_PARAMETER` if q.length < 2

**Response:**

```typescript
{
  data: {
    shows: {
      id: number;
      title: string | null;
      date: string;
      venue_name: string;
    }[];

    bands: {
      id: number;
      name: string;
      slug: string;
    }[];

    venues: {
      id: number;
      name: string;
      slug: string;
    }[];
  };
}
```

**SQL Notes:**
- Search shows: `to_tsvector('english', title) @@ plainto_tsquery('english', q)` WHERE scheduled + future
- Search bands: `to_tsvector('english', name || ' ' || COALESCE(bio, '')) @@ plainto_tsquery('english', q)`
- Search venues: `to_tsvector('english', name) @@ plainto_tsquery('english', q)`
- LIMIT each query to `limit` parameter
- Return empty arrays if no matches

---

## Health Endpoint

### `GET /health`

Health check for monitoring.

**Response:**

```typescript
{
  status: "ok";
  timestamp: string;           // ISO 8601
  database: "connected" | "disconnected";
}
```

**Implementation:**
- Test database connection with simple query: `SELECT 1`
- If query succeeds: `database: "connected"`
- If query fails: `database: "disconnected"`, return HTTP 503
- Always include current timestamp

---

## Implementation Notes

### Date Handling
- Store all dates in database with timezone: `TIMESTAMP WITH TIME ZONE`
- Return all dates as ISO 8601 strings with timezone offset
- Parse incoming dates as ISO 8601

### Slug Generation
- Generate slugs from names: lowercase, replace spaces with hyphens, remove special chars
- Example: "The Orange Peel" → "the-orange-peel"
- Ensure uniqueness by appending `-2`, `-3` if collision

### CORS
- Enable CORS for frontend domain
- Allow methods: GET, POST, OPTIONS
- Allow headers: Content-Type

### Performance
- Use indexes on frequently queried fields (see database-schema.md)
- Eager load relationships to avoid N+1 queries
- Use `COUNT(*) OVER()` for pagination total without separate query

### Null Handling
- Return `null` for optional fields that don't have values
- Never omit fields from response (always include with `null`)
