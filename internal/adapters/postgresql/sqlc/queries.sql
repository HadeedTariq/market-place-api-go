-- name: FindExistingUserByEmail :one
select is_verified from users where email = $1;

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

-- name: InsertEmailOtp :exec
INSERT INTO email_otps (email, otp, expires_at) VALUES ($1, $2, $3);

-- name: FindExistingOtp :one
SELECT 1 FROM email_otps WHERE email = $1 AND expires_at > NOW();

-- name: CheckOtp :one
SELECT 1 
FROM email_otps 
WHERE email = $1 
  AND otp = $2
  AND expires_at > NOW()
ORDER BY created_at DESC 
LIMIT 1;