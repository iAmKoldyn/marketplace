-- name: CreateAd :one
INSERT INTO ads (author_id, title, text, image_url, price)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, author_id, title, text, image_url, price, created_at;

-- name: ListAdsByDateDesc :many
SELECT a.id, a.author_id, u.username AS author_username,
       a.title, a.text, a.image_url, a.price, a.created_at
FROM ads a
JOIN users u ON u.id = a.author_id
WHERE a.price BETWEEN $1 AND $2
ORDER BY a.created_at DESC
LIMIT $3 OFFSET $4;

-- name: ListAdsByDateAsc :many
-- same with ORDER BY a.created_at ASC

-- name: ListAdsByPriceDesc :many
-- same with ORDER BY a.price DESC

-- name: ListAdsByPriceAsc :many
-- same with ORDER BY a.price ASC
