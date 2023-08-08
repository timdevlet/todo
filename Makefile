GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

include .env
export $(shell sed 's/=.*//' .env)

.PHONY: help
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@echo "  ${YELLOW}help            ${RESET} Show this help message"
	@echo "  ${YELLOW}build           ${RESET} Build application binary"
	@echo "  ${YELLOW}setup           ${RESET} Setup local environment"
	@echo "  ${YELLOW}check           ${RESET} Run tests, linters and tidy of the project"
	@echo "  ${YELLOW}test            ${RESET} Run tests only"
	@echo "  ${YELLOW}lint            ${RESET} Run linters via golangci-lint"
	@echo "  ${YELLOW}tidy            ${RESET} Run tidy for go module to remove unused dependencies"
	@echo "  ${YELLOW}run-web         ${RESET} Run web application example"
	@echo "  ${YELLOW}run-cli         ${RESET} Run cli application example"

.PHONY: build
build:
	OS="$(OS)" APP="web" ./hacks/build.sh
	OS="$(OS)" APP="cli" ./hacks/build.sh

.PHONY: build-docker
build-docker: 
	docker build -t timdevlet/todo:latest .

.PHONY: setup
setup:
	cp .env.example .env

.PHONY: check
check: %: tidy lint test

.PHONY: test
test:
	TEST_RUN_ARGS="$(TEST_RUN_ARGS)" TEST_DIR="$(TEST_DIR)" ./hacks/run-tests.sh

.PHONY: lint
lint:
	golangci-lint run

.PHONY: fix
fix:
	golangci-lint run --fix

.PHONY: tidy
tidy:
	go mod tidy -v

.PHONY: run-web
run-web:
	go run ./cmd/web/.

.PHONY: run-cli
run-cli:
	go run ./cmd/cli/.

.PHONY: github
github:
	sudo chmod 666 /var/run/docker.sock
	act