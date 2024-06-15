-- Миграция для добавления новых полей к таблице trainings
ALTER TABLE trainings
    ADD COLUMN trainer_id INTEGER,
    ADD COLUMN wants_public BOOLEAN DEFAULT FALSE,
    ADD COLUMN is_confirm BOOLEAN DEFAULT FALSE,
    ADD CONSTRAINT fk_trainer_id FOREIGN KEY (trainer_id) REFERENCES trainers (id) ON DELETE CASCADE;
