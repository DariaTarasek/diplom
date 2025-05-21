CREATE TABLE patient_medical_notes (
                                       id SERIAL PRIMARY KEY,
                                       patient_id INTEGER NOT NULL REFERENCES patients(user_id) ON DELETE CASCADE,
                                       type VARCHAR(20) NOT NULL CHECK (type IN ('allergy', 'chronic')),
                                       title TEXT NOT NULL,
                                       created_at TIMESTAMP,
                                       updated_at TIMESTAMP
);