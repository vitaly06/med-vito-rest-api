-- Очистить lastMessageId у чатов, где сообщения не существуют
UPDATE "Chat" 
SET "lastMessageId" = NULL 
WHERE "lastMessageId" IS NOT NULL 
  AND "lastMessageId" NOT IN (SELECT id FROM "Message");

-- Показать сколько записей было исправлено
SELECT COUNT(*) as fixed_chats FROM "Chat" WHERE "lastMessageId" IS NULL;
