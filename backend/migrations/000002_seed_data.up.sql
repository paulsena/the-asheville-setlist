-- The Asheville Setlist - Seed Data Migration
-- Loads initial genres and venues

BEGIN;

-- ============================================
-- GENRES
-- ============================================

INSERT INTO genres (name, slug, description) VALUES
-- Rock & Alternative
('Rock', 'rock', 'Rock music including classic rock, hard rock, and modern rock'),
('Indie', 'indie', 'Independent and alternative rock music'),
('Alternative', 'alternative', 'Alternative rock and experimental music'),
('Punk', 'punk', 'Punk rock and hardcore punk'),
('Metal', 'metal', 'Heavy metal and its subgenres'),
('Emo', 'emo', 'Emotional hardcore and emo rock'),

-- Americana & Roots
('Bluegrass', 'bluegrass', 'Traditional and progressive bluegrass'),
('Americana', 'americana', 'American roots music blending country, folk, and blues'),
('Folk', 'folk', 'Traditional and contemporary folk music'),
('Country', 'country', 'Country and country-western music'),
('Singer-Songwriter', 'singer-songwriter', 'Acoustic singer-songwriter performances'),

-- Jazz & Blues
('Jazz', 'jazz', 'Jazz music including traditional, modern, and fusion'),
('Blues', 'blues', 'Blues music and blues rock'),
('Soul', 'soul', 'Soul music and neo-soul'),
('R&B', 'rnb', 'Rhythm and blues music'),
('Funk', 'funk', 'Funk music and funk-influenced artists'),

-- Electronic & Dance
('Electronic', 'electronic', 'Electronic music including EDM, house, and techno'),
('DJ', 'dj', 'DJ sets and electronic dance music'),
('Hip-Hop', 'hip-hop', 'Hip-hop, rap, and urban music'),

-- Jam & Improvisational
('Jam Band', 'jam-band', 'Improvisational jam bands and extended live performances'),
('Psychedelic', 'psychedelic', 'Psychedelic rock and experimental music'),
('Prog Rock', 'prog-rock', 'Progressive rock and art rock'),

-- World & Reggae
('World', 'world', 'World music and international artists'),
('Reggae', 'reggae', 'Reggae, ska, and Caribbean music'),
('Latin', 'latin', 'Latin music including salsa, cumbia, and Latin rock'),

-- Classical & Other
('Classical', 'classical', 'Classical music and orchestral performances'),
('Gospel', 'gospel', 'Gospel and Christian music'),
('Motown', 'motown', 'Motown classics and vintage soul'),
('Cover Band', 'cover-band', 'Cover bands and tribute acts'),
('Comedy', 'comedy', 'Comedy shows and musical comedy'),
('Open Mic', 'open-mic', 'Open mic nights and jam sessions')

ON CONFLICT (slug) DO UPDATE SET
  name = EXCLUDED.name,
  description = EXCLUDED.description;

-- ============================================
-- VENUES
-- ============================================

INSERT INTO venues (name, slug, address, city, state, zip_code, region, capacity, website, metadata) VALUES

-- Major concert venues (500+ capacity)
('The Orange Peel', 'the-orange-peel',
 '101 Biltmore Ave', 'Asheville', 'NC', '28801',
 'downtown', 1050, 'https://theorangepeel.net',
 '{"lma_id": 44986, "parking": "street", "food": false, "bar": true, "type": "concert_hall"}'::jsonb),

('Asheville Yards', 'asheville-yards',
 '75 Coxe Ave', 'Asheville', 'NC', '28801',
 'downtown', 5000, 'https://www.ashevilleyards.com',
 '{"lma_id": null, "formerly": "Rabbit Rabbit", "parking": "lot", "outdoor": true, "type": "amphitheater"}'::jsonb),

('Harrahs Cherokee Center', 'harrahs-cherokee-center',
 '87 Haywood St', 'Asheville', 'NC', '28801',
 'downtown', 7600, 'https://www.harrahscherokeecenterasheville.com',
 '{"lma_id": null, "formerly": "US Cellular Center", "parking": "garage", "type": "arena"}'::jsonb),

('ExploreAsheville.com Arena', 'exploreasheville-arena',
 '87 Haywood St', 'Asheville', 'NC', '28801',
 'downtown', 7600, 'https://www.harrahscherokeecenterasheville.com',
 '{"lma_id": null, "same_as": "harrahs-cherokee-center", "type": "arena"}'::jsonb),

-- Mid-size venues (200-500 capacity)
('The Grey Eagle', 'the-grey-eagle',
 '185 Clingman Ave', 'Asheville', 'NC', '28801',
 'west asheville', 400, 'https://thegreyeagle.com',
 '{"lma_id": null, "parking": "lot", "food": true, "bar": true, "type": "music_hall"}'::jsonb),

('Salvage Station', 'salvage-station',
 '468 Riverside Dr', 'Asheville', 'NC', '28801',
 'west asheville', 1000, 'https://salvagestation.com',
 '{"lma_id": null, "parking": "lot", "outdoor": true, "food": true, "type": "outdoor_venue"}'::jsonb),

('Asheville Music Hall', 'asheville-music-hall',
 '31 Patton Ave', 'Asheville', 'NC', '28801',
 'downtown', 350, 'https://ashevillemusichall.com',
 '{"lma_id": null, "parking": "street", "bar": true, "type": "music_hall"}'::jsonb),

('The One Stop', 'the-one-stop',
 '29 Patton Ave', 'Asheville', 'NC', '28801',
 'downtown', 200, 'https://ashevillemusichall.com',
 '{"lma_id": null, "same_building": "asheville-music-hall", "bar": true, "food": true, "type": "music_hall"}'::jsonb),

('Pisgah Brewing Company', 'pisgah-brewing',
 '150 Eastside Dr', 'Black Mountain', 'NC', '28711',
 'black mountain', 500, 'https://pisgahbrewing.com',
 '{"lma_id": null, "parking": "lot", "outdoor": true, "food": true, "bar": true, "type": "brewery"}'::jsonb),

('Sierra Nevada Amphitheater', 'sierra-nevada-amphitheater',
 '100 Sierra Nevada Way', 'Mills River', 'NC', '28732',
 'mills river', 800, 'https://sierranevada.com/visit/mills-river/amphitheater',
 '{"lma_id": null, "parking": "lot", "outdoor": true, "food": true, "bar": true, "type": "brewery"}'::jsonb),

('Sierra Nevada High Gravity', 'sierra-nevada-high-gravity',
 '100 Sierra Nevada Way', 'Mills River', 'NC', '28732',
 'mills river', 350, 'https://sierranevada.com/visit/mills-river/high-gravity',
 '{"lma_id": null, "parking": "lot", "indoor": true, "food": true, "bar": true, "type": "brewery"}'::jsonb),

-- Small venues & clubs (50-200 capacity)
('The Mothlight', 'the-mothlight',
 '701 Haywood Rd', 'Asheville', 'NC', '28806',
 'west asheville', 150, 'https://themothlight.com',
 '{"lma_id": null, "parking": "street", "bar": true, "type": "club"}'::jsonb),

('Eulogy', 'eulogy',
 '10 Buxton Ave', 'Asheville', 'NC', '28801',
 'downtown', 100, NULL,
 '{"lma_id": null, "parking": "street", "bar": true, "type": "club"}'::jsonb),

('The Double Crown', 'the-double-crown',
 '375 Haywood Rd', 'Asheville', 'NC', '28806',
 'west asheville', 100, NULL,
 '{"lma_id": null, "parking": "street", "bar": true, "type": "bar"}'::jsonb),

('Fleetwoods', 'fleetwoods',
 '496 Haywood Rd', 'Asheville', 'NC', '28806',
 'west asheville', 100, 'https://fleetwoodsavl.com',
 '{"lma_id": null, "parking": "street", "bar": true, "type": "bar"}'::jsonb),

('Isis Music Hall', 'isis-music-hall',
 '743 Haywood Rd', 'Asheville', 'NC', '28806',
 'west asheville', 200, 'https://isisasheville.com',
 '{"lma_id": null, "parking": "street", "bar": true, "type": "music_hall"}'::jsonb),

('Sly Grog Lounge', 'sly-grog-lounge',
 '555 Haywood Rd', 'Asheville', 'NC', '28806',
 'west asheville', 75, NULL,
 '{"lma_id": null, "parking": "street", "bar": true, "type": "lounge"}'::jsonb),

('Third Room', 'third-room',
 '46 Wall St', 'Asheville', 'NC', '28803',
 'downtown', 80, NULL,
 '{"lma_id": null, "parking": "street", "bar": true, "type": "club"}'::jsonb),

('Lazy Diamond', 'lazy-diamond',
 '44 N French Broad Ave', 'Asheville', 'NC', '28801',
 'downtown', 100, NULL,
 '{"lma_id": null, "parking": "street", "bar": true, "type": "bar"}'::jsonb),

('Barleys Taproom', 'barleys-taproom',
 '42 Biltmore Ave', 'Asheville', 'NC', '28801',
 'downtown', 150, 'https://barleystaproom.com',
 '{"lma_id": null, "parking": "street", "bar": true, "type": "bar"}'::jsonb),

-- Breweries & taprooms with music
('Highland Brewing', 'highland-brewing',
 '12 Old Charlotte Hwy Suite H', 'Asheville', 'NC', '28803',
 'south asheville', 300, 'https://highlandbrewing.com',
 '{"lma_id": null, "parking": "lot", "outdoor": true, "food": false, "bar": true, "type": "brewery"}'::jsonb),

('New Belgium Brewing', 'new-belgium-brewing',
 '21 Craven St', 'Asheville', 'NC', '28806',
 'west asheville', 200, 'https://newbelgium.com/visit/asheville',
 '{"lma_id": null, "parking": "lot", "food": true, "bar": true, "type": "brewery"}'::jsonb),

('Wicked Weed Funkatorium', 'wicked-weed-funkatorium',
 '147 Coxe Ave', 'Asheville', 'NC', '28801',
 'south slope', 100, 'https://wickedweedbrewing.com',
 '{"lma_id": null, "parking": "street", "food": true, "bar": true, "type": "brewery"}'::jsonb),

('Burial Beer', 'burial-beer',
 '40 Collier Ave', 'Asheville', 'NC', '28801',
 'south slope', 100, 'https://burialbeer.com',
 '{"lma_id": null, "parking": "street", "food": true, "bar": true, "type": "brewery"}'::jsonb),

('Zillicoah Beer Company', 'zillicoah-beer',
 '870 Riverside Dr', 'Woodfin', 'NC', '28804',
 'north asheville', 150, 'https://zillicoahbeer.com',
 '{"lma_id": null, "parking": "lot", "outdoor": true, "food": true, "bar": true, "type": "brewery"}'::jsonb),

('French Broad River Brewery', 'french-broad-river-brewery',
 '101 Fairview Rd', 'Asheville', 'NC', '28803',
 'east asheville', 150, 'https://frenchbroadrivery.com',
 '{"lma_id": null, "parking": "lot", "outdoor": true, "food": false, "bar": true, "type": "brewery"}'::jsonb),

('Mills River Brewing', 'mills-river-brewing',
 '336 Banner Farm Rd', 'Mills River', 'NC', '28759',
 'mills river', 100, 'https://millsriverbrewery.com',
 '{"lma_id": null, "parking": "lot", "outdoor": true, "type": "brewery"}'::jsonb),

-- Jazz & specialty venues
('Hotel Eve', 'hotel-eve',
 '56 N Lexington Ave', 'Asheville', 'NC', '28801',
 'downtown', 75, 'https://www.hotelevejazz.com',
 '{"lma_id": 63722, "parking": "street", "type": "jazz_club"}'::jsonb),

('Sovereign Kava', 'sovereign-kava',
 '1 Page Ave Suite 135', 'Asheville', 'NC', '28801',
 'downtown', 50, 'https://sovereignkava.com',
 '{"lma_id": null, "parking": "street", "kava_bar": true, "type": "lounge"}'::jsonb),

-- Restaurants & cafes with music
('The Getaway', 'the-getaway',
 '108 N Lexington Ave', 'Asheville', 'NC', '28801',
 'downtown', 60, NULL,
 '{"lma_id": null, "parking": "street", "food": true, "bar": true, "type": "restaurant"}'::jsonb),

('White Horse Black Mountain', 'white-horse-black-mountain',
 '105 Montreat Rd', 'Black Mountain', 'NC', '28711',
 'black mountain', 100, 'https://whitehorseblackmountain.com',
 '{"lma_id": null, "parking": "street", "food": true, "bar": true, "type": "restaurant"}'::jsonb),

-- Coffee shops & listening rooms
('Allgood Coffee Weaverville', 'allgood-coffee-weaverville',
 '12 N Main St', 'Weaverville', 'NC', '28787',
 'weaverville', 40, NULL,
 '{"lma_id": 49759, "parking": "street", "food": true, "type": "coffee_shop"}'::jsonb),

('Dripolator Coffeehouse', 'dripolator-coffeehouse',
 '221 W State St', 'Black Mountain', 'NC', '28711',
 'black mountain', 40, NULL,
 '{"lma_id": null, "parking": "street", "food": true, "type": "coffee_shop"}'::jsonb)

ON CONFLICT (slug) DO UPDATE SET
  name = EXCLUDED.name,
  address = EXCLUDED.address,
  city = EXCLUDED.city,
  state = EXCLUDED.state,
  zip_code = EXCLUDED.zip_code,
  region = EXCLUDED.region,
  capacity = EXCLUDED.capacity,
  website = EXCLUDED.website,
  metadata = venues.metadata || EXCLUDED.metadata,
  updated_at = NOW();

COMMIT;
