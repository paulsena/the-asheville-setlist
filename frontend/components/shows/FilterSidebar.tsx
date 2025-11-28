"use client";

import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { Checkbox } from "@/components/ui/checkbox";
import { Slider } from "@/components/ui/slider";
import { Button } from "@/components/ui/button";
import { useState } from "react";

interface FilterSidebarProps {
  onFilterChange?: (filters: ShowFilters) => void;
}

export interface ShowFilters {
  venues: string[];
  genres: string[];
  regions: string[];
  priceRange: [number, number];
  dateFrom?: string;
  dateTo?: string;
}

// Mock data - will come from API later
const VENUES = [
  { id: "orange-peel", name: "The Orange Peel" },
  { id: "grey-eagle", name: "The Grey Eagle" },
  { id: "salvage-station", name: "Salvage Station" },
  { id: "mothlight", name: "The Mothlight" },
  { id: "asheville-music-hall", name: "Asheville Music Hall" },
];

const GENRES = [
  { id: "rock", name: "Rock" },
  { id: "indie", name: "Indie" },
  { id: "bluegrass", name: "Bluegrass" },
  { id: "folk", name: "Folk" },
  { id: "jazz", name: "Jazz" },
  { id: "electronic", name: "Electronic" },
  { id: "hip-hop", name: "Hip Hop" },
];

const REGIONS = [
  { id: "downtown", name: "Downtown" },
  { id: "west-asheville", name: "West Asheville" },
  { id: "river-arts", name: "River Arts District" },
  { id: "south-asheville", name: "South Asheville" },
];

export function FilterSidebar({ onFilterChange }: FilterSidebarProps) {
  const [selectedVenues, setSelectedVenues] = useState<string[]>([]);
  const [selectedGenres, setSelectedGenres] = useState<string[]>([]);
  const [selectedRegions, setSelectedRegions] = useState<string[]>([]);
  const [priceRange, setPriceRange] = useState<[number, number]>([0, 100]);

  const handleVenueToggle = (venueId: string) => {
    setSelectedVenues((prev) =>
      prev.includes(venueId)
        ? prev.filter((id) => id !== venueId)
        : [...prev, venueId]
    );
  };

  const handleGenreToggle = (genreId: string) => {
    setSelectedGenres((prev) =>
      prev.includes(genreId)
        ? prev.filter((id) => id !== genreId)
        : [...prev, genreId]
    );
  };

  const handleRegionToggle = (regionId: string) => {
    setSelectedRegions((prev) =>
      prev.includes(regionId)
        ? prev.filter((id) => id !== regionId)
        : [...prev, regionId]
    );
  };

  const handleClearFilters = () => {
    setSelectedVenues([]);
    setSelectedGenres([]);
    setSelectedRegions([]);
    setPriceRange([0, 100]);
  };

  const handleApplyFilters = () => {
    onFilterChange?.({
      venues: selectedVenues,
      genres: selectedGenres,
      regions: selectedRegions,
      priceRange,
    });
  };

  const hasActiveFilters =
    selectedVenues.length > 0 ||
    selectedGenres.length > 0 ||
    selectedRegions.length > 0 ||
    priceRange[0] > 0 ||
    priceRange[1] < 100;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h2 className="text-lg font-semibold">Filters</h2>
        {hasActiveFilters && (
          <Button
            variant="ghost"
            size="sm"
            onClick={handleClearFilters}
            className="h-auto p-0 text-sm text-muted-foreground hover:text-foreground"
          >
            Clear all
          </Button>
        )}
      </div>

      <Separator />

      {/* Regions */}
      <div className="space-y-3">
        <Label className="text-sm font-medium">Region</Label>
        <div className="space-y-2">
          {REGIONS.map((region) => (
            <div key={region.id} className="flex items-center space-x-2">
              <Checkbox
                id={`region-${region.id}`}
                checked={selectedRegions.includes(region.id)}
                onCheckedChange={() => handleRegionToggle(region.id)}
              />
              <label
                htmlFor={`region-${region.id}`}
                className="text-sm leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
              >
                {region.name}
              </label>
            </div>
          ))}
        </div>
      </div>

      <Separator />

      {/* Venues */}
      <div className="space-y-3">
        <Label className="text-sm font-medium">Venue</Label>
        <div className="space-y-2">
          {VENUES.map((venue) => (
            <div key={venue.id} className="flex items-center space-x-2">
              <Checkbox
                id={`venue-${venue.id}`}
                checked={selectedVenues.includes(venue.id)}
                onCheckedChange={() => handleVenueToggle(venue.id)}
              />
              <label
                htmlFor={`venue-${venue.id}`}
                className="text-sm leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
              >
                {venue.name}
              </label>
            </div>
          ))}
        </div>
      </div>

      <Separator />

      {/* Genres */}
      <div className="space-y-3">
        <Label className="text-sm font-medium">Genre</Label>
        <div className="space-y-2">
          {GENRES.map((genre) => (
            <div key={genre.id} className="flex items-center space-x-2">
              <Checkbox
                id={`genre-${genre.id}`}
                checked={selectedGenres.includes(genre.id)}
                onCheckedChange={() => handleGenreToggle(genre.id)}
              />
              <label
                htmlFor={`genre-${genre.id}`}
                className="text-sm leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
              >
                {genre.name}
              </label>
            </div>
          ))}
        </div>
      </div>

      <Separator />

      {/* Price Range */}
      <div className="space-y-3">
        <div className="flex items-center justify-between">
          <Label className="text-sm font-medium">Price Range</Label>
          <span className="text-sm text-muted-foreground">
            ${priceRange[0]} - ${priceRange[1]}
          </span>
        </div>
        <Slider
          value={priceRange}
          onValueChange={(value) => setPriceRange(value as [number, number])}
          max={100}
          step={5}
          className="w-full"
        />
        <div className="flex justify-between text-xs text-muted-foreground">
          <span>$0</span>
          <span>$100+</span>
        </div>
      </div>

      <Separator />

      {/* Apply Button */}
      <Button onClick={handleApplyFilters} className="w-full">
        Apply Filters
      </Button>
    </div>
  );
}
