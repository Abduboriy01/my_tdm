create table post_medias (
    id uuid primary key not null,
    type varchar(60),
    link varchar(60),
    post_id uuid,
    FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE 
);

create table post_users (
    id uuid primary key not null,
    first_name varchar(250),
    last_name varchar(250) 
);

CREATE TABLE if not exists posts (
    id uuid primary key not null,
    name varchar(60),
    description varchar(60),
    user_id uuid not null
);