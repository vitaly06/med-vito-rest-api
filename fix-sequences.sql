-- Сброс всех последовательностей для таблиц с автоинкрементом
SELECT setval(pg_get_serial_sequence('"Banner"', 'id'), COALESCE((SELECT MAX(id) FROM "Banner"), 1), true);
SELECT setval(pg_get_serial_sequence('"BannerView"', 'id'), COALESCE((SELECT MAX(id) FROM "BannerView"), 1), true);
SELECT setval(pg_get_serial_sequence('"User"', 'id'), COALESCE((SELECT MAX(id) FROM "User"), 1), true);
SELECT setval(pg_get_serial_sequence('"Product"', 'id'), COALESCE((SELECT MAX(id) FROM "Product"), 1), true);
SELECT setval(pg_get_serial_sequence('"Category"', 'id'), COALESCE((SELECT MAX(id) FROM "Category"), 1), true);
SELECT setval(pg_get_serial_sequence('"Subcategory"', 'id'), COALESCE((SELECT MAX(id) FROM "Subcategory"), 1), true);
SELECT setval(pg_get_serial_sequence('"SubcategoryType"', 'id'), COALESCE((SELECT MAX(id) FROM "SubcategoryType"), 1), true);
SELECT setval(pg_get_serial_sequence('"Review"', 'id'), COALESCE((SELECT MAX(id) FROM "Review"), 1), true);
SELECT setval(pg_get_serial_sequence('"Promotion"', 'id'), COALESCE((SELECT MAX(id) FROM "Promotion"), 1), true);
SELECT setval(pg_get_serial_sequence('"Address"', 'id'), COALESCE((SELECT MAX(id) FROM "Address"), 1), true);
SELECT setval(pg_get_serial_sequence('"Chat"', 'id'), COALESCE((SELECT MAX(id) FROM "Chat"), 1), true);
SELECT setval(pg_get_serial_sequence('"Message"', 'id'), COALESCE((SELECT MAX(id) FROM "Message"), 1), true);
SELECT setval(pg_get_serial_sequence('"Payment"', 'id'), COALESCE((SELECT MAX(id) FROM "Payment"), 1), true);
SELECT setval(pg_get_serial_sequence('"SupportTicket"', 'id'), COALESCE((SELECT MAX(id) FROM "SupportTicket"), 1), true);
SELECT setval(pg_get_serial_sequence('"SupportMessage"', 'id'), COALESCE((SELECT MAX(id) FROM "SupportMessage"), 1), true);
