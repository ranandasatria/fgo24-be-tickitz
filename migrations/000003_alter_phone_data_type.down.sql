ALTER TABLE users
ALTER COLUMN phone_number TYPE BIGINT USING phone_number::BIGINT;
