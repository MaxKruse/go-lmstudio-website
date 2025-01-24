-- PostgreSQL 
CREATE TABLE IF NOT EXISTS books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    description TEXT,
    image_url VARCHAR(255),
    published_date DATE,
    isbn VARCHAR(255),
    price DECIMAL(4,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

ALTER TABLE books
ADD CONSTRAINT unique_entry_constraint
UNIQUE (title, author, description, published_date);
