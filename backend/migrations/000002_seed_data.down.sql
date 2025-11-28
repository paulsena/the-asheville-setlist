-- The Asheville Setlist - Seed Data Rollback
-- Removes all seeded genres and venues

BEGIN;

-- Delete venues (order doesn't matter due to ON DELETE CASCADE on related tables)
DELETE FROM venues WHERE slug IN (
  'the-orange-peel', 'asheville-yards', 'harrahs-cherokee-center', 'exploreasheville-arena',
  'the-grey-eagle', 'salvage-station', 'asheville-music-hall', 'the-one-stop',
  'pisgah-brewing', 'sierra-nevada-amphitheater', 'sierra-nevada-high-gravity',
  'the-mothlight', 'eulogy', 'the-double-crown', 'fleetwoods', 'isis-music-hall',
  'sly-grog-lounge', 'third-room', 'lazy-diamond', 'barleys-taproom',
  'highland-brewing', 'new-belgium-brewing', 'wicked-weed-funkatorium', 'burial-beer',
  'zillicoah-beer', 'french-broad-river-brewery', 'mills-river-brewing',
  'hotel-eve', 'sovereign-kava', 'the-getaway', 'white-horse-black-mountain',
  'allgood-coffee-weaverville', 'dripolator-coffeehouse'
);

-- Delete genres (order doesn't matter due to ON DELETE CASCADE on related tables)
DELETE FROM genres WHERE slug IN (
  'rock', 'indie', 'alternative', 'punk', 'metal', 'emo',
  'bluegrass', 'americana', 'folk', 'country', 'singer-songwriter',
  'jazz', 'blues', 'soul', 'rnb', 'funk',
  'electronic', 'dj', 'hip-hop',
  'jam-band', 'psychedelic', 'prog-rock',
  'world', 'reggae', 'latin',
  'classical', 'gospel', 'motown', 'cover-band', 'comedy', 'open-mic'
);

COMMIT;
