CREATE TYPE "user_role" AS ENUM (
  'master',
  'vendor',
  'customer'
);

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY,
  "email" varchar NOT NULL,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "image" varchar NOT NULL,
  "password" varchar NOT NULL,
  "google_id" varchar,
  "role" user_role NOT NULL DEFAULT 'customer',
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "stores" (
  "id" uuid PRIMARY KEY,
  "user_id" uuid NOT NULL,
  "store_name" varchar NOT NULL,
  "description" text NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "products" (
  "id" uuid PRIMARY KEY,
  "store_id" uuid NOT NULL,
  "sku" varchar NOT NULL,
  "name" varchar NOT NULL,
  "description" text NOT NULL,
  "price" decimal NOT NULL,
  "stock_quantity" int NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "orders" (
  "id" uuid PRIMARY KEY,
  "user_id" uuid NOT NULL,
  "order_date" timestamp NOT NULL,
  "total_amount" decimal NOT NULL,
  "status" varchar NOT NULL,
  "payment_status" varchar NOT NULL,
  "shipping_address" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "order_items" (
  "id" uuid PRIMARY KEY,
  "order_id" uuid NOT NULL,
  "product_id" uuid NOT NULL,
  "quantity" int NOT NULL,
  "unit_price" decimal NOT NULL,
  "subtotal" decimal NOT NULL
);

CREATE TABLE "commissions" (
  "id" uuid PRIMARY KEY,
  "order_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "commission_amount" decimal NOT NULL,
  "paid_status" varchar NOT NULL
);

CREATE TABLE "links" (
  "id" uuid PRIMARY KEY,
  "user_id" uuid NOT NULL,
  "unique_link" varchar NOT NULL,
  "link_type" varchar NOT NULL
);

CREATE TABLE "customers" (
  "id" uuid PRIMARY KEY,
  "store_id" uuid NOT NULL,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "email" varchar NOT NULL,
  "phone" varchar NOT NULL,
  "address" varchar NOT NULL
);

CREATE INDEX "idx_stores_user_id" ON "stores" ("user_id");

CREATE INDEX "idx_products_store_id" ON "products" ("store_id");

CREATE INDEX "idx_orders_user_id" ON "orders" ("user_id");

CREATE INDEX "idx_order_items_order_id" ON "order_items" ("order_id");

CREATE INDEX "idx_order_items_product_id" ON "order_items" ("product_id");

CREATE INDEX "idx_commissions_order_id" ON "commissions" ("order_id");

CREATE INDEX "idx_commissions_user_id" ON "commissions" ("user_id");

CREATE INDEX "idx_links_user_id" ON "links" ("user_id");

CREATE INDEX "idx_customers_store_id" ON "customers" ("store_id");

ALTER TABLE "stores" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "products" ADD FOREIGN KEY ("store_id") REFERENCES "stores" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "order_items" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "order_items" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");

ALTER TABLE "commissions" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id");

ALTER TABLE "commissions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "links" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "customers" ADD FOREIGN KEY ("store_id") REFERENCES "stores" ("id");
