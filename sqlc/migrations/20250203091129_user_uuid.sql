-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ALTER TABLE users DROP CONSTRAINT users_pkey;
-- ALTER TABLE users ALTER COLUMN id SET DATA TYPE UUID using (uuid_generate_v4());
-- ALTER TABLE users ALTER COLUMN id SET DEFAULT uuid_generate_v4();
-- ALTER TABLE users ADD PRIMARY KEY (id);

ALTER TABLE users ADD COLUMN new_id UUID DEFAULT uuid_generate_v4();
UPDATE users SET new_id = uuid_generate_v4();
ALTER TABLE users DROP COLUMN id;
ALTER TABLE users RENAME COLUMN new_id TO id;
ALTER TABLE users ADD PRIMARY KEY (id);

-- +goose Down
-- ALTER TABLE users DROP CONSTRAINT users_pkey;
-- ALTER TABLE users ALTER COLUMN id SET DATA TYPE BIGSERIAL;
-- ALTER TABLE users ALTER COLUMN id DROP DEFAULT;
-- ALTER TABLE users ADD PRIMARY KEY (id);

ALTER TABLE users DROP CONSTRAINT users_pkey;
ALTER TABLE users ADD COLUMN old_id BIGSERIAL PRIMARY KEY;
CREATE SEQUENCE IF NOT EXISTS users_id_seq START 1;
UPDATE users SET old_id = nextval('users_id_seq');
ALTER TABLE users DROP COLUMN IF EXISTS new_id;
ALTER TABLE users DROP COLUMN id;
ALTER TABLE users RENAME COLUMN old_id TO id;
