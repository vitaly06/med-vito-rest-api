ALTER TABLE public."Product"
    ADD COLUMN IF NOT EXISTS quantity integer NOT NULL DEFAULT 1;

ALTER TABLE public."Product"
    DROP CONSTRAINT IF EXISTS "Product_quantity_chk";

ALTER TABLE public."Product"
    ADD CONSTRAINT "Product_quantity_chk" CHECK (quantity > 0);
