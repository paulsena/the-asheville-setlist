import Link from "next/link";

export function Header() {
  return (
    <header className="border-b bg-background">
      <div className="container mx-auto flex h-16 items-center justify-between px-4">
        <Link href="/" className="text-xl font-bold">
          The Asheville Setlist
        </Link>

        <nav aria-label="Main navigation" className="flex items-center gap-6">
          <Link
            href="/shows"
            className="text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
          >
            Shows
          </Link>
          <Link
            href="/venues"
            className="text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
          >
            Venues
          </Link>
          <Link
            href="/bands"
            className="text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
          >
            Bands
          </Link>
        </nav>
      </div>
    </header>
  );
}
