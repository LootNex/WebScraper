CREATE TABLE IF NOT EXISTS items(
    Id UUID PRIMARY KEY,
    user_id VARCHAR(100),
    link VARCHAR(100),
    product_name VARCHAR(100),
    start_price REAL,
    current_price REAL,
    creation_date TIMESTAMP
);