CREATE TABLE IF NOT EXISTS "SearchQueryStat" (
    id BIGSERIAL PRIMARY KEY,
    query TEXT NOT NULL,
    "region" TEXT NULL,
    "categorySlug" TEXT NULL,
    "subCategorySlug" TEXT NULL,
    "typeSlug" TEXT NULL,
    "resultsCount" INTEGER NOT NULL DEFAULT 0,
    "userId" INTEGER NULL REFERENCES "User"(id) ON DELETE SET NULL,
    "createdAt" TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS "SearchQueryStat_createdAt_idx" ON "SearchQueryStat"("createdAt");
CREATE INDEX IF NOT EXISTS "SearchQueryStat_query_idx" ON "SearchQueryStat"(query);

CREATE TABLE IF NOT EXISTS "TariffFunnelEvent" (
    id BIGSERIAL PRIMARY KEY,
    "userId" INTEGER NULL REFERENCES "User"(id) ON DELETE SET NULL,
    step TEXT NOT NULL,
    "promotionId" INTEGER NULL REFERENCES "Promotion"(id) ON DELETE SET NULL,
    "productId" INTEGER NULL REFERENCES "Product"(id) ON DELETE SET NULL,
    "createdAt" TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS "TariffFunnelEvent_createdAt_idx" ON "TariffFunnelEvent"("createdAt");
CREATE INDEX IF NOT EXISTS "TariffFunnelEvent_step_idx" ON "TariffFunnelEvent"(step);

