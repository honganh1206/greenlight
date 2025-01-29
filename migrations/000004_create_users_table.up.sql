CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    email citext UNIQUE NOT NULL, -- Case-insensitive
    password_hash bytea NOT NULL, -- binary string
    activated bool NOT NULL,
    version integer NOT NULL DEFAULT 1
);
