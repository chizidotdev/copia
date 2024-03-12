CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

CREATE TYPE "user_role" AS ENUM (
  'master',
  'vendor',
  'customer'
);

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "email" varchar UNIQUE NOT NULL,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "image" varchar NOT NULL,
  "google_id" varchar NOT NULL DEFAULT '',
  "role" user_role NOT NULL DEFAULT 'customer',
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "stores" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "user_id" uuid NOT NULL,
  "name" varchar UNIQUE NOT NULL,
  "description" text NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "products" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "store_id" uuid NOT NULL,
  "title" varchar NOT NULL,
  "description" text NOT NULL,
  "price" float NOT NULL,
  "out_of_stock" boolean NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "product_images" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "product_id" uuid NOT NULL,
  "is_primary" boolean NOT NULL DEFAULT false,
  "url" varchar NOT NULL
);

CREATE TABLE "orders" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "user_id" uuid NOT NULL,
  "order_date" timestamp NOT NULL,
  "total_amount" float NOT NULL,
  "status" varchar NOT NULL,
  "payment_status" varchar NOT NULL,
  "shipping_address" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "order_items" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "order_id" uuid NOT NULL,
  "product_id" uuid NOT NULL,
  "quantity" int NOT NULL,
  "unit_price" float NOT NULL,
  "subtotal" float NOT NULL
);

CREATE TABLE "commissions" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "order_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "commission_amount" float NOT NULL,
  "paid_status" varchar NOT NULL
);

CREATE TABLE "links" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "user_id" uuid NOT NULL,
  "unique_link" varchar UNIQUE NOT NULL,
  "link_type" varchar NOT NULL
);

CREATE TABLE "customers" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "store_id" uuid NOT NULL,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "email" varchar NOT NULL,
  "phone" varchar NOT NULL,
  "address" varchar NOT NULL
);

CREATE TABLE "cart_items" (
    "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
    "user_id" uuid NOT NULL,
    "product_id" uuid NOT NULL,
    "quantity" integer NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT (now()),
    "updated_at" timestamp NOT NULL DEFAULT (now()),
    UNIQUE ("user_id", "product_id")
);

CREATE INDEX "idx_stores_user_id" ON "stores" ("user_id");

CREATE INDEX "idx_stores_name_search" ON "stores" USING gin ("name" gin_trgm_ops);

CREATE INDEX "idx_products_store_id" ON "products" ("store_id");

CREATE INDEX "idx_products_title_search" ON "products" USING gin ("title" gin_trgm_ops);

CREATE INDEX "idx_products_description_search" ON "products" USING gin ("description" gin_trgm_ops);

CREATE INDEX "idx_product_images_product_id" ON "product_images" ("product_id");

CREATE INDEX "idx_orders_user_id" ON "orders" ("user_id");

CREATE INDEX "idx_order_items_order_id" ON "order_items" ("order_id");

CREATE INDEX "idx_order_items_product_id" ON "order_items" ("product_id");

CREATE INDEX "idx_commissions_order_id" ON "commissions" ("order_id");

CREATE INDEX "idx_commissions_user_id" ON "commissions" ("user_id");

CREATE INDEX "idx_links_user_id" ON "links" ("user_id");

CREATE INDEX "idx_customers_store_id" ON "customers" ("store_id");

CREATE INDEX "idx_cart_items_user_id" ON "cart_items" ("user_id");

CREATE INDEX "idx_cart_items_product_id" ON "cart_items" ("product_id");

ALTER TABLE "stores" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "order_items" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id") ON DELETE CASCADE;

ALTER TABLE "order_items" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "commissions" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id") ON DELETE CASCADE;

ALTER TABLE "commissions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "links" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "customers" ADD FOREIGN KEY ("store_id") REFERENCES "stores" ("id") ON DELETE CASCADE;

ALTER TABLE "products" ADD FOREIGN KEY ("store_id") REFERENCES "stores" ("id") ON DELETE CASCADE;

ALTER TABLE "product_images" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON DELETE CASCADE;

ALTER TABLE "cart_items" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "cart_items" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON DELETE CASCADE;
