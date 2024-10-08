-- name: GetGame :one
SELECT * FROM games
WHERE id = ? LIMIT 1;

-- name: ListGames :many
SELECT * FROM games
ORDER BY game;

-- name: CreateGame :one
INSERT INTO games (
  game, year , dev,  publisher  , platform
) VALUES (
  ?, ? , ? , ? , ?
)
RETURNING *;
