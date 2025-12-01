-- name: GetPendingDownloads :many
SELECT * FROM downloads WHERE finished_at IS NULL;

-- name: GetDownload :one
SELECT * FROM downloads WHERE id = ?;

-- name: InsertDownload :exec
INSERT INTO downloads (
  id, url, options, finished_at
  ) VALUES (
  ?, ?, ?, ?
);

-- name: UpdateDownloadFinishedAt :exec
UPDATE downloads SET finished_at = ? WHERE id = ?;

-- name: GetTrack :one
SELECT * FROM tracks WHERE id = ?;

-- name: GetPlaylistTracks :many
SELECT t.* FROM tracks t
JOIN playlists_tracks pt ON t.id = pt.track_id
WHERE pt.playlist_id = ?;

-- name: GetPlaylist :one
SELECT * FROM playlists WHERE id = ?;

-- name: GetArtist :one
SELECT * FROM artists WHERE id = ?;

-- name: GetAlbum :one
SELECT * FROM albums WHERE id = ?;

-- name: GetAlbumByTitle :many
SELECT * FROM albums WHERE title = ?;