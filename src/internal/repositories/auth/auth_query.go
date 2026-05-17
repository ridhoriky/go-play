package auth

//nolint:gosec // SQL queries are not credentials
const (
	createRefreshTokenQuery = `
		INSERT INTO refresh_tokens (token_hash, user_id, user_agent, ip_address, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	getRefreshTokenByHashQuery = `
	SELECT id, user_id, token_hash, COALESCE(user_agent, ''), COALESCE(ip_address::text, ''), expires_at, created_at, revoked_at
	FROM refresh_tokens
	WHERE token_hash = $1
	`

	revokeRefreshTokenQuery = `
		UPDATE refresh_tokens
		SET revoked_at = $1
		WHERE token_hash = $2 AND revoked_at IS NULL
	`
)
