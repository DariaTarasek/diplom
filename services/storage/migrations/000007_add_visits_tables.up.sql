-- 1. ICD-коды (справочник)
CREATE TABLE icd_codes (
                           id SERIAL PRIMARY KEY,
                           code VARCHAR(10) NOT NULL UNIQUE,
                           name TEXT NOT NULL,
                           description TEXT
);

-- 2. Факт проведения приема
CREATE TABLE appointment_visits (
                                    id SERIAL PRIMARY KEY,
                                    appointment_id INTEGER NOT NULL REFERENCES appointments(id) ON DELETE CASCADE,
                                    complaints TEXT,
                                    treatment_plan TEXT,
                                    created_at TIMESTAMP DEFAULT NOW(),
                                    updated_at TIMESTAMP DEFAULT NOW()
);

-- 3. Диагнозы, связанные с визитом
CREATE TABLE appointment_diagnoses (
                                       id SERIAL PRIMARY KEY,
                                       visit_id INTEGER NOT NULL REFERENCES appointment_visits(id) ON DELETE CASCADE,
                                       icd_code_id INTEGER REFERENCES icd_codes(id),
                                       diagnosis_note TEXT
);

-- 4. Проведенные услуги
CREATE TABLE appointment_services (
                                      id SERIAL PRIMARY KEY,
                                      visit_id INTEGER NOT NULL REFERENCES appointment_visits(id) ON DELETE CASCADE,
                                      service_id INTEGER NOT NULL REFERENCES services(id) ON DELETE RESTRICT,
                                      quantity INTEGER NOT NULL DEFAULT 1,
                                      price_per_unit INTEGER,
                                      total_price INTEGER
);

-- 5. Израсходованные материалы
CREATE TABLE appointment_materials (
                                       id SERIAL PRIMARY KEY,
                                       visit_id INTEGER NOT NULL REFERENCES appointment_visits(id) ON DELETE CASCADE,
                                       material_id INTEGER NOT NULL REFERENCES materials(id) ON DELETE RESTRICT,
                                       quantity_used INTEGER NOT NULL
);
