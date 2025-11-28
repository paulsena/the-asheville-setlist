import { ShowCard } from "@/components/shows/ShowCard";
import { Button } from "@/components/ui/button";
import { Show } from "@/lib/types";
import Link from "next/link";

// Mock data for testing
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
];

export default function HomePage() {
  return (
    <div className="container mx-auto px-4 py-12">
      {/* Hero Section */}
      <section className="mb-12 text-center">
        <h1 className="mb-4 text-4xl font-bold tracking-tight sm:text-5xl md:text-6xl">
          Discover Live Music in Asheville
        </h1>
        <p className="mx-auto mb-8 max-w-2xl text-lg text-muted-foreground">
          Your complete guide to upcoming concerts, local venues, and the
          vibrant music scene in Asheville, NC.
        </p>
        <div className="flex justify-center gap-4">
          <Button asChild size="lg">
            <Link href="/shows">Browse All Shows</Link>
          </Button>
          <Button asChild variant="outline" size="lg">
            <Link href="/venues">Explore Venues</Link>
          </Button>
        </div>
      </section>

      {/* Featured Shows */}
      <section>
        <div className="mb-6 flex items-center justify-between">
          <h2 className="text-3xl font-bold">Upcoming Shows</h2>
          <Button asChild variant="ghost">
            <Link href="/shows">View All →</Link>
          </Button>
        </div>

        <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {mockShows.map((show) => (
            <ShowCard key={show.id} show={show} />
          ))}
        </div>
      </section>

      {/* Quick Links */}
      <section className="mt-16 grid gap-6 sm:grid-cols-3">
        <div className="rounded-lg border bg-card p-6">
          <h3 className="mb-2 text-lg font-semibold">Find Your Vibe</h3>
          <p className="mb-4 text-sm text-muted-foreground">
            Filter shows by genre, venue, date, and price to find exactly what
            you&apos;re looking for.
          </p>
          <Button asChild variant="outline" size="sm">
            <Link href="/shows">Search Shows</Link>
          </Button>
        </div>

        <div className="rounded-lg border bg-card p-6">
          <h3 className="mb-2 text-lg font-semibold">Explore Venues</h3>
          <p className="mb-4 text-sm text-muted-foreground">
            From intimate clubs to large concert halls, discover Asheville&apos;s
            diverse music venues.
          </p>
          <Button asChild variant="outline" size="sm">
            <Link href="/venues">Browse Venues</Link>
          </Button>
        </div>

        <div className="rounded-lg border bg-card p-6">
          <h3 className="mb-2 text-lg font-semibold">For Bands</h3>
          <p className="mb-4 text-sm text-muted-foreground">
            Playing in Asheville? Submit your show to reach more music lovers.
          </p>
          <Button asChild variant="outline" size="sm">
            <Link href="/shows/submit">Submit a Show</Link>
          </Button>
        </div>
      </section>
    </div>
  );
}
