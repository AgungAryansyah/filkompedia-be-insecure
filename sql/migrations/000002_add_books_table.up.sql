CREATE TABLE books (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    introduction TEXT NOT NULL,
    image TEXT NOT NULL,
    file TEXT NOT NULL,
    author VARCHAR(255) NOT NULL,
    release_date DATE NOT NULL,
    price DECIMAL NOT NULL CHECK(price > 1000)
);
 