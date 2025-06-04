-- Muse Database Schema
-- Development initialization script with sample data

-- Import the main schema
\i /docker-entrypoint-initdb.d/init.sql

-- Insert sample data for development

-- Sample artists
INSERT INTO artists (name, bio, image_url, spotify_id) VALUES
    ('Taylor Swift', 'American singer-songwriter known for narrative songs about her personal life', 'https://example.com/taylor-swift.jpg', '06HL4z0CvFAxyc27GXpf02'),
    ('The Beatles', 'English rock band formed in Liverpool in 1960', 'https://example.com/beatles.jpg', '3WrFJ7ztbogyGnTHbHJFl2'),
    ('Drake', 'Canadian rapper, singer, and songwriter', 'https://example.com/drake.jpg', '3TVXtAsR1Inumwj472S9r4'),
    ('Billie Eilish', 'American singer and songwriter', 'https://example.com/billie-eilish.jpg', '6qqNVTkY8uBg9cP3Jd8DAH');

-- Sample albums
INSERT INTO albums (title, artist_id, release_date, genre, cover_image_url, spotify_id) VALUES
    ('1989 (Taylor''s Version)', (SELECT artist_id FROM artists WHERE name = 'Taylor Swift'), '2023-10-27', 'Pop', 'https://example.com/1989-tv.jpg', '1o59UpKw81iHR0HPiSkJR0'),
    ('Abbey Road', (SELECT artist_id FROM artists WHERE name = 'The Beatles'), '1969-09-26', 'Rock', 'https://example.com/abbey-road.jpg', '0ETFjACtuP2ADo6LFhL6HN'),
    ('Scorpion', (SELECT artist_id FROM artists WHERE name = 'Drake'), '2018-06-29', 'Hip Hop', 'https://example.com/scorpion.jpg', '1ATL5GLyefJaxhQzSPVrLX'),
    ('Happier Than Ever', (SELECT artist_id FROM artists WHERE name = 'Billie Eilish'), '2021-07-30', 'Alternative', 'https://example.com/happier-than-ever.jpg', '0JGOiO34nkfUdDruD7qmgf');

-- Sample songs
INSERT INTO songs (title, artist_id, album_id, duration_seconds, track_number, spotify_id) VALUES
    ('Shake It Off', (SELECT artist_id FROM artists WHERE name = 'Taylor Swift'), (SELECT album_id FROM albums WHERE title = '1989 (Taylor''s Version)'), 219, 6, '5AsrX0Shds3V1Mm8s8Qw1D'),
    ('Come Together', (SELECT artist_id FROM artists WHERE name = 'The Beatles'), (SELECT album_id FROM albums WHERE title = 'Abbey Road'), 259, 1, '2EqlS6tkEnglzr7tkKAAYD'),
    ('God''s Plan', (SELECT artist_id FROM artists WHERE name = 'Drake'), (SELECT album_id FROM albums WHERE title = 'Scorpion'), 198, 2, '6DCZcSspjsKoFjzjrWoCdn'),
    ('bad guy', (SELECT artist_id FROM artists WHERE name = 'Billie Eilish'), (SELECT album_id FROM albums WHERE title = 'Happier Than Ever'), 194, 2, '2Fxmhks0bxGSBdJ92vM42m');

-- Sample users
INSERT INTO users (email, username, first_name, last_name, password_hash, is_verified) VALUES
    ('john.doe@example.com', 'johndoe', 'John', 'Doe', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBPj9.3X7O5O8m', true),
    ('jane.smith@example.com', 'janesmith', 'Jane', 'Smith', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBPj9.3X7O5O8m', true),
    ('music.lover@example.com', 'musiclover', 'Music', 'Lover', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBPj9.3X7O5O8m', false);

-- Sample playlists
INSERT INTO playlists (user_id, title, description, is_public) VALUES
    ((SELECT user_id FROM users WHERE username = 'johndoe'), 'My Favorites', 'My all-time favorite songs', true),
    ((SELECT user_id FROM users WHERE username = 'janesmith'), 'Workout Playlist', 'High energy songs for working out', true),
    ((SELECT user_id FROM users WHERE username = 'musiclover'), 'Chill Vibes', 'Relaxing songs for studying', false);

-- Sample playlist songs
INSERT INTO playlist_songs (playlist_id, song_id, position) VALUES
    ((SELECT playlist_id FROM playlists WHERE title = 'My Favorites'), (SELECT song_id FROM songs WHERE title = 'Shake It Off'), 1),
    ((SELECT playlist_id FROM playlists WHERE title = 'My Favorites'), (SELECT song_id FROM songs WHERE title = 'bad guy'), 2),
    ((SELECT playlist_id FROM playlists WHERE title = 'Workout Playlist'), (SELECT song_id FROM songs WHERE title = 'God''s Plan'), 1),
    ((SELECT playlist_id FROM playlists WHERE title = 'Chill Vibes'), (SELECT song_id FROM songs WHERE title = 'Come Together'), 1);

-- Sample ratings
INSERT INTO ratings (user_id, song_id, rating, review) VALUES
    ((SELECT user_id FROM users WHERE username = 'johndoe'), (SELECT song_id FROM songs WHERE title = 'Shake It Off'), 4.5, 'Classic Taylor Swift! Love the energy.'),
    ((SELECT user_id FROM users WHERE username = 'janesmith'), (SELECT song_id FROM songs WHERE title = 'bad guy'), 5.0, 'Billie''s unique style is amazing!'),
    ((SELECT user_id FROM users WHERE username = 'musiclover'), (SELECT song_id FROM songs WHERE title = 'Come Together'), 4.8, 'Timeless Beatles classic.');

-- Sample album ratings
INSERT INTO ratings (user_id, album_id, rating, review) VALUES
    ((SELECT user_id FROM users WHERE username = 'johndoe'), (SELECT album_id FROM albums WHERE title = '1989 (Taylor''s Version)'), 4.7, 'Great re-recording of a pop masterpiece'),
    ((SELECT user_id FROM users WHERE username = 'janesmith'), (SELECT album_id FROM albums WHERE title = 'Abbey Road'), 5.0, 'One of the greatest albums ever made');

-- Sample user follows
INSERT INTO user_follows (follower_id, following_id) VALUES
    ((SELECT user_id FROM users WHERE username = 'johndoe'), (SELECT user_id FROM users WHERE username = 'janesmith'')),
    ((SELECT user_id FROM users WHERE username = 'janesmith'), (SELECT user_id FROM users WHERE username = 'musiclover'));

-- Sample recommendations
INSERT INTO recommendations (user_id, song_id, score, reason) VALUES
    ((SELECT user_id FROM users WHERE username = 'johndoe'), (SELECT song_id FROM songs WHERE title = 'God''s Plan'), 0.8532, 'Based on your high rating of pop songs'),
    ((SELECT user_id FROM users WHERE username = 'janesmith'), (SELECT song_id FROM songs WHERE title = 'Come Together'), 0.7821, 'Users with similar taste also liked this');

COMMIT; 