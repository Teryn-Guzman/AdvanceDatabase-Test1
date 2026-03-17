-- RESERVATIONS
CREATE TABLE reservations (
    reservation_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    timeslot_id BIGINT NOT NULL,
    party_size INT NOT NULL CHECK (party_size > 0),
    status VARCHAR(50) NOT NULL CHECK (
        status IN ('confirmed', 'cancelled', 'no_show', 'completed')
    ),
    is_walk_in BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    cancelled_at TIMESTAMP,
    notes TEXT,
    CONSTRAINT fk_reservation_customer
        FOREIGN KEY (customer_id)
        REFERENCES customers(customer_id),
    CONSTRAINT fk_reservation_timeslot
        FOREIGN KEY (timeslot_id)
        REFERENCES time_slots(timeslot_id)
);

CREATE INDEX idx_reservations_customer_id ON reservations(customer_id);
CREATE INDEX idx_reservations_timeslot_id ON reservations(timeslot_id);
