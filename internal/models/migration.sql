CREATE TABLE products (
    ID SERIAL PRIMARY KEY,
    Name VARCHAR(255) NOT NULL,
    Description TEXT,
    Stock BIT DEFAULT 1,
    Quantity INTEGER NOT NULL,
    Price FLOAT
)
