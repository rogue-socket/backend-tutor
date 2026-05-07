-- Loop 4 migration — users, sessions, and link ownership.

CREATE TABLE users (
    id             BIGSERIAL    PRIMARY KEY,
    email          TEXT         NOT NULL UNIQUE,
    password_hash  TEXT         NOT NULL,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE TABLE sessions (
    id          TEXT         PRIMARY KEY,
    user_id     BIGINT       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at  TIMESTAMPTZ  NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX sessions_user_id_idx  ON sessions(user_id);
CREATE INDEX sessions_expires_idx  ON sessions(expires_at);

-- Add owner_id to links. Keep nullable for now; backfill+contract per Loop 3
-- if you already have data. New apps: assign at insert time and skip backfill.
ALTER TABLE links ADD COLUMN owner_id BIGINT REFERENCES users(id) ON DELETE CASCADE;
CREATE INDEX links_owner_idx ON links(owner_id);
