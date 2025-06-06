-- Добавление новых прав
INSERT INTO permissions (id, name) VALUES
                                       (6, 'perm:view_superadmin_pages'),
                                       (7, 'perm:add_daily_override'),
                                       (8, 'perm:manage_patient'),
                                       (9, 'perm:manage_employee'),
                                       (10, 'perm:manage_payments'),
                                       (11, 'perm:manage_patient_medical_notes'),
                                       (12, 'perm:get_patient_visits'),
                                       (13, 'perm:add_consultation'),
                                       (14, 'perm:manage_appointment'),
                                       (15, 'perm:add_patient_doc'),
                                       (16, 'perm:get_patient_doc'),
                                       (17, 'perm:get_statistics'),
                                       (18, 'perm:manage_materials_and_services'),
                                       (19, 'perm:patient_delete');


INSERT INTO role_permission (role_id, permission_id) VALUES
                                                         (1, 6),
                                                         (1, 7),
                                                         (2, 7),
                                                         (1, 8),
                                                         (2, 8),
                                                         (1, 9),
                                                         (1, 10),
                                                         (2, 10),
                                                         (3, 11),
                                                         (3, 12),
                                                         (3, 13),
                                                         (4, 14),
                                                         (2, 14),
                                                         (1, 14),
                                                         (4, 15),
                                                         (4, 16),
                                                         (3, 16),
                                                         (1, 17),
                                                         (1, 18),
                                                         (1, 19);

