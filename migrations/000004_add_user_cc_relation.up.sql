BEGIN;

ALTER TABLE credit_cards ADD COLUMN user_id INT UNIQUE;

ALTER TABLE credit_cards ADD CONSTRAINT fk_credit_cards_users FOREIGN KEY (user_id) REFERENCES users(id);

COMMIT;