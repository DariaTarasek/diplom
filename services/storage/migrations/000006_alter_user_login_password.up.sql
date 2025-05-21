-- Для таблицы users: делаем login и password NULLABLE
ALTER TABLE users
    ALTER COLUMN login DROP NOT NULL,
    ALTER COLUMN password DROP NOT NULL;

-- Для таблицы patient: делаем phone_number NULLABLE
ALTER TABLE patients
    ALTER COLUMN phone_number DROP NOT NULL;