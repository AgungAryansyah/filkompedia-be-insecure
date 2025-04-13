CREATE TABLE users (
    id varchar(36) NOT NULL PRIMARY KEY,
    username varchar(36) NOT NULL UNIQUE,
    email varchar(255) NOT NULL UNIQUE,
    password varchar(60) NOT NULL,
    role_id integer NOT NULL DEFAULT 0,
    is_verified boolean NOT NULL DEFAULT FALSE
);

CREATE TABLE roles (
    id INTEGER NOT NULL PRIMARY KEY,
    role varchar(36) NOT NULL UNIQUE
);

INSERT INTO roles (id, role) VALUES 
(0, 'user'),
(1, 'admin');

CREATE TABLE sessions (
    user_id varchar(36),
    token varchar(255) UNIQUE,
    ip_address varchar(16),
    expires_at timestamp,
    user_agent varchar(255),
    device_id varchar(255) UNIQUE, 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
