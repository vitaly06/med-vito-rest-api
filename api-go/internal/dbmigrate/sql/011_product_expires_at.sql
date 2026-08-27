-- Добавляем колонку expiresAt в таблицу Product
ALTER TABLE "Product"
    ADD COLUMN IF NOT EXISTS "expiresAt" TIMESTAMPTZ;

-- Заполняем существующие записи: дата создания + 30 дней
UPDATE "Product"
SET "expiresAt" = "createdAt" + INTERVAL '30 days'
WHERE "expiresAt" IS NULL;
