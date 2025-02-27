CREATE TABLE IF NOT EXISTS Songs (
    song_id SERIAL PRIMARY KEY,
    group_name VARCHAR(64) NOT NULL,
    song_name VARCHAR(64) NOT NULL,
    release_date VARCHAR(10),
    lyrics TEXT,
    link VARCHAR(128)
);