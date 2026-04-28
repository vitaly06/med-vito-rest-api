ALTER TABLE public."Product"
    ADD COLUMN IF NOT EXISTS "allowReservations" boolean DEFAULT true NOT NULL,
    ADD COLUMN IF NOT EXISTS "reservationHours" integer DEFAULT 24 NOT NULL;

CREATE TABLE IF NOT EXISTS public."ProductReservation" (
    id bigserial PRIMARY KEY,
    "productId" integer NOT NULL,
    "buyerId" integer NOT NULL,
    "sellerId" integer NOT NULL,
    status text NOT NULL DEFAULT 'ACTIVE',
    "hours" integer NOT NULL,
    note text,
    "cancelReason" text,
    "extendedOnce" boolean DEFAULT false NOT NULL,
    "expiresAt" timestamp(3) without time zone NOT NULL,
    "cancelledAt" timestamp(3) without time zone,
    "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updatedAt" timestamp(3) without time zone NOT NULL,
    CONSTRAINT "ProductReservation_status_chk"
        CHECK (status IN ('ACTIVE','CANCELLED_BY_BUYER','CANCELLED_BY_SELLER','EXPIRED','DEAL_CREATED','COMPLETED'))
);

CREATE INDEX IF NOT EXISTS "ProductReservation_product_active_idx"
    ON public."ProductReservation" ("productId", status);

CREATE INDEX IF NOT EXISTS "ProductReservation_buyer_created_idx"
    ON public."ProductReservation" ("buyerId", "createdAt");

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'ProductReservation_productId_fkey') THEN
        ALTER TABLE public."ProductReservation"
            ADD CONSTRAINT "ProductReservation_productId_fkey"
            FOREIGN KEY ("productId") REFERENCES public."Product"(id) ON DELETE CASCADE;
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'ProductReservation_buyerId_fkey') THEN
        ALTER TABLE public."ProductReservation"
            ADD CONSTRAINT "ProductReservation_buyerId_fkey"
            FOREIGN KEY ("buyerId") REFERENCES public."User"(id) ON DELETE CASCADE;
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'ProductReservation_sellerId_fkey') THEN
        ALTER TABLE public."ProductReservation"
            ADD CONSTRAINT "ProductReservation_sellerId_fkey"
            FOREIGN KEY ("sellerId") REFERENCES public."User"(id) ON DELETE CASCADE;
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS public."ReservationUserPenalty" (
    "userId" integer PRIMARY KEY,
    "blockedUntil" timestamp(3) without time zone NOT NULL,
    reason text,
    "updatedAt" timestamp(3) without time zone NOT NULL
);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'ReservationUserPenalty_userId_fkey') THEN
        ALTER TABLE public."ReservationUserPenalty"
            ADD CONSTRAINT "ReservationUserPenalty_userId_fkey"
            FOREIGN KEY ("userId") REFERENCES public."User"(id) ON DELETE CASCADE;
    END IF;
END $$;
