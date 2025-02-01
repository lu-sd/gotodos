-- name: GetTodos :many
SELECT id, title, done FROM todos2;

-- name: CreateTodo :execresult
INSERT INTO todos2 (title, done) VALUES (?, ?) RETURNING id;

-- name: GetTodoByID :one
SELECT id, title, done FROM todos2 WHERE id = ?;

-- name: UpdateTodoStatus :exec
UPDATE todos2 SET done = ? WHERE id = ?;

-- name: DeleteTodo :exec
DELETE FROM todos2 WHERE id = ?;
