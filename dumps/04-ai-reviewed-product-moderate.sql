--
-- Incremental dump: add AI_REVIEWED to ProductModerate enum
--

ALTER TYPE "ProductModerate" ADD VALUE IF NOT EXISTS 'AI_REVIEWED';

INSERT INTO public._prisma_migrations (
    id,
    checksum,
    finished_at,
    migration_name,
    logs,
    rolled_back_at,
    started_at,
    applied_steps_count
)
SELECT
    '9d58f8ab-a4cc-45ec-af98-2dd2d4da2d89',
    '20c9306ede31cca559b1d6936092c554408648221531c6efb925125db69e1b07',
    CURRENT_TIMESTAMP,
    '20260406000100_add_ai_reviewed_product_moderate',
    '',
    NULL,
    CURRENT_TIMESTAMP,
    1
WHERE EXISTS (
    SELECT 1
    FROM information_schema.tables
    WHERE table_schema = 'public' AND table_name = '_prisma_migrations'
)
AND NOT EXISTS (
    SELECT 1
    FROM public._prisma_migrations
    WHERE migration_name = '20260406000100_add_ai_reviewed_product_moderate'
);
