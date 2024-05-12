CREATE TABLE orders(
    id bigserial not null primary key,
    user_id int not null,
    product_id integer[],
    product_count integer[],
    summ int not null,
    FOREIGN KEY (user_id) REFERENCES users(id)
)