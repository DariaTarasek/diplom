-- Создание таблицы clinic_weekly_schedule
CREATE TABLE clinic_weekly_schedule (
                                        id SERIAL PRIMARY KEY,
                                        weekday INTEGER NOT NULL CHECK (weekday BETWEEN 0 AND 6),
                                        start_time TIME,
                                        end_time TIME,
                                        slot_duration_minutes INTEGER,
                                        is_day_off BOOLEAN DEFAULT FALSE
);

-- Создание таблицы doctor_weekly_schedule
CREATE TABLE doctor_weekly_schedule (
                                        id SERIAL PRIMARY KEY,
                                        doctor_id INTEGER NOT NULL REFERENCES doctors(user_id) ON DELETE CASCADE,
                                        weekday INTEGER NOT NULL CHECK (weekday BETWEEN 0 AND 6),
                                        start_time TIME,
                                        end_time TIME,
                                        slot_duration_minutes INTEGER,
                                        is_day_off BOOLEAN DEFAULT FALSE
);

-- Создание таблицы clinic_daily_override
CREATE TABLE clinic_daily_override (
                                       id SERIAL PRIMARY KEY,
                                       date DATE NOT NULL,
                                       start_time TIME,
                                       end_time TIME,
                                       slot_duration_minutes INTEGER,
                                       is_day_off BOOLEAN DEFAULT FALSE
);

-- Создание таблицы doctor_daily_override
CREATE TABLE doctor_daily_override (
                                       id SERIAL PRIMARY KEY,
                                       doctor_id INTEGER NOT NULL REFERENCES doctors(user_id) ON DELETE CASCADE,
                                       date DATE NOT NULL,
                                       start_time TIME,
                                       end_time TIME,
                                       slot_duration_minutes INTEGER,
                                       is_day_off BOOLEAN DEFAULT FALSE
);

DROP TABLE weekly_schedule;
DROP TABLE daily_overrides;