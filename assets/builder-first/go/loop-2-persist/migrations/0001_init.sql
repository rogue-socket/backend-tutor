-- Loop 2 initial schema.
--
-- Notes:
--   - id BIGSERIAL because BIGINT is the right default in 2026; INT runs out at 2.1B
--   - TIMESTAMPTZ, never TIMESTAMP — read https://wiki.postgresql.org/wiki/Don%27t_Do_This
--   - NOT NULL where a NULL would be a bug, never as a habit
--   - No indexes yet beyond PKs and unique constraints; add them when EXPLAIN says you need them

CREATE TABLE links (
    id          BIGSERIAL    PRIMARY KEY,
    url         TEXT         NOT NULL,
    title       TEXT         NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE TABLE tags (
    id    BIGSERIAL  PRIMARY KEY,
    name  TEXT       NOT NULL UNIQUE
);

CREATE TABLE link_tags (
    link_id  BIGINT  NOT NULL REFERENCES links(id) ON DELETE CASCADE,
    tag_id   BIGINT  NOT NULL REFERENCES tags(id)  ON DELETE CASCADE,
    PRIMARY KEY (link_id, tag_id)
);

-- Index on tag_id so "links with this tag" stays fast as the table grows.
CREATE INDEX link_tags_tag_id_idx ON link_tags(tag_id);
