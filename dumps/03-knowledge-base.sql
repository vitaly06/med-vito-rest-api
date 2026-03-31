--
-- Incremental dump: knowledge base articles
--

CREATE SEQUENCE IF NOT EXISTS public."KnowledgeBaseArticle_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE IF NOT EXISTS public."KnowledgeBaseArticle" (
    id integer NOT NULL DEFAULT nextval('public."KnowledgeBaseArticle_id_seq"'::regclass),
    title text NOT NULL,
    content text NOT NULL,
    "createdAt" timestamp(3) without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" timestamp(3) without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER SEQUENCE public."KnowledgeBaseArticle_id_seq" OWNED BY public."KnowledgeBaseArticle".id;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'KnowledgeBaseArticle_pkey'
    ) THEN
        ALTER TABLE ONLY public."KnowledgeBaseArticle"
            ADD CONSTRAINT "KnowledgeBaseArticle_pkey" PRIMARY KEY (id);
    END IF;
END $$;

INSERT INTO public."KnowledgeBaseArticle" (title, content)
SELECT v.title, v.content
FROM (
    VALUES
        (
            'Как оформить заказ',
            'Выберите нужный товар, откройте карточку, свяжитесь с продавцом через чат и согласуйте условия покупки.'
        ),
        (
            'Как разместить объявление',
            'Авторизуйтесь, перейдите к созданию товара, заполните обязательные поля, добавьте фотографии и отправьте объявление на модерацию.'
        )
) AS v(title, content)
WHERE NOT EXISTS (
    SELECT 1
    FROM public."KnowledgeBaseArticle" k
    WHERE k.title = v.title
);

SELECT pg_catalog.setval(
    'public."KnowledgeBaseArticle_id_seq"',
    COALESCE((SELECT MAX(id) FROM public."KnowledgeBaseArticle"), 1),
    (SELECT COUNT(*) > 0 FROM public."KnowledgeBaseArticle")
);

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
    'b8f7f3e8-3a8a-4de2-b88c-1d5c797f4e69',
    '695072fdffb9c2103a8028ccbd9c6e976837e7821b429af5f879ee0b16c84b6a',
    CURRENT_TIMESTAMP,
    '20260331000100_add_knowledge_base_article',
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
    WHERE migration_name = '20260331000100_add_knowledge_base_article'
);
