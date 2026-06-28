ALTER TABLE categories ADD COLUMN parent_id UUID REFERENCES categories(id);
ALTER TABLE categories ADD COLUMN sort_order INT DEFAULT 0;
ALTER TABLE categories ADD COLUMN image_url VARCHAR(512);
