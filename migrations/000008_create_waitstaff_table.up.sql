-- WAITSTAFF
CREATE TABLE waitstaff (
    staff_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    hire_date DATE,
    is_active BOOLEAN DEFAULT TRUE
);