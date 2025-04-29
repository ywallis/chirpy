-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, password)
VALUES (
	gen_random_uuid (),
	NOW(),
	NOW(),
	$1,
	$2
)
RETURNING *;

-- name: DeleteAllUsers :exec

DELETE FROM users;

-- name: GetUserByEmail :one

SELECT * FROM users
WHERE $1 = email;


-- name: UpdateCredentials :one
UPDATE users
SET email = $1, password = $2, updated_at = $3
WHERE id = $4
RETURNING *;
