DOCKER_EXEC ?= docker
export DOCKER_EXEC

IMAGE_NAME ?= goboilerplate
IMAGE_TAG ?= latest

## https://docs.docker.com/reference/cli/docker/buildx/build/
## --output='type=docker'
## --output='type=image,push=true'
## --platform=linux/arm64
## --platform=linux/amd64,linux/arm64,linux/arm/v7
## --platform=local
## --progress='plain'
## make DOCKER_BUILD_PLATFORM=linux/arm64 image
DOCKER_BUILD_PLATFORM ?= local
DOCKER_BUILD_OUTPUT ?= type=docker
.PHONY: build-image
build-image:
	$(DOCKER_EXEC) buildx build \
		--output='type=docker' \
		--file='$(CURDIR)/build/docker/Dockerfile' \
		--tag='$(IMAGE_NAME):$(IMAGE_TAG)' \
		--platform='$(DOCKER_BUILD_PLATFORM)' \
		--output='$(DOCKER_BUILD_OUTPUT)' \
		.

DOCKER_COMPOSE_EXEC ?= $(DOCKER_EXEC) compose -f $(CURDIR)/deployments/local/docker-compose.yml

.PHONY: compose-up
compose-up:
	$(DOCKER_COMPOSE_EXEC) up --force-recreate --build

.PHONY: compose-up-detach
compose-up-detach:
	$(DOCKER_COMPOSE_EXEC) up --force-recreate --build --detach

.PHONY: compose-down
compose-down:
	$(DOCKER_COMPOSE_EXEC) down --volumes --rmi local --remove-orphans

# If the first target is "compose-exec"
# remove the first argument 'compose-exec' and store the rest in DOCKER_COMPOSE_ARGS
# and ignore the subsequent arguments as make targets.
# (using spaces for indentation)
ifeq (compose-exec,$(firstword $(MAKECMDGOALS)))
    DOCKER_COMPOSE_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
    $(eval $(DOCKER_COMPOSE_ARGS):;@:)
endif

.PHONY: compose-exec
compose-exec:
	$(DOCKER_COMPOSE_EXEC) exec $(DOCKER_COMPOSE_ARGS)
