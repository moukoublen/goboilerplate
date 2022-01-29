VERSION = 0.0.0
VER_FLAGS = \
		-X '${MODULE}/build.Version=${VERSION}' \
		-X '${MODULE}/build.Branch=$(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || true)' \
		-X '${MODULE}/build.Commit=$(shell git rev-parse HEAD 2>/dev/null || true)' \
		-X '${MODULE}/build.CommitShort=$(shell git rev-parse --short HEAD 2>/dev/null || true)' \
		-X '${MODULE}/build.Tag=$(shell git describe --tags 2>/dev/null || true)'

.PHONY: build
build:
	$(GO_EXEC) build -mod=vendor -ldflags "-extldflags -static ${VER_FLAGS}" ${MAINCMD}

.PHONY: image
image:
	$(DOCKER) build . -f $(CURDIR)/build/docker/Dockerfile -t ${NAME}:${IMAGE_TAG}

.PHONY: env
env:
	@echo "Module: ${MODULE}"
	@echo "Name  : ${NAME}"
	@echo "Cmd   : ${MAINCMD}"
	$(GO_EXEC) env
	@echo ""
	@echo ">>> Packages:"
	${PACKAGES}
	@echo ""
	@echo ">>> Folders:"
	${FOLDERS}

.PHONY: mod
mod:
	$(GO_EXEC) mod tidy -compat=1.17
	$(GO_EXEC) mod verify

.PHONY: vendor
vendor:
	$(GO_EXEC) mod vendor
