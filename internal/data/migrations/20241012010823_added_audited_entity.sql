-- +goose Up
-- modify "roles" table
ALTER TABLE "public"."roles" ADD COLUMN "deleted_by" bigint NULL;
-- modify "users" table
ALTER TABLE "public"."users" ADD COLUMN "deleted_by" bigint NULL;

-- +goose Down
-- reverse: modify "users" table
ALTER TABLE "public"."users" DROP COLUMN "deleted_by";
-- reverse: modify "roles" table
ALTER TABLE "public"."roles" DROP COLUMN "deleted_by";
