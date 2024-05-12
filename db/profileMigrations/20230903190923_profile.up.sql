CREATE TABLE profile (
    id bigserial not null primary key,
    user_id int unique,
    name varchar not null,
    surname varchar not null,
    patronymic varchar,
    contact varchar not null,
    photo varchar,
    FOREIGN KEY (user_id) REFERENCES users(id)
);