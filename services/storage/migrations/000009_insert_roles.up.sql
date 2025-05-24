-- insert_roles.sql

INSERT INTO roles (id, name)
VALUES
    (1, 'superadmin'),
    (2, 'admin'),
    (3, 'doctor'),
    (4, 'patient')
ON CONFLICT (id) DO NOTHING;
