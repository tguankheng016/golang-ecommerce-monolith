# RUN
run:
	cd cmd/app && go run .

# GOOSE
# https://github.com/pressly/goose
MIGRATION_NAME = added_user_role_permissions
add_migration:
	cd internal/data/migrations && goose create $(MIGRATION_NAME) sql


# Swaggo
run_swagger:
	swag init --parseDependency --parseInternal -g cmd/app/main.go -o internal/docs/v1