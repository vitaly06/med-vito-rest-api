-- Заполняем все оставшиеся записи Product без даты окончания публикации
UPDATE "Product"
SET "expiresAt" = "createdAt" + INTERVAL '30 days'
WHERE "expiresAt" IS NULL;
