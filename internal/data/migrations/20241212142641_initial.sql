-- +goose Up
-- +goose StatementBegin
CREATE TABLE "public"."users" (
    "id" bigserial NOT NULL,
    "first_name" character varying(64) NULL,
    "last_name" character varying(64) NULL,
    "user_name" character varying(256) NOT NULL,
    "normalized_user_name" character varying(256) NOT NULL,
    "email" character varying(256) NOT NULL,
    "normalized_email" character varying(256) NOT NULL,
    "password_hash" text NOT NULL,
    "security_stamp" text NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "created_by" bigint NULL,
    "updated_at" timestamptz NULL,
    "updated_by" bigint NULL,
    "is_deleted" boolean NOT NULL,
    "deleted_at" timestamptz NULL,
    "deleted_by" bigint NULL,
    PRIMARY KEY ("id")
);
CREATE INDEX IF NOT EXISTS "idx_user_is_deleted" ON "public"."users" ("is_deleted");
CREATE INDEX IF NOT EXISTS "idx_user_user_name" ON "public"."users" ("normalized_user_name" ASC NULLS LAST);
CREATE INDEX IF NOT EXISTS "idx_user_email" ON "public"."users" ("normalized_email" ASC NULLS LAST);

CREATE TABLE "public"."roles" (
    "id" bigserial NOT NULL,
    "name" character varying(256) NULL,
    "normalized_name" character varying(256) NULL,
    "is_default" boolean NOT NULL,
    "is_static" boolean NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "created_by" bigint NULL,
    "updated_at" timestamptz NULL,
    "updated_by" bigint NULL,
    "is_deleted" boolean NOT NULL,
    "deleted_at" timestamptz NULL,
    "deleted_by" bigint NULL,
    PRIMARY KEY ("id")
);
CREATE INDEX IF NOT EXISTS "idx_role_is_deleted" ON "public"."roles" ("is_deleted");
CREATE INDEX IF NOT EXISTS "idx_role_name" ON "public"."roles" ("normalized_name" ASC NULLS LAST);

CREATE TABLE "public"."user_roles" (
    "user_id" bigint NOT NULL,
    "role_id" bigint NOT NULL,
    PRIMARY KEY ("user_id", "role_id"),
    CONSTRAINT "fk_user_roles_role" FOREIGN KEY ("role_id") REFERENCES "public"."roles" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
    CONSTRAINT "fk_user_roles_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE INDEX IF NOT EXISTS "idx_user_roles_user_id" ON "public"."user_roles" ("user_id" ASC NULLS LAST);
CREATE INDEX IF NOT EXISTS "idx_user_roles_role_id" ON "public"."user_roles" ("role_id" ASC NULLS LAST);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "public"."idx_user_roles_user_id";
DROP INDEX IF EXISTS "public"."idx_user_roles_role_id";
DROP TABLE "public"."user_roles";

DROP INDEX IF EXISTS "public"."idx_user_is_deleted";
DROP INDEX IF EXISTS "public"."idx_user_user_name";
DROP INDEX IF EXISTS "public"."idx_user_email";
DROP TABLE "public"."users";

DROP INDEX IF EXISTS "public"."idx_role_is_deleted";
DROP INDEX IF EXISTS "public"."idx_role_name";
DROP TABLE "public"."roles";
-- +goose StatementEnd
