-- name: SaveUser :one
INSERT INTO _user.users
    (
        first_name,
        last_name,
        phone,
        email,
        password_hash,
        google_user_id,
        oauth_refresh_token_encrypted,
        last_ip_address
    )
VALUES
    (
        $1, $2, $3, $4, $5, $6, $7, $8
    )
RETURNING 
    id;

-- name: ToggleTermsAndConditionsAccepted :exec
UPDATE _user.users
SET terms_and_condition_accepted = true
WHERE id = $1;

-- name: CreateUser :one
INSERT INTO _user.users
    (
        first_name,
        last_name,
        phone,
        email,
        password_hash,
        google_user_id,
        oauth_refresh_token_encrypted,
        last_ip_address
    )
VALUES
    (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8
    )
RETURNING
    id;

-- name: GetUser :one
SELECT
    id,
    first_name,
    last_name,
    phone,
    email,
    password_hash,
    password_last_changed,
    google_user_id,
    oauth_refresh_token_encrypted,
    last_ip_address
FROM _user.users
WHERE id = $1;


-- name: CreateEmailVerificationCode :one
INSERT INTO _user.email_verification_codes
(
    email,
    verification_codes_id,
    created_at
)
VALUES
(
    $1,
    $2,
    now()
)
ON CONFLICT (email)
DO UPDATE SET 
    verification_codes_id = $2,
    created_at = now()
RETURNING
    id;

-- name: GetEmailVerificationCodeByEmail :one
SELECT 
    id,
    verification_codes_id,
    created_at
FROM _user.email_verification_codes
WHERE email = $1;


INSERT INTO users (email, name, last_login) 
VALUES ('user@example.com', 'Jane Doe', NOW())
ON CONFLICT (email) 
DO UPDATE SET 
    name = EXCLUDED.name, 
    last_login = EXCLUDED.last_login;

-- name: GetUserByEmail :one
SELECT
    id, 
    first_name,
    last_name,
    phone,
    email,
    password_hash,
    password_last_changed,
    google_user_id,
    oauth_refresh_token_encrypted,
    last_ip_address
FROM _user.users
WHERE email = sqlc.arg(email)::text;

-- name: UpdateUserPassword :exec
UPDATE _user.users
SET password_hash = $1,
password_last_changed = CURRENT_TIMESTAMP
WHERE id = $2;

-- name: UpdateUserLastIpAddress :exec
UPDATE _user.users
SET last_ip_address = $1
WHERE id = $2;

-- name: UpdateUserOAUthCredentials :exec
UPDATE _user.users
SET google_user_id = $1,
oauth_refresh_token_encrypted = $2
WHERE id = $2;

-- name: GetSystemAdministratorByUserId :one
SELECT
    id,
    created_at
FROM _user.system_administrators
WHERE user_id = $1;

-- name: GetSystemAdministrators :many
SELECT 
    id,
    user_id,
    created_at
FROM _user.system_administrators;

-- name: CreateVerificationSession :one
INSERT INTO _user.verification_sessions
    (
        user_id,
        verification_codes_id
    )
VALUES
    (
        $1,
        $2
    )
RETURNING
    ID;

-- name: GetVerificationSessionByUserId :one
SELECT 
    id,
    verification_codes_id,
    created_at
FROM _user.verification_sessions
WHERE user_id = $1
ORDER BY id DESC
LIMIT 1;

-- name: DeleteVerificationSessionByEmail :exec
DELETE FROM _user.verification_sessions
WHERE user_id IN (SELECT id FROM _user.users WHERE email = $1);

-- name: GetVerificationSessionBySessionToken :one
SELECT
    _user.verification_sessions.user_id,
    shared.verification_codes.verification_code,
    shared.verification_codes.failed_attempts,
    shared.verification_codes.created_at
FROM _user.verification_sessions
JOIN shared.verification_codes
ON _user.verification_sessions.verification_codes_id = shared.verification_codes.id
WHERE shared.verification_codes.id = sqlc.arg(session_token);

-- name: DeleteUserByEmail :exec
DELETE FROM _user.users WHERE email = $1;

-- name: DeleteEmailVerificationCodeByEmail :exec
DELETE FROM _user.email_verification_codes WHERE email = $1;