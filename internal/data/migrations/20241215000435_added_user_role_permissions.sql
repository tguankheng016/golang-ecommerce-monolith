-- +goose Up
-- +goose StatementBegin
CREATE TABLE "public"."user_role_permissions" (
    "id" serial NOT NULL,
    "name" character varying(256) NOT NULL,
    "user_id" bigint NULL,
    "role_id" bigint NULL,
    "is_granted" boolean NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);
CREATE INDEX IF NOT EXISTS "idx_user_role_permissions_user_id" ON "public"."user_role_permissions" ("user_id" ASC NULLS LAST);
CREATE INDEX IF NOT EXISTS "idx_user_role_permissions_role_id" ON "public"."user_role_permissions" ("role_id" ASC NULLS LAST);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "public"."idx_user_role_permissions_user_id";
DROP INDEX IF EXISTS "public"."idx_user_role_permissions_role_id";
DROP TABLE "public"."user_role_permissions";
-- +goose StatementEnd
