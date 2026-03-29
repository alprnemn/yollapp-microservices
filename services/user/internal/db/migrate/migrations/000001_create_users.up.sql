CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    first_name VARCHAR(20) NOT NULL,
    last_name  VARCHAR(20) NOT NULL,
    username   VARCHAR(20) UNIQUE NOT NULL,
    email      VARCHAR(55) UNIQUE NOT NULL,
    age        INT CHECK (age >= 0),
    password   TEXT NOT NULL
);