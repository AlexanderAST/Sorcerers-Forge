ALTER TABLE profile DROP CONSTRAINT profile_user_id_fkey;
ALTER TABLE cart_items DROP CONSTRAINT cart_items_user_id_fkey;
ALTER TABLE favorite_items DROP CONSTRAINT favorite_items_user_id_fkey;
ALTER TABLE likes DROP CONSTRAINT fk_user;
ALTER TABLE reviews DROP CONSTRAINT reviews_user_id_fkey;
ALTER TABLE orders DROP CONSTRAINT orders_user_id_fkey;
DROP TABLE IF EXISTS users;