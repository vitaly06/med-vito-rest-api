-- Параметры посылки, способ получения/передачи, статус логистики CDEK
ALTER TABLE "ProductDeal"
    ADD COLUMN IF NOT EXISTS "cdekPackageWeight" INTEGER,
    ADD COLUMN IF NOT EXISTS "cdekPackageLength" INTEGER,
    ADD COLUMN IF NOT EXISTS "cdekPackageWidth" INTEGER,
    ADD COLUMN IF NOT EXISTS "cdekPackageHeight" INTEGER,
    ADD COLUMN IF NOT EXISTS "cdekRecipientMode" TEXT,
    ADD COLUMN IF NOT EXISTS "cdekSellerHandoff" TEXT,
    ADD COLUMN IF NOT EXISTS "cdekFromAddress" TEXT,
    ADD COLUMN IF NOT EXISTS "cdekStatus" TEXT,
    ADD COLUMN IF NOT EXISTS "cdekStatusAt" TIMESTAMP(3);
