CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR UNIQUE NOT NULL,
    password VARCHAR NOT NULL,
    first_name VARCHAR NOT NULL,
    last_name VARCHAR NOT NULL,
    age INTEGER NOT NULL,
    sex INTEGER NOT NULL,
    url_photo VARCHAR NULL
);

CREATE TABLE trainers (
    id SERIAL PRIMARY KEY,
    email VARCHAR UNIQUE NOT NULL,
    password VARCHAR NOT NULL,
    first_name VARCHAR NOT NULL,
    last_name VARCHAR NOT NULL,
    age INTEGER NOT NULL,
    sex INTEGER NOT NULL,
    experience INTEGER NOT NULL,
    quote VARCHAR NULL,
    url_photo VARCHAR NULL
);

CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR UNIQUE NOT NULL
);

CREATE TABLE specializations (
    id SERIAL PRIMARY KEY,
    name VARCHAR UNIQUE NOT NULL
);

CREATE TABLE achievements (
    id SERIAL PRIMARY KEY,
    name VARCHAR UNIQUE NOT NULL
);

CREATE TABLE services (
    trainer_id INTEGER NOT NULL,
    name VARCHAR UNIQUE NOT NULL,
    price INTEGER UNIQUE NOT NULL,
    CONSTRAINT fk_trainer FOREIGN KEY (trainer_id) REFERENCES trainers(id) ON DELETE CASCADE
);

CREATE TABLE schedule_profile (
    id SERIAL PRIMARY KEY,
    trainer_id INTEGER NOT NULL,
    day_of_week VARCHAR(2) NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    CONSTRAINT fk_trainer FOREIGN KEY (trainer_id) REFERENCES trainers(id) ON DELETE CASCADE
);

CREATE TABLE trainers_roles (
    trainer_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,
    PRIMARY KEY (trainer_id, role_id),
    CONSTRAINT fk_trainer FOREIGN KEY (trainer_id) REFERENCES trainers(id) ON DELETE CASCADE,
    CONSTRAINT fk_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

CREATE TABLE trainers_specializations (
    trainer_id INTEGER NOT NULL,
    specialization_id INTEGER NOT NULL,
    PRIMARY KEY (trainer_id, specialization_id),
    CONSTRAINT fk_trainer FOREIGN KEY (trainer_id) REFERENCES trainers(id) ON DELETE CASCADE,
    CONSTRAINT fk_specialization FOREIGN KEY (specialization_id) REFERENCES specializations(id) ON DELETE CASCADE
);

CREATE TABLE trainers_achievements (
    trainer_id INTEGER NOT NULL,
    achievement_id INTEGER NOT NULL,
    PRIMARY KEY (trainer_id, achievement_id),
    CONSTRAINT fk_trainer FOREIGN KEY (trainer_id) REFERENCES trainers(id) ON DELETE CASCADE,
    CONSTRAINT fk_achievement FOREIGN KEY (achievement_id) REFERENCES achievements(id) ON DELETE CASCADE
);
