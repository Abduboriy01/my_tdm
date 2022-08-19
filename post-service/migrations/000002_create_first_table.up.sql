create table post_medias (
    id uuid,
    type varchar(60),
    link varchar(60),
    post_id uuid,
    PRIMARY KEY(id),
    CONSTRAINT fk_posts
    FOREIGN KEY(post_id) 
    REFERENCES posts(id)
    ON DELETE CASCADE 
);