BEGIN;

ALTER TABLE photos ADD COLUMN user_id INT;

ALTER TABLE photos ADD CONSTRAINT fk_photos_users FOREIGN KEY (user_id) REFERENCES users(id);

COMMIT;