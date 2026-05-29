ALTER TABLE "User"
ADD COLUMN IF NOT EXISTS "isPhoneVerified" BOOLEAN NOT NULL DEFAULT true;

CREATE TABLE IF NOT EXISTS "OAuthIdentity" (
  "id" SERIAL PRIMARY KEY,
  "provider" TEXT NOT NULL,
  "externalId" TEXT NOT NULL,
  "userId" INTEGER NOT NULL REFERENCES "User"("id") ON DELETE CASCADE,
  "createdAt" TIMESTAMP(3) NOT NULL DEFAULT NOW(),
  "updatedAt" TIMESTAMP(3) NOT NULL DEFAULT NOW(),
  UNIQUE ("provider", "externalId")
);

CREATE INDEX IF NOT EXISTS "OAuthIdentity_userId_idx" ON "OAuthIdentity" ("userId");
