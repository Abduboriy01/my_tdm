CREATE TABLE if not exists posts (
    id uuid not null,
    name varchar(60),
    description varchar(60),
    user_id uuid not null,
    PRIMARY KEY(id)
);