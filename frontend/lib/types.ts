/**
 * Type definitions for API responses
 */

export interface Show {
  id: number;
  title: string;
  date: string;
  venue_id: number;
  venue?: Venue;
  bands?: BandShow[];
  price_min?: number;
  price_max?: number;
  ticket_url?: string;
  status: 'scheduled' | 'cancelled' | 'postponed';
  created_at: string;
}

export interface Venue {
  id: number;
  name: string;
  slug: string;
  address?: string;
  city: string;
  state: string;
  region?: string;
  capacity?: number;
  website_url?: string;
}

export interface Band {
  id: number;
  name: string;
  slug: string;
  bio?: string;
  image_url?: string;
  website_url?: string;
  spotify_url?: string;
  bandcamp_url?: string;
}

export interface BandShow {
  band_id: number;
  band?: Band;
  is_headliner: boolean;
  performance_order: number;
}

export interface Genre {
  id: number;
  name: string;
  slug: string;
  description?: string;
}
