-- Добавляем колонку adminReadAt для отслеживания времени прочтения тикета администратором / поддержкой
ALTER TABLE "SupportTicket"
    ADD COLUMN IF NOT EXISTS "adminReadAt" TIMESTAMP(3) WITHOUT TIME ZONE;

CREATE INDEX IF NOT EXISTS "idx_support_ticket_admin_read_at"
    ON "SupportTicket" ("adminReadAt");

-- Очищаем возможные дублированные префиксы https://torguisam.ru/api в фото пользователей
UPDATE "User"
SET photo = REGEXP_REPLACE(photo, '^https?://torguisam\.ru/api', '')
WHERE photo LIKE '%torguisam.ru/apihttp%';
