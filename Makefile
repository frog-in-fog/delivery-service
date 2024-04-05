## builds all docker images
up_build:
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## up: starts all containers in the background without forcing build
up:
	@echo "Starting docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## down: stop docker compose
down:
	@echo "Stopping docker images..."
	docker-compose down
	@echo "Docker stopped!"

sync_db:
	@echo "Database synchronization started..."
	docker cp auth-service:/storage.db ~/golang-dev/delivery-system/auth-service/internal/storage/storage.db
	@echo "Done!"