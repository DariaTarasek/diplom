
ALTER TABLE doctors
    ADD CONSTRAINT fk_doctors_login
        FOREIGN KEY (email) REFERENCES users(login)
            ON DELETE CASCADE;



ALTER TABLE administrators
    ADD CONSTRAINT fk_administrators_login
        FOREIGN KEY (email) REFERENCES users(login)
            ON DELETE CASCADE;


ALTER TABLE patients
    ADD CONSTRAINT fk_patients_login
        FOREIGN KEY (phone_number) REFERENCES users(login)
            ON DELETE CASCADE;
