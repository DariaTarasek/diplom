CREATE TABLE patient_documents (
                                   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                   patient_id integer NOT NULL REFERENCES patients(user_id),
                                   file_name TEXT NOT NULL,
                                   modality TEXT,
                                   study_date DATE,
                                   description TEXT,
                                   storage_path TEXT NOT NULL,
                                   created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);