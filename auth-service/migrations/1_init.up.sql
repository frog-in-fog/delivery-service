CREATE TABLE IF NOT EXISTS users (
    id varchar(100) PRIMARY KEY,
    email varchar(50) NOT NULL UNIQUE,
    password_hash text NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_email ON users (email);