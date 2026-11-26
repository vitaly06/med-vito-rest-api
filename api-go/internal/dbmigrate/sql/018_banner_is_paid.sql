ALTER TABLE "Banner"
  ADD COLUMN IF NOT EXISTS "isPaid" BOOLEAN NOT NULL DEFAULT false;

UPDATE "Banner" SET "isPaid" = true WHERE "moderateState" = 'APPROVED';
