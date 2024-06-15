CREATE TABLE plans
(
    id          SERIAL PRIMARY KEY,
    user_id     INTEGER NOT NULL,
    name        VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE plans_trainings
(
    plan_id     INTEGER NOT NULL,
    training_id INTEGER NOT NULL,
    FOREIGN KEY (plan_id) REFERENCES plans (id) ON DELETE CASCADE,
    FOREIGN KEY (training_id) REFERENCES trainings (id) ON DELETE CASCADE
);
