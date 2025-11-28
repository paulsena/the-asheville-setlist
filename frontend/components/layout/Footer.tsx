import Link from "next/link";

export function Footer() {
  const currentYear = new Date().getFullYear();

  return (
    <footer className="border-t bg-background">
      <div className="container mx-auto px-4 py-8">
        <div className="grid grid-cols-1 gap-8 md:grid-cols-3">
          <div>
            <h3 className="mb-3 text-sm font-semibold">Discover</h3>
            <ul className="space-y-2 text-sm">
              <li>
                <Link
                  href="/shows"
                  className="text-muted-foreground hover:text-foreground"
                >
                  Upcoming Shows
                </Link>
              </li>
              <li>
                <Link
                  href="/venues"
                  className="text-muted-foreground hover:text-foreground"
                >
                  Venues
                </Link>
              </li>
              <li>
                <Link
                  href="/bands"
                  className="text-muted-foreground hover:text-foreground"
                >
                  Bands
                </Link>
              </li>
            </ul>
          </div>

          <div>
            <h3 className="mb-3 text-sm font-semibold">For Bands</h3>
            <ul className="space-y-2 text-sm">
              <li>
                <Link
                  href="/shows/submit"
                  className="text-muted-foreground hover:text-foreground"
                >
                  Submit a Show
                </Link>
              </li>
            </ul>
          </div>

          <div>
            <h3 className="mb-3 text-sm font-semibold">About</h3>
            <p className="text-sm text-muted-foreground">
              The Asheville Setlist helps music lovers discover concerts and
              live music in Asheville, NC.
            </p>
          </div>
        </div>

        <div className="mt-8 border-t pt-6 text-center text-sm text-muted-foreground">
          © {currentYear} The Asheville Setlist. All rights reserved.
        </div>
      </div>
    </footer>
  );
}
