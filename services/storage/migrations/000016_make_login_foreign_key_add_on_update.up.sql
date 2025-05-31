
ALTER TABLE doctors
DROP CONSTRAINT fk_doctors_login;

ALTER TABLE administrators
DROP CONSTRAINT fk_administrators_login;

ALTER TABLE patients
DROP CONSTRAINT IF EXISTS fk_patients_login;

ALTER TABLE patients
    ADD CONSTRAINT fk_patients_login
        FOREIGN KEY (phone_number) REFERENCES users(login)
            ON DELETE CASCADE
            ON UPDATE CASCADE;

ALTER TABLE doctors
    ADD CONSTRAINT fk_doctors_login
        FOREIGN KEY (email) REFERENCES users(login)
            ON DELETE CASCADE
            ON UPDATE CASCADE;

ALTER TABLE administrators
    ADD CONSTRAINT fk_administrators_login
        FOREIGN KEY (email) REFERENCES users(login)
            ON DELETE CASCADE
            ON UPDATE CASCADE;

ALTER TABLE patients
ALTER COLUMN phone_number TYPE VARCHAR(256);