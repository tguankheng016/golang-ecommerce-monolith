-- +goose Up
-- +goose StatementBegin
CREATE TABLE "public"."user_tokens" (
    "id" bigserial NOT NULL,
    "user_id" bigint NOT NULL,
    "token_key" text NOT NULL,
    "expiration_time" timestamptz NOT NULL,
    PRIMARY KEY ("id")
);
CREATE INDEX IF NOT EXISTS "idx_user_tokens_user_id_expiration_time" ON "public"."user_tokens" ("user_id", "expiration_time");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "public"."idx_user_tokens_user_id_expiration_time";
DROP TABLE "public"."user_tokens";
-- +goose StatementEnd
