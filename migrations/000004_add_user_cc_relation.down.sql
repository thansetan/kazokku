BEGIN;

ALTER TABLE credit_cards DROP CONSTRAINT IF EXISTS fk_credit_cards_users;

ALTER TABLE credit_cards DROP COLUMN IF EXISTS user_id;

COMMIT;