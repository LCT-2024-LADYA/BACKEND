CREATE TABLE exercises
(
    id                SERIAL PRIMARY KEY,
    name              VARCHAR UNIQUE NOT NULL,
    muscle            VARCHAR NOT NULL,
    additional_muscle VARCHAR NOT NULL,
    type              VARCHAR NOT NULL,
    equipment         VARCHAR NOT NULL,
    difficulty        VARCHAR NOT NULL,
    photos            VARCHAR[] NOT NULL
);

CREATE TABLE trainings
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
    user_id     INTEGER,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE trainings_exercises
(
    step        INTEGER NOT NULL,
    training_id INTEGER NOT NULL,
    exercise_id INTEGER NOT NULL,
    UNIQUE (training_id, exercise_id),
    FOREIGN KEY (training_id) REFERENCES trainings (id) ON DELETE CASCADE,
    FOREIGN KEY (exercise_id) REFERENCES exercises (id) ON DELETE CASCADE
);

CREATE TABLE users_trainings
(
    id          SERIAL PRIMARY KEY,
    user_id     INTEGER NOT NULL,
    training_id INTEGER NOT NULL,
    date        DATE    NOT NULL,
    time_start  TIME    NOT NULL,
    time_end    TIME    NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (training_id) REFERENCES trainings (id) ON DELETE CASCADE
);

CREATE TABLE user_trainings_exercises
(
    id                 SERIAL PRIMARY KEY,
    users_trainings_id INTEGER NOT NULL,
    exercise_id        INTEGER NOT NULL,
    sets               INTEGER NOT NULL,
    reps               INTEGER NOT NULL,
    weight             INTEGER NOT NULL,
    status             BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (users_trainings_id) REFERENCES users_trainings (id) ON DELETE CASCADE,
    FOREIGN KEY (exercise_id) REFERENCES exercises (id) ON DELETE CASCADE
);
