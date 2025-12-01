PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS downloads(
  id VARCHAR(255) PRIMARY KEY,
  options TEXT NOT NULL,
  url TEXT NOT NULL,
  finished_at DATETIME
);

CREATE TABLE IF NOT EXISTS tracks(
  id INTEGER PRIMARY KEY,
  path TEXT NOT NULL,
  album_id INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS albums(
  id INTEGER PRIMARY KEY,
  title TEXT NOT NULL,
  artwork BLOB
);

CREATE TABLE IF NOT EXISTS playlists(
  id INTEGER PRIMARY KEY,
  title TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS playlists_tracks(
  track_id INTEGER NOT NULL,
  playlist_id INTEGER NOT NULL,

  PRIMARY KEY (track_id,playlist_id),
  FOREIGN KEY (track_id) REFERENCES tracks(id),
  FOREIGN KEY (playlist_id) REFERENCES playlists(id)
);