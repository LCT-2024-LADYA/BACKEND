CREATE TABLE users_trainers_services
(
    id              SERIAL PRIMARY KEY,
    user_id         INTEGER               NOT NULL,
    trainer_id      INTEGER               NOT NULL,
    service_id      INTEGER,
    is_payed        BOOLEAN DEFAULT FALSE NOT NULL,
    trainer_confirm BOOLEAN,
    user_confirm    BOOLEAN,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (trainer_id) REFERENCES trainers (id) ON DELETE CASCADE,
    FOREIGN KEY (service_id) REFERENCES services (id) ON DELETE CASCADE
);

CREATE TABLE trainer_users_trainers_services
(
    id                         SERIAL PRIMARY KEY,
    users_trainers_services_id INTEGER NOT NULL,
    date                       DATE    NOT NULL,
    time_start                 TIME    NOT NULL,
    time_end                   TIME    NOT NULL,
    FOREIGN KEY (users_trainers_services_id) REFERENCES users_trainers_services (id) ON DELETE CASCADE
);
