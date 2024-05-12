CREATE TABLE reviews(
    id bigserial not null primary key,
    product_id int,
    user_id int not null,
    stars int check (stars>=1 and stars<=5),
    message varchar,
    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);