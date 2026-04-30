CREATE TABLE IF NOT EXISTS "ReviewAppeal" (
    id BIGSERIAL PRIMARY KEY,
    "reviewId" INTEGER NOT NULL REFERENCES "Review"(id) ON DELETE CASCADE,
    "userId" INTEGER NOT NULL REFERENCES "User"(id) ON DELETE CASCADE,
    reason TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'OPEN',
    "moderatorId" INTEGER NULL REFERENCES "User"(id) ON DELETE SET NULL,
    "moderatorNote" TEXT NULL,
    "createdAt" TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMP(3) WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS "ReviewAppeal_review_idx" ON "ReviewAppeal"("reviewId");
CREATE INDEX IF NOT EXISTS "ReviewAppeal_user_idx" ON "ReviewAppeal"("userId");
