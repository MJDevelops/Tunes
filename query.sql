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