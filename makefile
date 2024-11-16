DOCKER_COMPOSE_FILE = ./docker-compose.yml

# Commands
up:
	docker-compose -f $(DOCKER_COMPOSE_FILE) --env-file ./.env up -d --build

down:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

restart:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down
	docker-compose -f $(DOCKER_COMPOSE_FILE) --env-file ./.env up -d --build
