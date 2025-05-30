-- Добавить уникальный индекс по дате
CREATE UNIQUE INDEX IF NOT EXISTS clinic_daily_override_date_key
    ON clinic_daily_override (date);

-- Добавить уникальный индекс по doctor_id и date
CREATE UNIQUE INDEX IF NOT EXISTS doctor_daily_override_doctor_id_date_key
    ON doctor_daily_override (doctor_id, date);
