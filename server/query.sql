-- name: CreateMailbox :one
INSERT INTO
    mailserver.virtual_users (email, domain_id, password, temporary)
VALUES
    ($1, $2, $3, TRUE)
RETURNING *;

-- name: DeleteMailbox :exec
DELETE FROM mailserver.virtual_users WHERE id = $1;