CREATE TABLE IF NOT EXISTS feeds  (
    feed_id BLOB PRIMARY KEY,
    title TEXT DEFAULT "" NOT NULL,
    link TEXT NOT NULL UNIQUE,
    image TEXT DEFAULT "" NOT NULL,
    alt_text TEXT DEFAULT "" NOT NULL,
    media_type TEXT DEFAULT "" NOT NULL,
    categories TEXT DEFAULT "" NOT NULL,
    description TEXT DEFAULT "" NOT NULL,
    language TEXT DEFAULT "" NOT NULL,
    last_build_date TEXT DEFAULT CURRENT_TIMESTAMP,
    followed INTEGER DEFAULT 0,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS feed_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    feed_id BLOB NOT NULL,
    feed_name TEXT NOT NULL,
    title TEXT NOT NULL,
    link TEXT NOT NULL UNIQUE,
    description TEXT DEFAULT "" NOT NULL,
    image TEXT DEFAULT "" NOT NULL,
    alt_text TEXT DEFAULT "" NOT NULL,
    media_type TEXT DEFAULT "" NOT NULL,
    categories TEXT DEFAULT "" NOT NULL,
    pub_date TEXT DEFAULT CURRENT_TIMESTAMP,
    guid TEXT UNIQUE,
    saved INTEGER DEFAULT 0,
    read INTEGER DEFAULT 0,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (feed_id) REFERENCES feeds(feed_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS collections (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS collection_items (
    collection_id INTEGER NOT NULL,
    feed_item_id INTEGER NOT NULL,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (feed_item_id) REFERENCES feed_items(id) ON DELETE CASCADE,
    FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE CASCADE,
    UNIQUE (collection_id, feed_item_id)
);
