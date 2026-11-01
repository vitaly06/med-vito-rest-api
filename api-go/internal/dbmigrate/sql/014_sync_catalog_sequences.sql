-- Sync all primary key sequences in public schema to MAX(id)
DO $$
DECLARE
    tbl text;
    col text;
    seq text;
    max_id bigint;
BEGIN
    FOR tbl, col IN
        SELECT c.table_name, c.column_name
        FROM information_schema.columns c
        WHERE c.table_schema = 'public'
          AND c.column_default LIKE 'nextval%'
    LOOP
        seq := pg_get_serial_sequence('"' || tbl || '"', col);
        IF seq IS NOT NULL THEN
            EXECUTE format('SELECT COALESCE(MAX(%I), 0) FROM %I', col, tbl) INTO max_id;
            IF max_id > 0 THEN
                PERFORM setval(seq, max_id, true);
            ELSE
                PERFORM setval(seq, 1, false);
            END IF;
        END IF;
    END LOOP;
END $$;

-- Explicitly ensure catalog sequences are at MAX(id)
SELECT setval(pg_get_serial_sequence('"TypeField"', 'id'), GREATEST(COALESCE((SELECT MAX(id) FROM "TypeField"), 1), 1));
SELECT setval(pg_get_serial_sequence('"SubcategotyType"', 'id'), GREATEST(COALESCE((SELECT MAX(id) FROM "SubcategotyType"), 1), 1));
SELECT setval(pg_get_serial_sequence('"SubCategory"', 'id'), GREATEST(COALESCE((SELECT MAX(id) FROM "SubCategory"), 1), 1));
SELECT setval(pg_get_serial_sequence('"Category"', 'id'), GREATEST(COALESCE((SELECT MAX(id) FROM "Category"), 1), 1));
SELECT setval(pg_get_serial_sequence('"ProductFieldValue"', 'id'), GREATEST(COALESCE((SELECT MAX(id) FROM "ProductFieldValue"), 1), 1));
SELECT setval(pg_get_serial_sequence('"ReviewAppeal"', 'id'), GREATEST(COALESCE((SELECT MAX(id) FROM "ReviewAppeal"), 1), 1));
SELECT setval(pg_get_serial_sequence('"Product"', 'id'), GREATEST(COALESCE((SELECT MAX(id) FROM "Product"), 1), 1));
SELECT setval(pg_get_serial_sequence('"User"', 'id'), GREATEST(COALESCE((SELECT MAX(id) FROM "User"), 1), 1));
SELECT setval(pg_get_serial_sequence('"Banner"', 'id'), GREATEST(COALESCE((SELECT MAX(id) FROM "Banner"), 1), 1));
