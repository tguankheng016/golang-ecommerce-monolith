# swagger
swagger_app:
	@echo Starting swagger generating
	swag init --parseDependency --parseInternal -g cmd/app/main.go -o docs

# atlas
# execute at root directory with atlas.hcl
atlas_cmd:
	atlas migrate diff migration_name --env gorm