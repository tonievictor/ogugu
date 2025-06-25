ALTER TABLE IF EXISTS users
ADD CONSTRAINT users_username_key UNIQUE (username);
