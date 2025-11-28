-- Test Data for Development
-- Creates sample bands, shows, and relationships for testing

-- ==============================================================================
-- SAMPLE BANDS
-- ==============================================================================

INSERT INTO bands (name, slug, bio, hometown, image_url, spotify_url, instagram, website) VALUES
-- Rock/Indie
('Moon Taxi', 'moon-taxi', 'Nashville-based indie rock band known for their energetic live shows.', 'Nashville, TN', 'https://example.com/images/moon-taxi.jpg', 'https://open.spotify.com/artist/2ZRQcPiMvJkelhjJjUNw3h', '@moontaxi', 'https://moontaxi.com'),
('The Revivalists', 'the-revivalists', 'New Orleans rock band blending soul, folk, and alternative.', 'New Orleans, LA', 'https://example.com/images/revivalists.jpg', 'https://open.spotify.com/artist/4MzXwWMhyBbmu6hmi498y9', '@therevivalists', 'https://therevivalists.com'),
('Tame Impala', 'tame-impala', 'Australian psychedelic rock project led by Kevin Parker.', 'Perth, Australia', 'https://example.com/images/tame-impala.jpg', 'https://open.spotify.com/artist/5INjqkS1o8h1imAzPqGZBb', '@tameimpala', 'https://tameimpala.com'),

-- Jam Band
('String Cheese Incident', 'string-cheese-incident', 'Colorado-based jam band with electronic influences.', 'Boulder, CO', 'https://example.com/images/sci.jpg', 'https://open.spotify.com/artist/6JvyLnKKx0mlEJlNkPqO4k', '@stringcheeseincident', 'https://stringcheeseincident.com'),
('Umphrey''s McGee', 'umphreys-mcgee', 'Progressive rock jam band from Chicago.', 'Chicago, IL', NULL, 'https://open.spotify.com/artist/3lxhC5BYJLQIHz5rV3Ukvp', '@umphreysmcgee', 'https://umphreys.com'),
('Pigeons Playing Ping Pong', 'pigeons-playing-ping-pong', 'High-energy psychedelic funk jam band.', 'Baltimore, MD', NULL, NULL, '@pigeons_pppp', 'https://pigeonsplayingpingpong.com'),

-- Electronic/Dance
('ODESZA', 'odesza', 'Electronic music duo from Seattle.', 'Seattle, WA', 'https://example.com/images/odesza.jpg', 'https://open.spotify.com/artist/62zFWGJ7hRPZOh9K4ibgGU', '@odesza', 'https://odesza.com'),
('Big Gigantic', 'big-gigantic', 'Live electronic duo combining EDM with live saxophone.', 'Boulder, CO', NULL, NULL, '@biggigantic', 'https://biggigantic.net'),

-- Hip Hop
('Anderson .Paak', 'anderson-paak', 'Grammy-winning rapper, singer, and drummer.', 'Oxnard, CA', 'https://example.com/images/anderson-paak.jpg', 'https://open.spotify.com/artist/3jK9MiCrA42lLAdMGUZpwa', '@anderson._paak', NULL),
('GRiZ', 'griz', 'Electronic producer and saxophonist blending funk and dubstep.', 'Detroit, MI', NULL, NULL, '@griz', 'https://mygriz.com'),

-- Bluegrass/Folk
('Billy Strings', 'billy-strings', 'Progressive bluegrass guitarist and vocalist.', 'Michigan', 'https://example.com/images/billy-strings.jpg', 'https://open.spotify.com/artist/67tgMW4WzcE1ykeQXXiMPx', '@billy_strings', 'https://billystrings.com'),
('Mandolin Orange', 'mandolin-orange', 'Folk duo from Chapel Hill, NC.', 'Chapel Hill, NC', NULL, NULL, '@mandolinorange', NULL),

-- Metal
('Whitechapel', 'whitechapel', 'Deathcore band from Knoxville, Tennessee.', 'Knoxville, TN', NULL, NULL, '@whitechapelband', 'https://whitechapelmetal.com'),
('Bodysnatcher', 'bodysnatcher', 'Heavy deathcore band from Melbourne, Florida.', 'Melbourne, FL', NULL, NULL, '@bodysnatcherhc', NULL),

-- Local Asheville Bands
('The Get Right Band', 'the-get-right-band', 'Asheville-based funk and soul collective.', 'Asheville, NC', NULL, NULL, '@getrightnc', 'https://thegetrightband.com'),
('Town Mountain', 'town-mountain', 'Bluegrass band rooted in Asheville.', 'Asheville, NC', NULL, NULL, '@townmountain', 'https://townmountain.com');

-- ==============================================================================
-- BAND GENRES (Link bands to genres)
-- ==============================================================================

-- Moon Taxi: Rock, Indie, Psychedelic
INSERT INTO band_genres (band_id, genre_id) VALUES
((SELECT id FROM bands WHERE slug = 'moon-taxi'), (SELECT id FROM genres WHERE slug = 'rock')),
((SELECT id FROM bands WHERE slug = 'moon-taxi'), (SELECT id FROM genres WHERE slug = 'indie')),
((SELECT id FROM bands WHERE slug = 'moon-taxi'), (SELECT id FROM genres WHERE slug = 'psychedelic'));

-- The Revivalists: Rock, Indie, Soul
INSERT INTO band_genres (band_id, genre_id) VALUES
((SELECT id FROM bands WHERE slug = 'the-revivalists'), (SELECT id FROM genres WHERE slug = 'rock')),
((SELECT id FROM bands WHERE slug = 'the-revivalists'), (SELECT id FROM genres WHERE slug = 'indie')),
((SELECT id FROM bands WHERE slug = 'the-revivalists'), (SELECT id FROM genres WHERE slug = 'soul'));

-- Tame Impala: Psychedelic, Indie, Electronic
INSERT INTO band_genres (band_id, genre_id) VALUES
((SELECT id FROM bands WHERE slug = 'tame-impala'), (SELECT id FROM genres WHERE slug = 'psychedelic')),
((SELECT id FROM bands WHERE slug = 'tame-impala'), (SELECT id FROM genres WHERE slug = 'indie')),
((SELECT id FROM bands WHERE slug = 'tame-impala'), (SELECT id FROM genres WHERE slug = 'electronic'));

-- String Cheese Incident: Jam Band, Bluegrass, Electronic
INSERT INTO band_genres (band_id, genre_id) VALUES
((SELECT id FROM bands WHERE slug = 'string-cheese-incident'), (SELECT id FROM genres WHERE slug = 'jam-band')),
((SELECT id FROM bands WHERE slug = 'string-cheese-incident'), (SELECT id FROM genres WHERE slug = 'bluegrass')),
((SELECT id FROM bands WHERE slug = 'string-cheese-incident'), (SELECT id FROM genres WHERE slug = 'electronic'));

-- Umphrey's McGee: Jam Band, Rock
INSERT INTO band_genres (band_id, genre_id) VALUES
((SELECT id FROM bands WHERE slug = 'umphreys-mcgee'), (SELECT id FROM genres WHERE slug = 'jam-band')),
((SELECT id FROM bands WHERE slug = 'umphreys-mcgee'), (SELECT id FROM genres WHERE slug = 'rock'));

-- Pigeons Playing Ping Pong: Jam Band, Funk
INSERT INTO band_genres (band_id, genre_id) VALUES
((SELECT id FROM bands WHERE slug = 'pigeons-playing-ping-pong'), (SELECT id FROM genres WHERE slug = 'jam-band')),
((SELECT id FROM bands WHERE slug = 'pigeons-playing-ping-pong'), (SELECT id FROM genres WHERE slug = 'funk'));

-- ODESZA: Electronic
INSERT INTO band_genres (band_id, genre_id) VALUES
((SELECT id FROM bands WHERE slug = 'odesza'), (SELECT id FROM genres WHERE slug = 'electronic'));

-- Big Gigantic: Electronic, Funk
INSERT INTO band_genres (band_id, genre_id) VALUES
((SELECT id FROM bands WHERE slug = 'big-gigantic'), (SELECT id FROM genres WHERE slug = 'electronic')),
((SELECT id FROM bands WHERE slug = 'big-gigantic'), (SELECT id FROM genres WHERE slug = 'funk'));

-- Anderson .Paak: Hip Hop, R&B, Soul
INSERT INTO band_genres (band_id, genre_id) VALUES
((SELECT id FROM bands WHERE slug = 'anderson-paak'), (SELECT id FROM genres WHERE slug = 'hip-hop')),
((SELECT id FROM bands WHERE slug = 'anderson-paak'), (SELECT id FROM genres WHERE slug = 'soul'));

-- GRiZ: Electronic, Funk
INSERT INTO band_genres (band_id, genre_id) VALUES
((SELECT id FROM bands WHERE slug = 'griz'), (SELECT id FROM genres WHERE slug = 'electronic')),
((SELECT id FROM bands WHERE slug = 'griz'), (SELECT id FROM genres WHERE slug = 'funk'));

-- Billy Strings: Bluegrass
INSERT INTO band_genres (band_id, genre_id) VALUES
((SELECT id FROM bands WHERE slug = 'billy-strings'), (SELECT id FROM genres WHERE slug = 'bluegrass'));

-- Mandolin Orange: Folk, Bluegrass
INSERT INTO band_genres (band_id, genre_id) VALUES
((SELECT id FROM bands WHERE slug = 'mandolin-orange'), (SELECT id FROM genres WHERE slug = 'folk')),
((SELECT id FROM bands WHERE slug = 'mandolin-orange'), (SELECT id FROM genres WHERE slug = 'bluegrass'));

-- Whitechapel: Metal
INSERT INTO band_genres (band_id, genre_id) VALUES
((SELECT id FROM bands WHERE slug = 'whitechapel'), (SELECT id FROM genres WHERE slug = 'metal'));

-- Bodysnatcher: Metal
INSERT INTO band_genres (band_id, genre_id) VALUES
((SELECT id FROM bands WHERE slug = 'bodysnatcher'), (SELECT id FROM genres WHERE slug = 'metal'));

-- The Get Right Band: Funk, Soul
INSERT INTO band_genres (band_id, genre_id) VALUES
((SELECT id FROM bands WHERE slug = 'the-get-right-band'), (SELECT id FROM genres WHERE slug = 'funk')),
((SELECT id FROM bands WHERE slug = 'the-get-right-band'), (SELECT id FROM genres WHERE slug = 'soul'));

-- Town Mountain: Bluegrass, Folk
INSERT INTO band_genres (band_id, genre_id) VALUES
((SELECT id FROM bands WHERE slug = 'town-mountain'), (SELECT id FROM genres WHERE slug = 'bluegrass')),
((SELECT id FROM bands WHERE slug = 'town-mountain'), (SELECT id FROM genres WHERE slug = 'folk'));

-- ==============================================================================
-- SAMPLE SHOWS (Past, Present, Future)
-- ==============================================================================

-- Get venue IDs (assumes venues from 002_venues.sql are loaded)
DO $$
DECLARE
    orange_peel_id INTEGER;
    grey_eagle_id INTEGER;
    salvage_station_id INTEGER;

    moon_taxi_id INTEGER;
    revivalists_id INTEGER;
    tame_impala_id INTEGER;
    sci_id INTEGER;
    umphreys_id INTEGER;
    pigeons_id INTEGER;
    odesza_id INTEGER;
    big_gig_id INTEGER;
    paak_id INTEGER;
    griz_id INTEGER;
    billy_id INTEGER;
    mandolin_id INTEGER;
    whitechapel_id INTEGER;
    bodysnatcher_id INTEGER;
    get_right_id INTEGER;
    town_mountain_id INTEGER;
BEGIN
    -- Get venue IDs
    SELECT id INTO orange_peel_id FROM venues WHERE slug = 'the-orange-peel';
    SELECT id INTO grey_eagle_id FROM venues WHERE slug = 'the-grey-eagle';
    SELECT id INTO salvage_station_id FROM venues WHERE slug = 'salvage-station';

    -- Get band IDs
    SELECT id INTO moon_taxi_id FROM bands WHERE slug = 'moon-taxi';
    SELECT id INTO revivalists_id FROM bands WHERE slug = 'the-revivalists';
    SELECT id INTO tame_impala_id FROM bands WHERE slug = 'tame-impala';
    SELECT id INTO sci_id FROM bands WHERE slug = 'string-cheese-incident';
    SELECT id INTO umphreys_id FROM bands WHERE slug = 'umphreys-mcgee';
    SELECT id INTO pigeons_id FROM bands WHERE slug = 'pigeons-playing-ping-pong';
    SELECT id INTO odesza_id FROM bands WHERE slug = 'odesza';
    SELECT id INTO big_gig_id FROM bands WHERE slug = 'big-gigantic';
    SELECT id INTO paak_id FROM bands WHERE slug = 'anderson-paak';
    SELECT id INTO griz_id FROM bands WHERE slug = 'griz';
    SELECT id INTO billy_id FROM bands WHERE slug = 'billy-strings';
    SELECT id INTO mandolin_id FROM bands WHERE slug = 'mandolin-orange';
    SELECT id INTO whitechapel_id FROM bands WHERE slug = 'whitechapel';
    SELECT id INTO bodysnatcher_id FROM bands WHERE slug = 'bodysnatcher';
    SELECT id INTO get_right_id FROM bands WHERE slug = 'the-get-right-band';
    SELECT id INTO town_mountain_id FROM bands WHERE slug = 'town-mountain';

    -- PAST SHOWS (for testing historical data)
    INSERT INTO shows (venue_id, title, date, doors_time, show_time, price_min, price_max, ticket_url, age_restriction, status, source) VALUES
    (orange_peel_id, 'Moon Taxi', NOW() - INTERVAL '30 days', '19:00', '20:00', 25.00, 35.00, 'https://example.com/tickets/1', 'All Ages', 'completed', 'scraped'),
    (grey_eagle_id, 'Billy Strings', NOW() - INTERVAL '15 days', '18:30', '19:30', 45.00, 55.00, 'https://example.com/tickets/2', 'All Ages', 'completed', 'scraped');

    -- Link bands to past shows
    INSERT INTO show_bands (show_id, band_id, is_headliner, performance_order) VALUES
    ((SELECT id FROM shows WHERE title = 'Moon Taxi' AND status = 'completed'), moon_taxi_id, true, 1);

    INSERT INTO show_bands (show_id, band_id, is_headliner, performance_order) VALUES
    ((SELECT id FROM shows WHERE title = 'Billy Strings' AND status = 'completed'), billy_id, true, 1);

    -- UPCOMING SHOWS (various dates)

    -- Show 1: Tonight (free show)
    INSERT INTO shows (venue_id, title, date, doors_time, show_time, price_min, price_max, ticket_url, age_restriction, status, source) VALUES
    (grey_eagle_id, 'The Get Right Band', NOW() + INTERVAL '8 hours', '19:00', '20:00', NULL, NULL, NULL, '21+', 'scheduled', 'manual')
    RETURNING id INTO @show_id;

    INSERT INTO show_bands (show_id, band_id, is_headliner, performance_order) VALUES
    (LASTVAL(), get_right_id, true, 1);

    -- Show 2: This weekend - Multi-band lineup
    INSERT INTO shows (venue_id, title, image_url, date, doors_time, show_time, price_min, price_max, ticket_url, age_restriction, status, source) VALUES
    (orange_peel_id, 'Whitechapel with Bodysnatcher', 'https://example.com/posters/whitechapel.jpg', NOW() + INTERVAL '3 days', '19:00', '20:00', 30.00, 40.00, 'https://example.com/tickets/3', 'All Ages', 'scheduled', 'scraped');

    INSERT INTO show_bands (show_id, band_id, is_headliner, performance_order) VALUES
    (LASTVAL(), whitechapel_id, true, 2),
    (LASTVAL(), bodysnatcher_id, false, 1);

    -- Show 3: Next week - Jam band
    INSERT INTO shows (venue_id, title, date, doors_time, show_time, price_min, price_max, ticket_url, age_restriction, status, source) VALUES
    (salvage_station_id, 'Pigeons Playing Ping Pong', NOW() + INTERVAL '7 days', '19:30', '20:30', 20.00, 25.00, 'https://example.com/tickets/4', 'All Ages', 'scheduled', 'scraped');

    INSERT INTO show_bands (show_id, band_id, is_headliner, performance_order) VALUES
    (LASTVAL(), pigeons_id, true, 1);

    -- Show 4: Two weeks out - Big electronic show
    INSERT INTO shows (venue_id, title, image_url, date, doors_time, show_time, price_min, price_max, ticket_url, age_restriction, status, source) VALUES
    (orange_peel_id, 'ODESZA', 'https://example.com/posters/odesza.jpg', NOW() + INTERVAL '14 days', '20:00', '21:00', 55.00, 75.00, 'https://example.com/tickets/5', 'All Ages', 'scheduled', 'scraped');

    INSERT INTO show_bands (show_id, band_id, is_headliner, performance_order) VALUES
    (LASTVAL(), odesza_id, true, 1);

    -- Show 5: Next month - Psychedelic rock
    INSERT INTO shows (venue_id, title, image_url, date, doors_time, show_time, price_min, price_max, ticket_url, age_restriction, status, source) VALUES
    (orange_peel_id, 'Tame Impala', 'https://example.com/posters/tame-impala.jpg', NOW() + INTERVAL '30 days', '19:00', '20:00', 65.00, 85.00, 'https://example.com/tickets/6', 'All Ages', 'scheduled', 'scraped');

    INSERT INTO show_bands (show_id, band_id, is_headliner, performance_order) VALUES
    (LASTVAL(), tame_impala_id, true, 1);

    -- Show 6: Next month - Multi-night residency (Night 1)
    INSERT INTO shows (venue_id, title, date, doors_time, show_time, price_min, price_max, ticket_url, age_restriction, status, source) VALUES
    (orange_peel_id, 'String Cheese Incident - Night 1', NOW() + INTERVAL '45 days', '19:00', '20:00', 50.00, 60.00, 'https://example.com/tickets/7', 'All Ages', 'scheduled', 'scraped');

    INSERT INTO show_bands (show_id, band_id, is_headliner, performance_order) VALUES
    (LASTVAL(), sci_id, true, 1);

    -- Show 7: Next month - Multi-night residency (Night 2)
    INSERT INTO shows (venue_id, title, date, doors_time, show_time, price_min, price_max, ticket_url, age_restriction, status, source) VALUES
    (orange_peel_id, 'String Cheese Incident - Night 2', NOW() + INTERVAL '46 days', '19:00', '20:00', 50.00, 60.00, 'https://example.com/tickets/8', 'All Ages', 'scheduled', 'scraped');

    INSERT INTO show_bands (show_id, band_id, is_headliner, performance_order) VALUES
    (LASTVAL(), sci_id, true, 1);

    -- Show 8: Far future - Big hip hop show
    INSERT INTO shows (venue_id, title, date, doors_time, show_time, price_min, price_max, ticket_url, age_restriction, status, source) VALUES
    (orange_peel_id, 'Anderson .Paak & The Free Nationals', NOW() + INTERVAL '60 days', '20:00', '21:00', 70.00, 90.00, 'https://example.com/tickets/9', 'All Ages', 'scheduled', 'scraped');

    INSERT INTO show_bands (show_id, band_id, is_headliner, performance_order) VALUES
    (LASTVAL(), paak_id, true, 1);

    -- Show 9: Local show - Bluegrass
    INSERT INTO shows (venue_id, title, date, doors_time, show_time, price_min, price_max, ticket_url, age_restriction, status, source) VALUES
    (grey_eagle_id, 'Town Mountain with Mandolin Orange', NOW() + INTERVAL '20 days', '18:30', '19:30', 18.00, 22.00, 'https://example.com/tickets/10', 'All Ages', 'scheduled', 'manual');

    INSERT INTO show_bands (show_id, band_id, is_headliner, performance_order) VALUES
    (LASTVAL(), town_mountain_id, true, 2),
    (LASTVAL(), mandolin_id, false, 1);

    -- Show 10: Cancelled show (for testing status filter)
    INSERT INTO shows (venue_id, title, date, doors_time, show_time, price_min, price_max, ticket_url, age_restriction, status, source) VALUES
    (salvage_station_id, 'GRiZ', NOW() + INTERVAL '25 days', '19:00', '20:00', 35.00, 45.00, NULL, 'All Ages', 'cancelled', 'scraped');

    INSERT INTO show_bands (show_id, band_id, is_headliner, performance_order) VALUES
    (LASTVAL(), griz_id, true, 1);

    -- Show 11: Band submitted show (pending approval)
    INSERT INTO shows (venue_id, title, date, doors_time, show_time, price_min, price_max, ticket_url, age_restriction, status, source) VALUES
    (grey_eagle_id, 'The Revivalists', NOW() + INTERVAL '50 days', '19:00', '20:00', 30.00, 40.00, 'https://example.com/tickets/11', 'All Ages', 'pending', 'band_submitted');

    INSERT INTO show_bands (show_id, band_id, is_headliner, performance_order) VALUES
    (LASTVAL(), revivalists_id, true, 1);

END $$;

-- ==============================================================================
-- Summary
-- ==============================================================================

-- Verify data loaded
DO $$
DECLARE
    band_count INTEGER;
    show_count INTEGER;
    upcoming_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO band_count FROM bands;
    SELECT COUNT(*) INTO show_count FROM shows;
    SELECT COUNT(*) INTO upcoming_count FROM shows WHERE status = 'scheduled' AND date >= NOW();

    RAISE NOTICE 'Test data loaded successfully:';
    RAISE NOTICE '  Bands: %', band_count;
    RAISE NOTICE '  Shows (total): %', show_count;
    RAISE NOTICE '  Shows (upcoming): %', upcoming_count;
END $$;
