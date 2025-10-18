CREATE TABLE downloads(
  id VARCHAR(255) PRIMARY KEY,
  options TEXT NOT NULL,
  url VARCHAR(255) NOT NULL,
  finished_at DATETIME
);