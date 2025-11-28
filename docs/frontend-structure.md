# Frontend Structure

## Tech Stack

- **Next.js 15** - App Router, TypeScript strict mode
- **Tailwind CSS** - Utility-first styling
- **shadcn/ui** - Copy-paste component library
- **TanStack Query** - Server state management
- **React Hook Form + Zod** - Form handling and validation
- **Lucide React** - Icons

---

## Pages & Routes

### Public Pages

```
/                           - Homepage (featured shows, quick links)
/shows                      - Shows listing with filters
/shows/[id]                 - Show detail page
/shows/submit               - Band submission form
/venues                     - Venues listing
/venues/[slug]              - Venue detail page
/bands                      - Bands listing (paginated)
/bands/[slug]               - Band detail page
/search                     - Search results page
```

### Rendering Strategy

| Route | Strategy | Revalidate | Reason |
|-------|----------|------------|--------|
| `/` | SSG + ISR | 3600s (1hr) | Static content, refresh hourly for featured shows |
| `/shows` | SSR | N/A | Query params for filters, needs fresh data |
| `/shows/[id]` | SSR | N/A | Dynamic, needs latest ticket/status info |
| `/shows/submit` | CSR | N/A | Form with client-side validation |
| `/venues` | SSG + ISR | 86400s (1day) | Venue list changes rarely |
| `/venues/[slug]` | SSG + ISR | 21600s (6hr) | Static venue info, upcoming shows refresh 4x/day |
| `/bands` | SSR | N/A | Pagination and filters |
| `/bands/[slug]` | SSG + ISR | 86400s (1day) | Band info static, upcoming shows refresh daily |
| `/search` | SSR | N/A | Dynamic search query |

**Implementation:**
- Add `export const revalidate = <seconds>` to enable ISR
- No revalidate export = SSR by default
- Use `'use client'` directive only when needed (forms, interactivity)

---

## State Management

### 1. Server State (TanStack Query)

Use for all API data fetching and mutations.

```typescript
// Example: Fetching shows
'use client';
import { useQuery } from '@tanstack/react-query';

function ShowsList({ filters }: { filters: ShowFilters }) {
  const { data, isLoading } = useQuery({
    queryKey: ['shows', filters],
    queryFn: () => fetchShows(filters),
  });

  return <div>{/* render shows */}</div>;
}
```

**When to use:**
- Fetching data from API
- Mutations (POST /api/shows)
- Automatic caching and refetching
- Loading/error states

### 2. URL State (Query Parameters)

All filter state lives in URL.

```typescript
// Example: Filter state from URL
const searchParams = useSearchParams();
const genre = searchParams.getAll('genre'); // ['rock', 'indie']
const dateFrom = searchParams.get('date_from');
```

**Stored in URL:**
- Show filters (genre, venue, date range, price)
- Pagination (page, per_page)
- Search query (q)
- Sort order

**Benefits:**
- Shareable URLs
- Browser back/forward works
- No need for local storage
- State persists on refresh

### 3. Local State (useState)

Only for UI-only state.

```typescript
// Mobile filter drawer open/closed
const [isFilterOpen, setIsFilterOpen] = useState(false);

// Form input state (managed by React Hook Form)
const form = useForm<ShowSubmission>();
```

**Use for:**
- Modal/drawer open state
- Form inputs (via React Hook Form)
- Accordion expand/collapse
- Temporary UI state

### 4. React Context (Minimal)

Avoid unless necessary. Only for truly global state.

```typescript
// Example: Theme provider (if adding dark mode later)
<ThemeProvider>
  {children}
</ThemeProvider>
```

**Avoid for:**
- API data (use TanStack Query)
- Filter state (use URL)
- Form state (use React Hook Form)

---

## Component Architecture

### Page Structure Pattern

```
app/shows/page.tsx                    (Server Component - default)
├── Fetch data from API
├── Parse searchParams for filters
└── Return JSX with data

app/shows/ShowsPageClient.tsx         (Client Component - 'use client')
├── FilterSidebar (Client)
│   ├── DateRangePicker
│   ├── VenueMultiSelect
│   ├── GenreMultiSelect
│   ├── PriceRangeSlider
│   └── ClearFiltersButton
├── ShowsGrid (Server Component can be nested)
│   └── ShowCard[] (Server Component)
├── Pagination (Client - updates URL)
└── MobileFilterDrawer (Client)
    └── (same filters as sidebar)
```

**Key Principles:**
- Server Components by default (better performance)
- Use `'use client'` only when needed (interactivity, hooks)
- Server Components can render Client Components
- Client Components can't render Server Components (but can pass as children)

### Component Organization

```
components/
├── ui/                     # shadcn/ui components
│   ├── button.tsx
│   ├── card.tsx
│   ├── input.tsx
│   └── ...
├── shows/
│   ├── ShowCard.tsx        # Server Component
│   ├── ShowsGrid.tsx       # Server Component
│   ├── FilterSidebar.tsx   # Client Component
│   └── ShowSubmissionForm.tsx  # Client Component
├── venues/
│   ├── VenueCard.tsx
│   └── VenueMap.tsx        # Client (if using map library)
├── bands/
│   ├── BandCard.tsx
│   └── SimilarBands.tsx
├── shared/
│   ├── Header.tsx
│   ├── Footer.tsx
│   ├── Pagination.tsx
│   └── SearchBar.tsx
└── layout/
    ├── PageHeader.tsx
    └── Container.tsx
```

---

## Filter UI Behavior

### Desktop (>= 768px)

**Layout:**
- Sidebar on left (sticky)
- Main content on right
- Filters always visible

**Interaction:**
- Filters apply immediately on change
- No "Apply" button needed
- URL updates on each filter change
- Page re-fetches data automatically (TanStack Query watches URL)

**Example Flow:**
1. User selects "Rock" genre
2. URL updates: `/shows?genre=rock`
3. TanStack Query detects URL change
4. New API request: `GET /api/shows?genre=rock`
5. Results update automatically

### Mobile (< 768px)

**Layout:**
- Floating "Filters" button (bottom-right, shows active count)
- Full-screen drawer when opened

**Interaction:**
- Click "Filters" → Drawer slides up
- User adjusts filters (local state)
- Click "Apply" → Update URL + close drawer
- Click "Clear" → Reset filters + close drawer
- Click outside → Close without applying

**Why different behavior?**
- Desktop: Immediate feedback, space for sidebar
- Mobile: Batch changes to avoid constant re-renders, save screen space

### Filter Persistence

**All filter state in URL query params:**
```
/shows?genre=rock&genre=indie&venue=orange-peel&date_from=2025-11-01&price_max=30&page=2
```

**No localStorage** - URL is single source of truth

**Benefits:**
- Share filtered URL with friends
- Bookmark filtered views
- Browser back/forward works correctly
- No sync issues between tabs

### Clear Filters

**"Clear all" button:**
- Removes all query params
- Navigates to `/shows` (no params)
- Resets to default view (all upcoming shows)

---

## Data Fetching Patterns

### Server Component (SSR/SSG)

```typescript
// app/venues/[slug]/page.tsx
export const revalidate = 21600; // ISR: 6 hours

export default async function VenuePage({ params }: { params: { slug: string } }) {
  // Fetch directly in component
  const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/venues/${params.slug}`, {
    next: { revalidate: 21600 }
  });

  if (!response.ok) {
    notFound(); // Show 404 page
  }

  const { data: venue } = await response.json();

  return <VenueDetails venue={venue} />;
}
```

### Client Component (TanStack Query)

```typescript
// components/shows/ShowsList.tsx
'use client';

import { useQuery } from '@tanstack/react-query';
import { useSearchParams } from 'next/navigation';

export function ShowsList() {
  const searchParams = useSearchParams();

  const { data, isLoading, error } = useQuery({
    queryKey: ['shows', searchParams.toString()],
    queryFn: () => fetchShows(searchParams),
  });

  if (isLoading) return <ShowsGridSkeleton />;
  if (error) return <ErrorMessage error={error} />;

  return <ShowsGrid shows={data.data} meta={data.meta} />;
}
```

### API Client

```typescript
// lib/api.ts
const API_URL = process.env.NEXT_PUBLIC_API_URL;

export async function fetchShows(params: URLSearchParams) {
  const response = await fetch(`${API_URL}/api/shows?${params}`);

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error.message);
  }

  return response.json();
}
```

---

## Form Handling

### Show Submission Form

```typescript
// app/shows/submit/page.tsx
'use client';

import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { useMutation } from '@tanstack/react-query';

const showSchema = z.object({
  venue_id: z.number(),
  date: z.string().datetime(),
  bands: z.array(z.object({
    name: z.string().min(1, 'Band name required'),
    is_headliner: z.boolean().optional(),
  })).min(1, 'At least one band required'),
  price_min: z.number().min(0).optional(),
});

type ShowSubmission = z.infer<typeof showSchema>;

export default function SubmitShowPage() {
  const form = useForm<ShowSubmission>({
    resolver: zodResolver(showSchema),
  });

  const mutation = useMutation({
    mutationFn: (data: ShowSubmission) =>
      fetch('/api/shows', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      }),
    onSuccess: () => {
      // Show success message
      // Redirect to show page
    },
  });

  return (
    <form onSubmit={form.handleSubmit(data => mutation.mutate(data))}>
      {/* form fields */}
    </form>
  );
}
```

---

## URL Structure & Dynamic Routes

### Route Parameters

```typescript
// Dynamic segments use [param] syntax

// app/shows/[id]/page.tsx
export default function ShowPage({ params }: { params: { id: string } }) {
  const showId = params.id;
}

// app/venues/[slug]/page.tsx
export default function VenuePage({ params }: { params: { slug: string } }) {
  const venueSlug = params.slug;
}

// app/bands/[slug]/page.tsx
export default function BandPage({ params }: { params: { slug: string } }) {
  const bandSlug = params.slug;
}
```

### Query Parameters (Filters)

```typescript
// app/shows/page.tsx
import { redirect } from 'next/navigation';

export default async function ShowsPage({
  searchParams,
}: {
  searchParams: { [key: string]: string | string[] | undefined };
}) {
  // Access query params
  const genre = searchParams.genre; // Can be string or string[]
  const dateFrom = searchParams.date_from as string | undefined;
  const page = parseInt(searchParams.page as string) || 1;

  // Build API URL
  const params = new URLSearchParams();
  if (genre) {
    const genres = Array.isArray(genre) ? genre : [genre];
    genres.forEach(g => params.append('genre', g));
  }
  if (dateFrom) params.set('date_from', dateFrom);
  params.set('page', page.toString());

  const shows = await fetchShows(params);

  return <ShowsPageClient shows={shows} />;
}
```

### Updating URL (Client Component)

```typescript
'use client';

import { useRouter, useSearchParams } from 'next/navigation';

function FilterSidebar() {
  const router = useRouter();
  const searchParams = useSearchParams();

  function updateFilter(key: string, value: string) {
    const params = new URLSearchParams(searchParams);
    params.set(key, value);
    router.push(`/shows?${params.toString()}`);
  }

  function addGenre(genre: string) {
    const params = new URLSearchParams(searchParams);
    params.append('genre', genre); // Append for multiple values
    router.push(`/shows?${params.toString()}`);
  }
}
```

---

## SEO & Metadata

### Dynamic Metadata

```typescript
// app/shows/[id]/page.tsx
import type { Metadata } from 'next';

export async function generateMetadata({ params }: { params: { id: string } }): Promise<Metadata> {
  const show = await fetchShow(params.id);

  return {
    title: `${show.title} - ${show.venue.name}`,
    description: `${show.title} at ${show.venue.name} on ${formatDate(show.date)}`,
    openGraph: {
      title: show.title,
      description: `Live at ${show.venue.name}`,
      images: [show.image_url || show.venue.image_url],
    },
  };
}
```

### Structured Data (JSON-LD)

```typescript
// Add to show detail page for Google search results
<script type="application/ld+json">
  {JSON.stringify({
    "@context": "https://schema.org",
    "@type": "MusicEvent",
    "name": show.title,
    "startDate": show.date,
    "location": {
      "@type": "Place",
      "name": show.venue.name,
      "address": show.venue.address,
    },
    "performer": show.bands.map(b => ({
      "@type": "MusicGroup",
      "name": b.name,
    })),
  })}
</script>
```

---

## Image Handling

### Display Priority (Fallback Logic)

```typescript
function getShowImage(show: Show): string {
  return (
    show.image_url ||                           // 1. Event poster
    show.bands.find(b => b.is_headliner)?.image_url ||  // 2. Headliner
    show.venue.image_url ||                     // 3. Venue
    '/placeholder-concert.jpg'                  // 4. Placeholder
  );
}
```

### Next.js Image Component

```typescript
import Image from 'next/image';

<Image
  src={getShowImage(show)}
  alt={show.title}
  width={400}
  height={300}
  className="rounded-lg"
  priority={false}  // Set true for above-fold images
/>
```

**Benefits:**
- Automatic image optimization
- Lazy loading by default
- Responsive images
- WebP format when supported

---

## Error Handling

### Not Found (404)

```typescript
// app/shows/[id]/page.tsx
import { notFound } from 'next/navigation';

export default async function ShowPage({ params }: { params: { id: string } }) {
  const response = await fetch(`/api/shows/${params.id}`);

  if (response.status === 404) {
    notFound(); // Renders app/not-found.tsx
  }

  const { data: show } = await response.json();
  return <ShowDetails show={show} />;
}
```

### Error Boundary

```typescript
// app/shows/error.tsx
'use client';

export default function ShowsError({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  return (
    <div>
      <h2>Something went wrong!</h2>
      <p>{error.message}</p>
      <button onClick={reset}>Try again</button>
    </div>
  );
}
```

---

## Loading States

### Page-level Loading

```typescript
// app/shows/loading.tsx
export default function ShowsLoading() {
  return <ShowsGridSkeleton />;
}
```

### Component-level Loading (TanStack Query)

```typescript
const { data, isLoading } = useQuery({ ... });

if (isLoading) return <Skeleton />;
```

---

## Performance Optimizations

### 1. Server Components by Default
- Faster initial page load
- Less JavaScript sent to client
- Better SEO

### 2. Code Splitting
- Next.js automatically splits by route
- Use dynamic imports for heavy components:
  ```typescript
  const Map = dynamic(() => import('./Map'), { ssr: false });
  ```

### 3. Image Optimization
- Use Next.js `<Image>` component
- Lazy load images below fold
- Proper width/height to prevent layout shift

### 4. Caching Strategy
- ISR for static content (venues, bands)
- TanStack Query caching for client-side
- Stale-while-revalidate pattern

### 5. Prefetching
- Next.js prefetches visible links automatically
- Can disable with `prefetch={false}` if needed

---

## Implementation Checklist

### Initial Setup
- [ ] Initialize Next.js 15 with TypeScript
- [ ] Install dependencies (Tailwind, shadcn/ui, TanStack Query)
- [ ] Configure TanStack Query provider in root layout
- [ ] Set up API base URL in environment variables

### Layout & Navigation
- [ ] Create root layout with header/footer
- [ ] Implement navigation component
- [ ] Add mobile menu

### Pages (in order)
- [ ] Homepage (SSG + ISR)
- [ ] Shows listing (SSR)
- [ ] Show detail (SSR)
- [ ] Venues listing (SSG + ISR)
- [ ] Venue detail (SSG + ISR)
- [ ] Bands listing (SSR)
- [ ] Band detail (SSG + ISR)
- [ ] Search (SSR)
- [ ] Show submission form (CSR)

### Components
- [ ] ShowCard component
- [ ] Filter components (sidebar + mobile drawer)
- [ ] Pagination component
- [ ] VenueCard component
- [ ] BandCard component
- [ ] SearchBar component

### Features
- [ ] Filter state in URL
- [ ] Form validation with Zod
- [ ] Image fallback logic
- [ ] Error boundaries
- [ ] Loading skeletons
- [ ] SEO metadata
- [ ] Structured data (JSON-LD)
