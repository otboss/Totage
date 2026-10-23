-- name: GetNull :many
SELECT *;

-- name: CreateVerificationSession :one
INSERT INTO shared.verification_codes DEFAULT VALUES
RETURNING 
    id,
    verification_code;

-- name: GetVerificationCode :one
SELECT 
    verification_code,
    failed_attempts,
    is_completed,
    created_at
FROM shared.verification_codes
WHERE id = $1;

-- -- name: IsVerificationSessionSuccessful :one
-- SELECT EXISTS (
--     SELECT 1
--     FROM shared.verification_codes
--     WHERE verification_code = $1
--     AND id = sqlc.arg(session_token)
--     AND created_at >= NOW() - INTERVAL '5 minutes'
--     AND failed_attempts < 3
-- );

-- -- name: IncrementVerificationSessionFailedAttempts :exec
-- UPDATE shared.verification_codes
-- SET failed_attempts = failed_attempts + 1
-- WHERE id = sqlc.arg(session_token);

-- -- name: DeleteVerificationSession :exec
-- DELETE FROM shared.verification_codes
-- WHERE id = sqlc.arg(session_token);


-- name: IsVerificationSessionSuccessful :one
WITH eligible AS (
    SELECT id, verification_code
    FROM shared.verification_codes
    WHERE shared.verification_codes.id = $1
      AND shared.verification_codes.is_completed = FALSE
      AND shared.verification_codes.created_at >= now() - interval '5 minutes'
      AND shared.verification_codes.failed_attempts < 3
    FOR UPDATE
),
success AS (
    UPDATE shared.verification_codes
    SET is_completed = TRUE
    WHERE id = (SELECT id FROM eligible WHERE verification_codes.verification_code = $2)
    RETURNING true AS result
),
failure AS (
    UPDATE shared.verification_codes
    SET failed_attempts = failed_attempts + 1
    WHERE id = (SELECT id FROM eligible WHERE verification_codes.verification_code <> $2)
    RETURNING false AS result
)
SELECT COALESCE(
    (SELECT result FROM success),
    (SELECT result FROM failure),
    false
)::boolean AS verified;

-- WITH eligible AS (
--     SELECT id, verification_code
--     FROM shared.verification_codes
--     WHERE shared.verification_codes.id = $1
--       AND shared.verification_codes.is_completed = FALSE
--       AND shared.verification_codes.created_at >= now() - interval '5 minutes'
--       AND shared.verification_codes.failed_attempts < 3
--     FOR UPDATE
-- ),
-- success AS (
--     UPDATE shared.verification_codes
--     SET is_completed = TRUE
--     WHERE id = (SELECT id FROM eligible WHERE verification_code = $2)
--     RETURNING true AS result
-- ),
-- failure AS (
--     UPDATE shared.verification_codes
--     SET shared.verification_codes.failed_attempts = shared.verification_codes.failed_attempts + 1
--     WHERE id = (SELECT id FROM eligible WHERE verification_code <> $2)
--     RETURNING false AS result
-- )
-- SELECT COALESCE(
--     (SELECT result FROM success),
--     (SELECT result FROM failure),
--     false
-- )::boolean AS verified;