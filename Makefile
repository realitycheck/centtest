# Package Makefile
GO ?= go
DEP ?= dep
DOCKER ?= docker
COMPOSE ?= docker-compose

GO_PACKAGE ?= $(shell echo ${PWD} | sed -e "s/.*src\///")
GO_OUTPUT_FILE ?= $(notdir ${GO_PACKAGE})

DOCKER_TAG ?= ${GO_OUTPUT_FILE}

COMPOSE_FILE ?= configs/docker-compose.yaml

.SHELLFLAGS = -c # Run commands in a -c flag
.PHONY: help clean docker run

${GO_OUTPUT_FILE}: vendor
	go build -o ${GO_OUTPUT_FILE} -a .

vendor:
	$(GO) get -u github.com/golang/dep/cmd/dep
	$(DEP) ensure -vendor-only

help: ## Help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

clean: ## Clean built binary
	rm -f ${GO_OUTPUT_FILE}

docker: ## Create docker's image
	$(DOCKER) build . -t ${DOCKER_TAG} --build-arg GO_PACKAGE=${GO_PACKAGE}

run: ## Run the testing layout
	GO_PACKAGE=${GO_PACKAGE} $(COMPOSE) -f ${COMPOSE_FILE} up ${q}