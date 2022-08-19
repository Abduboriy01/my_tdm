CREATE TABLE IF NOT EXISTS users (
    id uuid primary key,
    first_name varchar(255),
    last_name varchar(255),
    email text,
    password text,
    phone_number text,
    PRIMARY KEY(id)
);