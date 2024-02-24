ALTER TABLE "customers" DROP CONSTRAINT "customers_store_id_fkey";
ALTER TABLE "links" DROP CONSTRAINT "links_user_id_fkey";
ALTER TABLE "commissions" DROP CONSTRAINT "commissions_user_id_fkey";
ALTER TABLE "commissions" DROP CONSTRAINT "commissions_order_id_fkey";
ALTER TABLE "order_items" DROP CONSTRAINT "order_items_product_id_fkey";
ALTER TABLE "order_items" DROP CONSTRAINT "order_items_order_id_fkey";
ALTER TABLE "orders" DROP CONSTRAINT "orders_user_id_fkey";
ALTER TABLE "products" DROP CONSTRAINT "products_store_id_fkey";
ALTER TABLE "stores" DROP CONSTRAINT "stores_user_id_fkey";

DROP INDEX IF EXISTS "idx_customers_store_id";
DROP INDEX IF EXISTS "idx_links_user_id";
DROP INDEX IF EXISTS "idx_commissions_user_id";
DROP INDEX IF EXISTS "idx_commissions_order_id";
DROP INDEX IF EXISTS "idx_order_items_product_id";
DROP INDEX IF EXISTS "idx_order_items_order_id";
DROP INDEX IF EXISTS "idx_orders_user_id";
DROP INDEX IF EXISTS "idx_products_store_id";
DROP INDEX IF EXISTS "idx_stores_user_id";

DROP TABLE IF EXISTS "customers";
DROP TABLE IF EXISTS "links";
DROP TABLE IF EXISTS "commissions";
DROP TABLE IF EXISTS "order_items";
DROP TABLE IF EXISTS "orders";
DROP TABLE IF EXISTS "products";
DROP TABLE IF EXISTS "stores";
DROP TABLE IF EXISTS "users";

DROP TYPE IF EXISTS "user_role";

DROP EXTENSION IF EXISTS "uuid-ossp";
