BEGIN;

ALTER TABLE photos DROP CONSTRAINT IF EXISTS fk_photos_users;

ALTER TABLE photos DROP COLUMN IF EXISTS user_id;

COMMIT;