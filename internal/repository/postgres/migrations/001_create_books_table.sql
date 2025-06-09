CREATE TABLE IF NOT EXISTS books (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    published INTEGER NOT NULL CHECK (published > 0),
    pages INTEGER NOT NULL CHECK (pages > 0),
    status VARCHAR(20) NOT NULL CHECK (status IN ('available', 'taken')),
    taken_by VARCHAR(255),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
); 