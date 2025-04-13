CREATE TABLE checkouts (
    id VARCHAR(36) PRIMARY KEY,
    user_id varchar(36) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

ALTER TABLE carts ADD COLUMN checkout_id VARCHAR(36);
ALTER TABLE carts ADD FOREIGN KEY (checkout_id) REFERENCES checkouts(id) ON DELETE CASCADE;