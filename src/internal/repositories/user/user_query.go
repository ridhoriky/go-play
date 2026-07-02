package user

const (
	createUserQuery = `
	INSERT INTO users (id, name, email, password_hash, role, is_active, is_verified, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	getAllUsersQuery = `
	SELECT u.id, u.name, u.email, u.role, u.is_active, u.is_verified, u.created_at, u.updated_at,
		COUNT(*) OVER() AS total_count
	FROM users u
	WHERE 1=1
	`

	getUserByIDQuery = `
	SELECT id, name, email, password_hash, role, is_active, is_verified, created_at, updated_at
	FROM users
	WHERE id = $1
	`

	getUserByEmailQuery = `
	SELECT id, name, email, password_hash, role, is_active, is_verified, created_at, updated_at
	FROM users
	WHERE email = $1
	`

	updateUserQuery = `
	UPDATE users
	SET name = $1, email = $2, role = $3, is_verified = $4, is_active = $5, updated_at = NOW()
	WHERE id = $6
	`

	deleteUserQuery = `
	UPDATE users
	SET is_active = FALSE, updated_at = NOW()
	WHERE id = $1
	`
)
