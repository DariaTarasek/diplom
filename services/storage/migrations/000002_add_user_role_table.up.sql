-- 1. Создание новой таблицы для связи user <-> role
CREATE TABLE user_role (
                            user_id INTEGER NOT NULL,
                            role_id INTEGER NOT NULL,
                            PRIMARY KEY (user_id),
                            FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                            FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

-- 2. Перенос данных из users.role_id в user_roles
INSERT INTO user_role (user_id, role_id)
SELECT id AS user_id, role
FROM users
WHERE role IS NOT NULL;

-- 3. Удаление внешнего ключа и колонки role_id из users
-- (Если внешний ключ был явно назван, укажи имя в DROP CONSTRAINT)
ALTER TABLE users DROP COLUMN role;
