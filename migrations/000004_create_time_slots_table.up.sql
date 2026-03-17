-- TIME SLOTS
CREATE TABLE time_slots (
    timeslot_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    shift_id BIGINT NOT NULL,
    start_datetime TIMESTAMP NOT NULL,
    end_datetime TIMESTAMP NOT NULL,
    is_peak BOOLEAN DEFAULT FALSE,
    CONSTRAINT fk_timeslot_shift
        FOREIGN KEY (shift_id)
        REFERENCES shifts(shift_id)
        ON DELETE CASCADE,
    CHECK (end_datetime > start_datetime)
);

CREATE INDEX idx_time_slots_shift_id ON time_slots(shift_id);