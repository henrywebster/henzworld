-- migrate:up
PRAGMA foreign_keys = ON;

CREATE TABLE posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    --content BLOB NOT NULL,
    created_at TEXT NOT NULL,
    description TEXT NOT NULL
) STRICT;

INSERT INTO posts(id, title, created_at, description) VALUES 
(NULL, 'Trying to Post', '2024-02-02 13:30', 'A post about posting'),
(NULL, 'Another post', '2023-07-12 08:00', 'Thinking about interesting things and commmenting on them.'),
(NULL, 'Doing it again', '2025-12-22 20:27', 'Life, the universe, and everything else.');


-- migrate:down
DROP TABLE posts;

PRAGMA foreign_keys = OFF;
