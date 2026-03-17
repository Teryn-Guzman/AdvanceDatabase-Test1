-- SHIFTS
CREATE TABLE shifts (
    shift_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    shift_name VARCHAR(100),
    shift_date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    CHECK (end_time > start_time)
);