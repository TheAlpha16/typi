CREATE TABLE IF NOT EXISTS videos (
    vid text PRIMARY KEY,
    title text,
    description text,
    thumbnail text,
    published_at text
);

CREATE TABLE IF NOT EXISTS config (
    key text PRIMARY KEY,
    value text
);
