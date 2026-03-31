CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name VARCHAR(25) NOT NULL,
    last_name  VARCHAR(25) NOT NULL,
    username VARCHAR(25) NOT NULL UNIQUE,
    email    VARCHAR(55) NOT NULL UNIQUE,
    phone VARCHAR(20) NOT NULL UNIQUE,
    age SMALLINT CHECK (age >= 0),
    password TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);