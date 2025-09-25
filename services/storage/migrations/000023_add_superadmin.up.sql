INSERT INTO users (id, login, password)
VALUES (1, 'ivanov@mail.ru', '$2a$12$3qYvaeDS/978hNfopWGCu.ZxHWF3lN64fty7DxTMkXL/HyaiKhQyi'); -- jGKmYkSt7Iy1

INSERT INTO user_role (user_id, role_id)
VALUES (1, 1);

INSERT INTO administrators
    (user_id, first_name, second_name, surname, phone_number, email, gender)
VALUES (1, 'Иван', 'Иванов', 'Иванович', '71111111111', 'ivanov@mail.ru', 'м');

