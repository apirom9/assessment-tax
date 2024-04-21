CREATE TABLE IF NOT EXISTS allowance (
    allowance_type VARCHAR(255) PRIMARY KEY, allowance_amount DECIMAL(10, 2) NOT NULL
);

INSERT INTO
    allowance (
        allowance_type, allowance_amount
    )
VALUES ('personal_default', 60000.00);