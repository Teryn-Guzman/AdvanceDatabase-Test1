-- SPECIAL REQUESTS
CREATE TABLE special_requests (
    request_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    reservation_id BIGINT NOT NULL,
    request_type VARCHAR(100),
    description TEXT,
    CONSTRAINT fk_special_request_reservation
        FOREIGN KEY (reservation_id)
        REFERENCES reservations(reservation_id)
        ON DELETE CASCADE
);

CREATE INDEX idx_special_requests_reservation_id ON special_requests(reservation_id);