CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" bigint NOT NULL,
  "profile_pic" text,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  FOREIGN KEY ("account_id") REFERENCES "accounts" ("id")
);

CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "from_account_id" bigint NOT NULL,
  "to_account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id"),
  FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id")
);

CREATE TABLE "blog_categories" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "status" boolean,
  "created_at" timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE "blogs" (
  "id" bigserial PRIMARY KEY,
  "title" varchar NOT NULL,
  "description_article" varchar NOT NULL,
  "category_id" bigint NOT NULL,
  "account_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  FOREIGN KEY ("category_id") REFERENCES "blog_categories" ("id"),
  FOREIGN KEY ("account_id") REFERENCES "accounts" ("id")
);

-- Indexes
CREATE INDEX ON "accounts" ("owner");
CREATE INDEX ON "entries" ("account_id");
CREATE INDEX ON "transfers" ("from_account_id");
CREATE INDEX ON "transfers" ("to_account_id");
CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");
CREATE INDEX ON "blogs" ("title");
CREATE INDEX ON "blogs" ("category_id", "account_id");
CREATE INDEX ON "blog_categories" ("name");

-- Comments
COMMENT ON COLUMN "entries"."amount" IS 'Can be negative or positive value';
COMMENT ON COLUMN "transfers"."amount" IS 'Must be positive';
