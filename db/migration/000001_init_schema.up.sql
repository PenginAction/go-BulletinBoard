CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "user_str_id" varchar UNIQUE NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "posts" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "text" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);


ALTER TABLE "posts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
