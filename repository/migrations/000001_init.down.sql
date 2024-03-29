DROP TABLE IF EXISTS "cart_items";
DROP TABLE IF EXISTS "customers" CASCADE;
DROP TABLE IF EXISTS "order_items" CASCADE;
DROP TABLE IF EXISTS "orders" CASCADE;
DROP TABLE IF EXISTS "product_images" CASCADE;
DROP TABLE IF EXISTS "products" CASCADE;
DROP TABLE IF EXISTS "stores" CASCADE;
DROP TABLE IF EXISTS "users" CASCADE;

DROP TYPE IF EXISTS "user_role" CASCADE;
DROP TYPE IF EXISTS "order_status" CASCADE;
DROP TYPE IF EXISTS "payment_status" CASCADE;

DROP EXTENSION IF EXISTS "pg_trgm" CASCADE;
DROP EXTENSION IF EXISTS "uuid-ossp" CASCADE;
