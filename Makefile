
BACKEND_BUILD_FLAGS=

ifeq ($(NO_CACHE), 1)
BACKEND_BUILD_FLAGS += --no-cache
endif

.PHONY: default
default: backend

.PHONY: backend
backend:
	docker compose --env-file=.env -f docker/docker-compose.yml build $(BACKEND_BUILD_FLAGS)
	-docker compose --env-file=.env -f docker/docker-compose.yml up postgresql_db backend
	docker compose --env-file=.env -f docker/docker-compose.yml down

