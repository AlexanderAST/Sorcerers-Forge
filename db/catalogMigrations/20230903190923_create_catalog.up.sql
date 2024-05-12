CREATE TABLE product_category (
                                 id bigserial not null primary key,
                                 name varchar not null unique
);

CREATE TABLE products (
                          id bigserial not null primary key,
                          name varchar not null,
                          description varchar not null,
                          price int not null,
                          reviews_mid double precision,
                          reviews_count int,
                          quantity int,
                          work_time varchar,
                          photo varchar,
                          category_id int not null,
                          is_active boolean,
                          FOREIGN KEY (category_id) REFERENCES product_category(id)
);

CREATE TABLE cart_items (
                            id bigserial not null primary key,
                            user_id int not null,
                            product_id int not null,
                            count int not null,
                            photo varchar,
                            FOREIGN KEY (product_id) REFERENCES products(id),
                            FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE favorite_items (
                            id bigserial not null primary key,
                            user_id int not null,
                            product_id int not null,
                            photo varchar,
                            FOREIGN KEY (product_id) REFERENCES products(id),
                            FOREIGN KEY (user_id) REFERENCES users(id)
);