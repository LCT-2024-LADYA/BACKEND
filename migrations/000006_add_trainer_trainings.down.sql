ALTER TABLE trainings
    DROP CONSTRAINT fk_trainer_id,
    DROP COLUMN trainer_id,
    DROP COLUMN wants_public,
    DROP COLUMN is_confirm;
