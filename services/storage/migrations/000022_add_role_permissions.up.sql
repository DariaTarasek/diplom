-- Добавление новых прав
INSERT INTO permissions (id, name) VALUES
                                       (20, 'perm:get_materials_and_services');


INSERT INTO role_permission (role_id, permission_id) VALUES
                                                         (1, 20),
                                                         (2, 20),
                                                         (3, 20);

