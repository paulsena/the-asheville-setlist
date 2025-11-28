-- Seed file: genres
-- Music genre taxonomy for The Asheville Setlist
-- Run: psql $DATABASE_URL < seeds/001_genres.sql

BEGIN;

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

COMMIT;
