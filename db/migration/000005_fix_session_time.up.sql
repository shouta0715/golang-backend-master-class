ALTER TABLE "sessions"
ALTER COLUMN "expires_at" TYPE timestamptz,
  ALTER COLUMN "created_at" TYPE timestamptz;