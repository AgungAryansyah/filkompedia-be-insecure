CREATE TABLE payment_status (
    id INTEGER PRIMARY KEY,
    status VARCHAR(36) NOT NULL
);

INSERT INTO payment_status (id, status) VALUES (0, 'pending');
INSERT INTO payment_status (id, status) VALUES (1, 'accepted');
INSERT INTO payment_status (id, status) VALUES (2, 'deny');
INSERT INTO payment_status (id, status) VALUES (3, 'failed');
INSERT INTO payment_status (id, status) VALUES (4, 'challenge');
INSERT INTO payment_status (id, status) VALUES (5, 'settlement');

CREATE TABLE payments (
    id VARCHAR(36) PRIMARY KEY,
    token VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    checkout_id VARCHAR(36) NOT NULL, 
    total_price DECIMAL NOT NULL CHECK(total_price > 0),
    status_id INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (checkout_id) REFERENCES checkouts(id) ON DELETE CASCADE,
    FOREIGN KEY (status_id) REFERENCES payment_status(id) ON DELETE CASCADE
);