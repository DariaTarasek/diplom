INSERT INTO permissions (id, name) VALUES
                                       (1, 'perm:view_admin_pages'),
                                       (2, 'perm:view_doctor_pages'),
                                       (3, 'perm:view_patient_pages'),
                                       (4, 'perm:employee_add');

INSERT INTO role_permission (role_id, permission_id) VALUES
                                                         (1, 1),
                                                         (2, 1),
                                                         (3, 2),
                                                         (4, 3),
                                                         (1, 4);