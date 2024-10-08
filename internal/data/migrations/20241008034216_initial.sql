-- +goose Up
-- create "roles" table
CREATE TABLE "public"."roles" (
  "id" bigserial NOT NULL,
  "name" character varying(256) NOT NULL,
  "created_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  "created_by" bigint NULL,
  "updated_at" timestamptz NULL,
  "updated_by" bigint NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- create "user_role_permissions" table
CREATE TABLE "public"."user_role_permissions" (
  "id" bigserial NOT NULL,
  "name" character varying(256) NOT NULL,
  "user_id" bigint NULL,
  "role_id" bigint NULL,
  "is_granted" boolean NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_user_role_permissions_role_id" to table: "user_role_permissions"
CREATE INDEX "idx_user_role_permissions_role_id" ON "public"."user_role_permissions" ("role_id");
-- create index "idx_user_role_permissions_user_id" to table: "user_role_permissions"
CREATE INDEX "idx_user_role_permissions_user_id" ON "public"."user_role_permissions" ("user_id");
-- create "user_tokens" table
CREATE TABLE "public"."user_tokens" (
  "id" bigserial NOT NULL,
  "user_id" bigint NULL,
  "token_key" text NULL,
  "expiration_time" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_user_tokens_user_id" to table: "user_tokens"
CREATE INDEX "idx_user_tokens_user_id" ON "public"."user_tokens" ("user_id");
-- create "users" table
CREATE TABLE "public"."users" (
  "id" bigserial NOT NULL,
  "first_name" character varying(64) NULL,
  "last_name" character varying(64) NULL,
  "user_name" character varying(256) NOT NULL,
  "normalized_user_name" character varying(256) NOT NULL,
  "email" character varying(256) NOT NULL,
  "normalized_email" character varying(256) NOT NULL,
  "password" text NOT NULL,
  "security_stamp" text NOT NULL,
  "created_at" timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
  "created_by" bigint NULL,
  "updated_at" timestamptz NULL,
  "updated_by" bigint NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- create "user_roles" table
CREATE TABLE "public"."user_roles" (
  "user_id" bigint NOT NULL,
  "role_id" bigint NOT NULL,
  PRIMARY KEY ("user_id", "role_id"),
  CONSTRAINT "fk_user_roles_role" FOREIGN KEY ("role_id") REFERENCES "public"."roles" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_user_roles_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);

-- +goose Down
-- reverse: create "user_roles" table
DROP TABLE "public"."user_roles";
-- reverse: create "users" table
DROP TABLE "public"."users";
-- reverse: create index "idx_user_tokens_user_id" to table: "user_tokens"
DROP INDEX "public"."idx_user_tokens_user_id";
-- reverse: create "user_tokens" table
DROP TABLE "public"."user_tokens";
-- reverse: create index "idx_user_role_permissions_user_id" to table: "user_role_permissions"
DROP INDEX "public"."idx_user_role_permissions_user_id";
-- reverse: create index "idx_user_role_permissions_role_id" to table: "user_role_permissions"
DROP INDEX "public"."idx_user_role_permissions_role_id";
-- reverse: create "user_role_permissions" table
DROP TABLE "public"."user_role_permissions";
-- reverse: create "roles" table
DROP TABLE "public"."roles";
