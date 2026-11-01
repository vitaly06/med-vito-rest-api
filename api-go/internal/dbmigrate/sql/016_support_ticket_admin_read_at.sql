-- Добавляем колонку adminReadAt для отслеживания времени прочтения тикета администратором / поддержкой
ALTER TABLE "SupportTicket"
    ADD COLUMN IF NOT EXISTS "adminReadAt" TIMESTAMP(3) WITHOUT TIME ZONE;

CREATE INDEX IF NOT EXISTS "idx_support_ticket_admin_read_at"
    ON "SupportTicket" ("adminReadAt");

