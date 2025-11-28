import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Show } from "@/lib/types";
import { format, parseISO } from "date-fns";
import Link from "next/link";

interface ShowCardProps {
  show: Show;
}

export function ShowCard({ show }: ShowCardProps) {
  const showDate = parseISO(show.date);
  const headliner = show.bands?.find((b) => b.is_headliner);
  const openers = show.bands?.filter((b) => !b.is_headliner) || [];

  return (
    <Card className="hover:shadow-lg transition-shadow">
      <CardHeader>
        <div className="flex items-start justify-between gap-4">
          <div className="space-y-1">
            <div className="text-sm text-muted-foreground">
              {format(showDate, "EEEE, MMMM d, yyyy")}
            </div>
            <h3 className="text-xl font-bold">
              {headliner?.band?.name || show.title}
            </h3>
            {openers.length > 0 && (
              <p className="text-sm text-muted-foreground">
                with {openers.map((b) => b.band?.name).join(", ")}
              </p>
            )}
          </div>
          <div className="text-right">
            <div className="text-xs text-muted-foreground">
              {format(showDate, "h:mm a")}
            </div>
          </div>
        </div>
      </CardHeader>

      <CardContent>
        <div className="space-y-3">
          <div className="flex items-center justify-between text-sm">
            <div>
              <div className="font-medium">{show.venue?.name}</div>
              {show.venue?.region && (
                <div className="text-xs text-muted-foreground capitalize">
                  {show.venue.region}
                </div>
              )}
            </div>
            {(show.price_min || show.price_max) && (
              <div className="font-medium">
                {show.price_min === 0 ? (
                  "Free"
                ) : show.price_min === show.price_max ? (
                  `$${show.price_min}`
                ) : (
                  `$${show.price_min}-${show.price_max}`
                )}
              </div>
            )}
          </div>

          <div className="flex gap-2">
            <Button asChild variant="default" size="sm" className="flex-1">
              <Link href={`/shows/${show.id}`}>View Details</Link>
            </Button>
            {show.ticket_url && (
              <Button asChild variant="outline" size="sm">
                <a
                  href={show.ticket_url}
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  Tickets
                </a>
              </Button>
            )}
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
