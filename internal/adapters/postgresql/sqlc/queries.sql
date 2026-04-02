-- name: FindExistingUserByEmail :one
select 1 from users where email = $1;

-- name: InsertUser :exec
INSERT INTO users (
    user_name, 
    email, 
    password_hash, 
    role, 
    source, 
    country_code, 
    gender
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);