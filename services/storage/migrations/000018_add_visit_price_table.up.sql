CREATE TABLE visit_payments (
                                visit_id INTEGER PRIMARY KEY,
                                price INTEGER NOT NULL,
                                status VARCHAR(20) NOT NULL CHECK (status IN ('confirmed', 'unconfirmed')),
                                CONSTRAINT fk_visit
                                    FOREIGN KEY (visit_id)
                                        REFERENCES appointment_visits(id)
                                        ON DELETE CASCADE
);