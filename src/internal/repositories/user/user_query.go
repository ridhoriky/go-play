package user

const (
	createUserQuery = `
	INSERT INTO users (id, name, email, password_hash, role, is_active, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	getAllUsersQuery = `
	SELECT u.id, u.name, u.email, u.role, u.is_active, u.created_at, u.updated_at,
		COUNT(*) OVER() AS total_count
	FROM users u
	WHERE 1=1
	`

	getUserByIDQuery = `
	SELECT id, name, email, password_hash, role, is_active, created_at, updated_at
	FROM users
	WHERE id = $1
	`

	getUserByEmailQuery = `
	SELECT id, name, email, password_hash, role, is_active, created_at, updated_at
	FROM users
	WHERE email = $1
	`

	updateUserQuery = `
	UPDATE users
	SET name = $1, email = $2, role = $3, updated_at = NOW()
	WHERE id = $4
	`

	deleteUserQuery = `
	UPDATE users
	SET is_active = FALSE, updated_at = NOW()
	WHERE id = $1
	`
)
