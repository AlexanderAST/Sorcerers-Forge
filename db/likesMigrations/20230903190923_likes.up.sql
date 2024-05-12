CREATE TABLE likes(
    id bigserial not null primary key,
    user_id int not null,
    product_id int not null,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES products(id),
    UNIQUE (user_id, product_id)
)