-- SHIFT TABLE ASSIGNMENTS
CREATE TABLE shift_table_assignments (
    shift_id BIGINT NOT NULL,
    table_id BIGINT NOT NULL,
    staff_id BIGINT NOT NULL,
    PRIMARY KEY (shift_id, table_id),
    CONSTRAINT fk_sta_shift
        FOREIGN KEY (shift_id)
        REFERENCES shifts(shift_id)
        ON DELETE CASCADE,
    CONSTRAINT fk_sta_table
        FOREIGN KEY (table_id)
        REFERENCES tables(table_id),
    CONSTRAINT fk_sta_staff
        FOREIGN KEY (staff_id)
        REFERENCES waitstaff(staff_id)
);