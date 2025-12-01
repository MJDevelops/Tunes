PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS downloads(
  id VARCHAR(255) PRIMARY KEY NOT NULL,
  options TEXT NOT NULL,
  url TEXT NOT NULL,
  finished_at DATETIME
);

CREATE TABLE IF NOT EXISTS artists(
  id INTEGER PRIMARY KEY NOT NULL,
  name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS tracks(
  id INTEGER PRIMARY KEY NOT NULL,
  path TEXT NOT NULL,
  album_id INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS albums(
  id INTEGER PRIMARY KEY NOT NULL,
  title TEXT NOT NULL,
  artwork BLOB
);

CREATE TABLE IF NOT EXISTS playlists(
  id INTEGER PRIMARY KEY NOT NULL,
  title VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS playlists_tracks(
  track_id INTEGER NOT NULL,
  playlist_id INTEGER NOT NULL,

  PRIMARY KEY (track_id,playlist_id),
  FOREIGN KEY (track_id) REFERENCES tracks(id),
  FOREIGN KEY (playlist_id) REFERENCES playlists(id)
);

CREATE TABLE IF NOT EXISTS artists_albums(
  artist_id INTEGER NOT NULL,
  album_id INTEGER NOT NULL,

  PRIMARY KEY (artist_id,album_id),
  FOREIGN KEY (artist_id) REFERENCES artists(id),
  FOREIGN KEY (album_id) REFERENCES albums(id)
);

CREATE TABLE IF NOT EXISTS artists_tracks(
  artist_id INTEGER NOT NULL,
  track_id INTEGER NOT NULL,

  PRIMARY KEY (artist_id,track_id),
  FOREIGN KEY (artist_id) REFERENCES artists(id),
  FOREIGN KEY (track_id) REFERENCES tracks(id)
);