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
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE trainings_exercises
(
    id          SERIAL PRIMARY KEY,
    training_id INTEGER NOT NULL,
    exercise_id INTEGER NOT NULL,
    sets        INTEGER NOT NULL DEFAULT 10,
    reps        INTEGER NOT NULL DEFAULT 4,
    weight      INTEGER NOT NULL DEFAULT 48,
    FOREIGN KEY (training_id) REFERENCES trainings (id),
    FOREIGN KEY (exercise_id) REFERENCES exercises (id)
);

CREATE TABLE users_trainings
(
    id          SERIAL PRIMARY KEY,
    user_id     INTEGER NOT NULL,
    training_id INTEGER NOT NULL,
    date        DATE    NOT NULL,
    time_start  TIME    NOT NULL,
    time_end    TIME    NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (training_id) REFERENCES trainings (id)
);

CREATE TABLE user_trainings_exercises
(
    id                     SERIAL PRIMARY KEY,
    users_trainings_id     INTEGER NOT NULL,
    trainings_exercises_id INTEGER NOT NULL,
    status                 BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (users_trainings_id) REFERENCES users_trainings (id),
    FOREIGN KEY (trainings_exercises_id) REFERENCES trainings_exercises (id)
);
