INSERT INTO "Role" (name)
SELECT x.name
FROM (VALUES
    ('SUPERADMIN'),
    ('ADMIN'),
    ('SENIOR_MODERATOR'),
    ('USER_VERIFIED'),
    ('USER'),
    ('GUEST')
) AS x(name)
WHERE NOT EXISTS (SELECT 1 FROM "Role" r WHERE r.name = x.name);

CREATE TABLE IF NOT EXISTS "ModerationAuditLog" (
    id BIGSERIAL PRIMARY KEY,
    "actorUserId" INTEGER NULL REFERENCES "User"(id) ON DELETE SET NULL,
    "actorRole" TEXT NULL,
    "targetType" TEXT NOT NULL,
    "targetId" BIGINT NOT NULL,
    action TEXT NOT NULL,
    payload JSONB NULL,
    "createdAt" TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS "ModerationAuditLog_createdAt_idx" ON "ModerationAuditLog"("createdAt");
CREATE INDEX IF NOT EXISTS "ModerationAuditLog_target_idx" ON "ModerationAuditLog"("targetType","targetId");

CREATE TABLE IF NOT EXISTS "ModerationAppeal" (
    id BIGSERIAL PRIMARY KEY,
    "productId" INTEGER NOT NULL REFERENCES "Product"(id) ON DELETE CASCADE,
    "userId" INTEGER NOT NULL REFERENCES "User"(id) ON DELETE CASCADE,
    reason TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'OPEN',
    "reviewedByUserId" INTEGER NULL REFERENCES "User"(id) ON DELETE SET NULL,
    "reviewComment" TEXT NULL,
    "createdAt" TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS "ModerationAppeal_product_idx" ON "ModerationAppeal"("productId");
CREATE INDEX IF NOT EXISTS "ModerationAppeal_user_idx" ON "ModerationAppeal"("userId");

