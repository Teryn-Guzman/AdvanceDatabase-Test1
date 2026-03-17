-- RESERVATION TABLE ASSIGNMENTS
CREATE TABLE reservation_table_assignments (
    reservation_id BIGINT NOT NULL,
    table_id BIGINT NOT NULL,
    PRIMARY KEY (reservation_id, table_id),
    CONSTRAINT fk_rta_reservation
        FOREIGN KEY (reservation_id)
        REFERENCES reservations(reservation_id)
        ON DELETE CASCADE,
    CONSTRAINT fk_rta_table
        FOREIGN KEY (table_id)
        REFERENCES tables(table_id)
);

CREATE INDEX idx_rta_table_id ON reservation_table_assignments(table_id);