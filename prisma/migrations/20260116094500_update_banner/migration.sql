-- AlterTable
ALTER TABLE "Banner" ADD COLUMN "name" TEXT;

-- AlterTable
ALTER TABLE "User" ALTER COLUMN "freeAdsLimit" SET DEFAULT 6;
