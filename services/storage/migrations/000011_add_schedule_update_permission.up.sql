INSERT INTO permissions (id, name) VALUES
                                       (5, 'perm:schedule_update');

INSERT INTO role_permission (role_id, permission_id) VALUES
                                                         (1, 5);
