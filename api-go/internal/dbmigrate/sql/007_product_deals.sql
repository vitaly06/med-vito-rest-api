DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'DealStatus') THEN
        CREATE TYPE "DealStatus" AS ENUM (
            'CREATED',
            'PAID',
            'SHIPPED',
            'DELIVERED',
            'COMPLETED',
            'CANCELLED',
            'REFUNDED',
            'DISPUTE'
        );
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS public."ProductDeal" (
    id BIGSERIAL PRIMARY KEY,
    "productId" INTEGER NOT NULL,
    "buyerId" INTEGER NOT NULL,
    "sellerId" INTEGER NOT NULL,
    status "DealStatus" NOT NULL DEFAULT 'CREATED',

    "productAmount" INTEGER NOT NULL DEFAULT 0,
    "deliveryCost" INTEGER NOT NULL DEFAULT 0,
    "platformFee" INTEGER NOT NULL DEFAULT 0,
    "sellerAmount" INTEGER NOT NULL DEFAULT 0,
    "totalAmount" INTEGER NOT NULL DEFAULT 0,

    "paymentId" TEXT NULL,
    "orderId" TEXT NULL,
    "paymentUrl" TEXT NULL,

    "cdekTariffCode" INTEGER NULL,
    "cdekTariffName" TEXT NULL,
    "cdekFromCityCode" INTEGER NULL,
    "cdekToCityCode" INTEGER NULL,
    "cdekFromPvzCode" TEXT NULL,
    "cdekToPvzCode" TEXT NULL,
    "cdekOrderUuid" TEXT NULL,
    "cdekTrackNumber" TEXT NULL,

    "disputeReason" TEXT NULL,

    "paidAt" TIMESTAMP(3) NULL,
    "shippedAt" TIMESTAMP(3) NULL,
    "deliveredAt" TIMESTAMP(3) NULL,
    "payoutAvailableAt" TIMESTAMP(3) NULL,
    "completedAt" TIMESTAMP(3) NULL,
    "cancelledAt" TIMESTAMP(3) NULL,
    "refundedAt" TIMESTAMP(3) NULL,

    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMP(3) NOT NULL DEFAULT NOW(),

    CONSTRAINT "ProductDeal_amounts_chk"
        CHECK (
            "productAmount" >= 0 AND
            "deliveryCost" >= 0 AND
            "platformFee" >= 0 AND
            "sellerAmount" >= 0 AND
            "totalAmount" >= 0
        )
);

CREATE INDEX IF NOT EXISTS "ProductDeal_buyer_created_idx"
    ON public."ProductDeal" ("buyerId", "createdAt" DESC);

CREATE INDEX IF NOT EXISTS "ProductDeal_seller_created_idx"
    ON public."ProductDeal" ("sellerId", "createdAt" DESC);

CREATE INDEX IF NOT EXISTS "ProductDeal_product_created_idx"
    ON public."ProductDeal" ("productId", "createdAt" DESC);

CREATE INDEX IF NOT EXISTS "ProductDeal_status_payout_idx"
    ON public."ProductDeal" (status, "payoutAvailableAt");

CREATE UNIQUE INDEX IF NOT EXISTS "ProductDeal_paymentId_uq"
    ON public."ProductDeal" ("paymentId")
    WHERE "paymentId" IS NOT NULL AND BTRIM("paymentId") <> '';

CREATE UNIQUE INDEX IF NOT EXISTS "ProductDeal_orderId_uq"
    ON public."ProductDeal" ("orderId")
    WHERE "orderId" IS NOT NULL AND BTRIM("orderId") <> '';

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'ProductDeal_productId_fkey') THEN
        ALTER TABLE public."ProductDeal"
            ADD CONSTRAINT "ProductDeal_productId_fkey"
            FOREIGN KEY ("productId") REFERENCES public."Product"(id) ON DELETE CASCADE;
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'ProductDeal_buyerId_fkey') THEN
        ALTER TABLE public."ProductDeal"
            ADD CONSTRAINT "ProductDeal_buyerId_fkey"
            FOREIGN KEY ("buyerId") REFERENCES public."User"(id) ON DELETE CASCADE;
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'ProductDeal_sellerId_fkey') THEN
        ALTER TABLE public."ProductDeal"
            ADD CONSTRAINT "ProductDeal_sellerId_fkey"
            FOREIGN KEY ("sellerId") REFERENCES public."User"(id) ON DELETE CASCADE;
    END IF;
END $$;

