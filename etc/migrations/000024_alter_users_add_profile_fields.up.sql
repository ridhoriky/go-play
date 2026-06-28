ALTER TABLE users ADD COLUMN avatar_url VARCHAR(512);
ALTER TABLE users ADD COLUMN phone VARCHAR(20);
ALTER TABLE users ADD COLUMN address JSONB;
ALTER TABLE users ADD COLUMN deleted_at TIMESTAMPTZ;

-- Update role constraint to allow 'buyer'
ALTER TABLE users DROP CONSTRAINT chk_users_role;
ALTER TABLE users ADD CONSTRAINT chk_users_role CHECK (role IN ('buyer', 'seller', 'admin', 'user'));
