CREATE TYPE "PaymentPurpose" AS ENUM ('TOP_UP', 'DEAL_ESCROW');

CREATE TYPE "DealStatus" AS ENUM (
    'CREATED',
    'PAID',
    'SHIPPED',
    'DELIVERED',
    'DISPUTE',
    'COMPLETED',
    'CANCELLED',
    'REFUNDED'
);

ALTER TABLE "Payment"
    ADD COLUMN "purpose" "PaymentPurpose" NOT NULL DEFAULT 'TOP_UP',
    ADD COLUMN "dealId" INTEGER;

CREATE TABLE "ProductDeal" (
    "id" SERIAL NOT NULL,
    "productId" INTEGER NOT NULL,
    "buyerId" INTEGER NOT NULL,
    "sellerId" INTEGER NOT NULL,
    "status" "DealStatus" NOT NULL DEFAULT 'CREATED',
    "productAmount" INTEGER NOT NULL,
    "deliveryCost" INTEGER NOT NULL DEFAULT 0,
    "platformFee" INTEGER NOT NULL,
    "sellerAmount" INTEGER NOT NULL,
    "totalAmount" INTEGER NOT NULL,
    "paymentId" TEXT,
    "orderId" TEXT,
    "paymentUrl" TEXT,
    "cdekTariffCode" INTEGER,
    "cdekTariffName" TEXT,
    "cdekFromCityCode" INTEGER,
    "cdekToCityCode" INTEGER,
    "cdekFromPvzCode" TEXT,
    "cdekToPvzCode" TEXT,
    "cdekOrderUuid" TEXT,
    "cdekTrackNumber" TEXT,
    "disputeReason" TEXT,
    "paidAt" TIMESTAMP(3),
    "shippedAt" TIMESTAMP(3),
    "deliveredAt" TIMESTAMP(3),
    "payoutAvailableAt" TIMESTAMP(3),
    "completedAt" TIMESTAMP(3),
    "cancelledAt" TIMESTAMP(3),
    "refundedAt" TIMESTAMP(3),
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "ProductDeal_pkey" PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX "ProductDeal_paymentId_key" ON "ProductDeal"("paymentId");
CREATE UNIQUE INDEX "ProductDeal_orderId_key" ON "ProductDeal"("orderId");
CREATE INDEX "ProductDeal_productId_idx" ON "ProductDeal"("productId");
CREATE INDEX "ProductDeal_buyerId_idx" ON "ProductDeal"("buyerId");
CREATE INDEX "ProductDeal_sellerId_idx" ON "ProductDeal"("sellerId");
CREATE INDEX "ProductDeal_status_idx" ON "ProductDeal"("status");

ALTER TABLE "ProductDeal" ADD CONSTRAINT "ProductDeal_productId_fkey"
    FOREIGN KEY ("productId") REFERENCES "Product"("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "ProductDeal" ADD CONSTRAINT "ProductDeal_buyerId_fkey"
    FOREIGN KEY ("buyerId") REFERENCES "User"("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "ProductDeal" ADD CONSTRAINT "ProductDeal_sellerId_fkey"
    FOREIGN KEY ("sellerId") REFERENCES "User"("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "Payment" ADD CONSTRAINT "Payment_dealId_fkey"
    FOREIGN KEY ("dealId") REFERENCES "ProductDeal"("id") ON DELETE SET NULL ON UPDATE CASCADE;
