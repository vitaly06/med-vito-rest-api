-- 015_product_paid_content_flags.sql: флаги платного контента для объявлений (более 8 фото, более 2000 символов)
ALTER TABLE "Product"
    ADD COLUMN IF NOT EXISTS "hasPaidExtraPhotos" BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS "hasPaidExtraDesc" BOOLEAN NOT NULL DEFAULT FALSE;

-- Инициализируем существующие товары: если у них уже было > 8 фото или > 2000 символов,
-- помечаем как оплаченные, чтобы с них повторно не списывали деньги.
UPDATE "Product"
SET "hasPaidExtraPhotos" = TRUE
WHERE array_length(images, 1) > 8;

UPDATE "Product"
SET "hasPaidExtraDesc" = TRUE
WHERE length(COALESCE(description, '')) > 2000;

