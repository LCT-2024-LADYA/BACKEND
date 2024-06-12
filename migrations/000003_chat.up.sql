CREATE TABLE messages
(
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER NOT NULL,
    trainer_id INTEGER NOT NULL,
    message    VARCHAR,
    service_id INTEGER,
    is_to_user BOOLEAN NOT NULL,
    time       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (trainer_id) REFERENCES trainers (id),
    FOREIGN KEY (service_id) REFERENCES services (id)
);
