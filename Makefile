.PHONY: dep test lint mock build vendor run

# load env variables from .env
ENV_PATH ?= ./.env
ifneq ($(wildcard $(ENV_PATH)),)
    include .env
    export
endif

# service code
SERVICE = trading

# current version
DOCKER_TAG ?= latest
# docker registry url
DOCKER_URL = gitlab.secreate.dev/cryptocare

# database migrations
DB_ADMIN_USER ?= admin
DB_ADMIN_PASSWORD ?= admin
DB_HOST ?= localhost
DB_NAME ?= cryptocare
DB_AUTH_USER ?= $(SERVICE)
DB_AUTH_PASSWORD ?= $(SERVICE)

DB_DRIVER = postgres
DB_STRING = "user=$(DB_AUTH_USER) password=$(DB_AUTH_PASSWORD) dbname=$(DB_NAME) host=$(DB_HOST) sslmode=disable"
DB_ADMIN_STRING = "user=$(DB_ADMIN_USER) password=$(DB_ADMIN_PASSWORD) dbname=$(DB_NAME) host=$(DB_HOST) sslmode=disable"
DB_INIT_FOLDER = "./src/db/init"
DB_MIG_FOLDER = "./src/db/migrations"

export GOFLAGS=-mod=vendor

vendor:
	go mod vendor

dep:
	go env -w GO111MODULE=on
	go env -w GOPRIVATE=gitlab.secreate.dev/devteam/cryptocarep2p/*
	go mod tidy

test: ## run the tests
	@echo "running tests (skipping integration)"
	go test -count=1 ./...

test-with-coverage: ## run the tests with coverage
	@echo "running tests with coverage file creation (skipping integration)"
	go test -count=1 -coverprofile .testCoverage.txt -v ./...

test-integration: ## run the integration tests
	@echo "running integration tests"
	go test -count=1 -tags integration ./...

check-lint-installed:
	@if ! [ -x "$$(command -v golangci-lint)" ]; then \
		echo "golangci-lint is not installed"; \
		exit 1; \
	fi; \

lint: check-lint-installed
	@echo Running golangci-lint
	golangci-lint --modules-download-mode vendor run --skip-dirs-use-default ./...
	go fmt -mod=vendor ./...

mock: # generate mocks
	@rm -R ./src/mocks 2> /dev/null; \
	find ./src -maxdepth 1 -type d \( ! -path ./*vendor ! -name . \) -exec mockery --output ./src/mocks --all --dir {} \;

build: lint # build library
	@mkdir -p bin
	go build -o bin/ src/cmd/main.go

artifacts: dep vendor mock build swagger ## builds and generates all artifacts

run: ## run the service
	./bin/main

# Database commands ====================================================================================================

check-goose-installed:
	@if ! [ -x "$$(command -v goose)" ]; then \
		echo "goose is not installed"; \
		exit 1; \
	fi; \

db-init-schema:
	GOOSE_DRIVER=$(DB_DRIVER) GOOSE_DBSTRING=$(DB_ADMIN_STRING) goose -dir $(DB_INIT_FOLDER) up

db-status: check-goose-installed
	GOOSE_DRIVER=$(DB_DRIVER) GOOSE_DBSTRING=$(DB_STRING) goose -dir $(DB_MIG_FOLDER) status

db-up: check-goose-installed
	GOOSE_DRIVER=$(DB_DRIVER) GOOSE_DBSTRING=$(DB_STRING) goose -dir $(DB_MIG_FOLDER) up

db-down: check-goose-installed
	GOOSE_DRIVER=$(DB_DRIVER) GOOSE_DBSTRING=$(DB_STRING) goose -dir $(DB_MIG_FOLDER) down

db-create: check-goose-installed
	@if [ -z $(name) ]; then \
      	echo "usage: make db-create name=<you-migration-name>"; \
    else \
		GOOSE_DRIVER=$(DB_DRIVER) GOOSE_DBSTRING=$(DB_STRING) goose -dir $(DB_MIG_FOLDER) create $(name) sql; \
	fi

# CI/CD gitlab commands =================================================================================================

ci-check-mocks:
	@cd ./src
	mv ./mocks ./mocks-init \;
	find . -maxdepth 1 -type d \( ! -path ./*vendor ! -name . \) -exec mockery --all --dir {} \;
	mockshash=$$(find ./mocks -type f -print0 | sort -z | xargs -r0 md5sum | awk '{print $$1}' | md5sum | awk '{print $$1}'); \
	mocksinithash=$$(find ./mocks-init -type f -print0 | sort -z | xargs -r0 md5sum | awk '{print $$1}' | md5sum | awk '{print $$1}'); \
	rm -frd ./mocks-init; \
	echo $$mockshash $$mocksinithash; \
	if ! [ "$$mockshash" = "$$mocksinithash" ] ; then \
	  echo "Mocks should be updated!" ; \
	  cd .. ; \
	  exit 1 ; \
	fi; \
	cd ..

ci-check: ci-check-mocks

ci-build: test-with-coverage build

# infrastructure =======================================================================================================
init-infra:
	@docker network create --driver bridge cc > /dev/null; \
	sudo mkdir -p /var/cryptocare/docker/volumes/aerospike/data; \
	sudo mkdir -p /var/cryptocare/docker/volumes/aerospike/etc; \
	sudo mkdir -p /var/cryptocare/docker/volumes/pg/data; \
	sudo chmod -R g+rwx /var/cryptocare/docker/volumes; \
	sudo chgrp -R 1000 /var/cryptocare/docker/volumes; \
	sudo ls -l /var/cryptocare/docker/volumes/

rm-infra:
	@docker-compose -f ./docker-compose-infra.yml down -v; \
	sudo rm -rfd /var/cryptocare/docker/volumes/pg; \
	sudo rm -rfd /var/cryptocare/docker/volumes/aerospike

run-infra:
	docker-compose -f ./docker-compose-infra.yml up -d --build --remove-orphans

stop-infra:
	docker-compose -f ./docker-compose-infra.yml down

check-infra:
	docker container ls -a --format 'table {{.ID}}\t{{.Names}}{{.Image}}\t{{.Status}}\t{{.Ports}}\t{{.RunningFor}}' | grep -e 'cc-'

# aerospike ============================================================================================================
aql:
	docker exec -it cc-aerospike bash -c aql

# Docker commands =======================================================================================================

docker-build: ## Build the docker images for all services (build inside)
	@echo Building images
	docker build . -f ./Dockerfile -t $(DOCKER_URL)/$(SERVICE):$(DOCKER_TAG) --build-arg SERVICE=$(SERVICE)

docker-build-test: ## Build the docker images for all services (build inside)
	@echo Building images
	docker build . -f ./Dockerfile-test -t $(DOCKER_URL)/$(SERVICE)-test:$(DOCKER_TAG) --build-arg SERVICE=$(SERVICE)

docker-push: docker-build ## Build and push docker images to the repository
	@echo Pushing images
	docker push $(DOCKER_URL)/$(SERVICE):$(DOCKER_TAG)

docker-push-test: docker-build-test ## Build and push docker images to the repository
	@echo Pushing images
	docker push $(DOCKER_URL)/$(SERVICE)-test:$(DOCKER_TAG)

docker-run:
	@echo Running container
	docker run $(DOCKER_URL)/$(SERVICE):$(DOCKER_TAG)

# dev environment ======================================================================================================================

ssh-dev:
	ssh ubuntu@88.99.88.38 -p 5022

# Swagger commands =======================================================================================================

swagger:
	@echo Generating swagger documentation
	swag init -d ./src/cmd,./src/http,./src/kit/http -o ./src/swagger --parseInternal