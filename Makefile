# Package Makefile
GO ?= go
DEP ?= dep
DOCKER ?= docker
DOCKER_COMPOSE ?= docker-compose

COMPOSE_FILE ?= configs/docker-compose.yaml

SRC_DIR ?= ${GOPATH}/src/
PKG_DIR ?= $(shell pwd)

PACKAGE ?= $(subst ${SRC_DIR},,${PKG_DIR})
OUTPUT ?= $(notdir ${PACKAGE})
TAG ?= ${OUTPUT}

.SHELLFLAGS = -c # Run commands in a -c flag
.PHONY: help clean docker run

${OUTPUT}: vendor
	go build -o ${OUTPUT} -a .

vendor:
	$(GO) get -u github.com/golang/dep/cmd/dep
	$(DEP) ensure -vendor-only

help: ## Help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

clean: ## Clean built binary
	rm -f ${OUTPUT}

docker: ## Create docker's image
	$(DOCKER) build . -t ${TAG} --build-arg package=${PACKAGE}

run: ## Run the testing layout
	package=${PACKAGE} $(DOCKER_COMPOSE) -f ${COMPOSE_FILE} up ${flags}