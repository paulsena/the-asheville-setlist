import { ShowCard } from "@/components/shows/ShowCard";
import { FilterSidebar, ShowFilters } from "@/components/shows/FilterSidebar";
import { Show } from "@/lib/types";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Shows",
  description: "Browse upcoming concerts and live music shows in Asheville, NC",
};

// Mock data - same as homepage for now
const mockShows: Show[] = [
  {
    id: 1,
    title: "Moon Taxi with Special Guests",
    date: "2025-12-15T20:00:00Z",
    venue_id: 1,
    venue: {
      id: 1,
      name: "The Orange Peel",
      slug: "orange-peel",
      city: "Asheville",
      state: "NC",
      region: "downtown",
    },
    bands: [
      {
        band_id: 1,
        band: { id: 1, name: "Moon Taxi", slug: "moon-taxi" },
        is_headliner: true,
        performance_order: 1,
      },
      {
        band_id: 2,
        band: { id: 2, name: "The Wild Feathers", slug: "wild-feathers" },
        is_headliner: false,
        performance_order: 2,
      },
    ],
    price_min: 25,
    price_max: 30,
    ticket_url: "https://example.com/tickets",
    status: "scheduled",
    created_at: "2025-11-20T10:00:00Z",
  },
  {
    id: 2,
    title: "Greensky Bluegrass",
    date: "2025-12-20T19:30:00Z",
    venue_id: 2,
    venue: {
      id: 2,
      name: "Salvage Station",
      slug: "salvage-station",
      city: "Asheville",
      state: "NC",
      region: "river arts district",
    },
    bands: [
      {
        band_id: 3,
        band: { id: 3, name: "Greensky Bluegrass", slug: "greensky-bluegrass" },
        is_headliner: true,
        performance_order: 1,
      },
    ],
    price_min: 35,
    price_max: 35,
    status: "scheduled",
    created_at: "2025-11-18T10:00:00Z",
  },
  {
    id: 3,
    title: "Local Showcase Night",
    date: "2025-12-10T21:00:00Z",
    venue_id: 3,
    venue: {
      id: 3,
      name: "The Grey Eagle",
      slug: "grey-eagle",
      city: "Asheville",
      state: "NC",
      region: "west asheville",
    },
    bands: [
      {
        band_id: 4,
        band: { id: 4, name: "Mountain Heart", slug: "mountain-heart" },
        is_headliner: true,
        performance_order: 1,
      },
    ],
    price_min: 0,
    price_max: 0,
    status: "scheduled",
    created_at: "2025-11-25T10:00:00Z",
  },
  {
    id: 4,
    title: "Jazz Night at The Mothlight",
    date: "2025-12-08T20:00:00Z",
    venue_id: 4,
    venue: {
      id: 4,
      name: "The Mothlight",
      slug: "mothlight",
      city: "Asheville",
      state: "NC",
      region: "south asheville",
    },
    bands: [
      {
        band_id: 5,
        band: { id: 5, name: "The Jazz Collective", slug: "jazz-collective" },
        is_headliner: true,
        performance_order: 1,
      },
    ],
    price_min: 15,
    price_max: 20,
    status: "scheduled",
    created_at: "2025-11-22T10:00:00Z",
  },
];

export default function ShowsPage() {
  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-8">
        <h1 className="mb-2 text-3xl font-bold">Upcoming Shows</h1>
        <p className="text-muted-foreground">
          Find concerts and live music in Asheville, NC
        </p>
      </div>

      <div className="flex gap-8">
        {/* Sidebar - Desktop Only */}
        <aside className="hidden w-64 shrink-0 lg:block">
          <div className="sticky top-4">
            <FilterSidebar
              onFilterChange={(filters: ShowFilters) => {
                console.log("Filters changed:", filters);
                // In real app, this would update URL params or trigger API refetch
              }}
            />
          </div>
        </aside>

        {/* Main Content */}
        <main className="flex-1">
          <div className="mb-4 flex items-center justify-between">
            <p className="text-sm text-muted-foreground">
              {mockShows.length} shows found
            </p>
            {/* Mobile filter button would go here */}
          </div>

          <div className="grid gap-6 sm:grid-cols-2 xl:grid-cols-2">
            {mockShows.map((show) => (
              <ShowCard key={show.id} show={show} />
            ))}
          </div>

          {/* Pagination would go here */}
        </main>
      </div>
    </div>
  );
}
