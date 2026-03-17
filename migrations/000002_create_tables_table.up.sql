CREATE TABLE tables (
    table_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    table_number VARCHAR(20) UNIQUE NOT NULL,
    capacity INT NOT NULL CHECK (capacity > 0),
    location VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE
);