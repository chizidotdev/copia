CREATE TABLE "cart_items" (
    "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
    "user_id" uuid NOT NULL,
    "product_id" uuid NOT NULL,
    "quantity" integer NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT (now()),
    "updated_at" timestamp NOT NULL DEFAULT (now()),
    UNIQUE ("user_id", "product_id")
);

CREATE INDEX "idx_cart_items_user_id" ON "cart_items" ("user_id");

CREATE INDEX "idx_cart_items_product_id" ON "cart_items" ("product_id");

ALTER TABLE "cart_items" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "cart_items" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON DELETE CASCADE;
