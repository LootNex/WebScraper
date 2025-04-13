CREATE TABLE IF NOT EXISTS users (
    id TEXT NOT NULL PRIMARY KEY,
    telegram_login TEXT NOT NULL,
    login TEXT NOT NULL,
    pass_hash TEXT NOT NULL
);