-- +goose Up
-- modify "user_tokens" table
ALTER TABLE "public"."user_tokens" ALTER COLUMN "token_key" TYPE character varying(64);

-- +goose Down
-- reverse: modify "user_tokens" table
ALTER TABLE "public"."user_tokens" ALTER COLUMN "token_key" TYPE text;
