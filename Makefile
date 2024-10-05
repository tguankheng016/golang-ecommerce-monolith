# swagger
swagger_app:
	@echo Starting swagger generating
	swag init --parseDependency --parseInternal -g cmd/app/main.go -o docs