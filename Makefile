.PHONY: swag_init

swag_init:
	swag init --parseDependency --parseInternal

dev:
	docker compose up