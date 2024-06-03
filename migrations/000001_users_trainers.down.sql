ALTER TABLE trainers_achievements DROP CONSTRAINT IF EXISTS fk_trainer;
ALTER TABLE trainers_achievements DROP CONSTRAINT IF EXISTS fk_achievement;

ALTER TABLE trainers_specializations DROP CONSTRAINT IF EXISTS fk_trainer;
ALTER TABLE trainers_specializations DROP CONSTRAINT IF EXISTS fk_specialization;

ALTER TABLE trainers_roles DROP CONSTRAINT IF EXISTS fk_trainer;
ALTER TABLE trainers_roles DROP CONSTRAINT IF EXISTS fk_role;

ALTER TABLE schedule_profile DROP CONSTRAINT IF EXISTS fk_trainer;

ALTER TABLE services DROP CONSTRAINT IF EXISTS fk_trainer;


DROP TABLE IF EXISTS trainers_achievements, trainers_specializations, trainers_roles, schedule_profile, services, achievements, specializations, roles, trainers, users;
