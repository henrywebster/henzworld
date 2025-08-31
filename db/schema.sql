CREATE TABLE IF NOT EXISTS "schema_migrations" (version varchar(128) primary key);
CREATE TABLE posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    --content BLOB NOT NULL,
    created_at TEXT NOT NULL,
    description TEXT NOT NULL
) STRICT;
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('20250831041232');
