CREATE TABLE IF NOT EXISTS users (
    user_id text primary key,
    username text not null unique,
    password_hash text not null,
    name text,
    refresh_token text
);

CREATE UNIQUE INDEX IF NOT EXISTS users_user_id_idx ON users (user_id);
CREATE UNIQUE INDEX IF NOT EXISTS users_username_idx ON users (username);