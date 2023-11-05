ifneq (,$(wildcard ./.env))
	include .env
	export
endif

VERSION ?= $(shell git describe --tags --always --dirty)
IMAGE_NAME ?= ghcr.io/benc-uk/htmx-go-chat

REPO_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))

.EXPORT_ALL_VARIABLES:
.PHONY: help install-tools run build
.DEFAULT_GOAL := help

# Tools installed locally into repo, don't change
GOLINT_PATH := $(REPO_DIR)/.tools/golangci-lint
AIR_PATH := $(REPO_DIR)/.tools/air

help: ## 💬 This help message :)
	@figlet $@ || true
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install-tools: ## 🔧 Install dev tools into local project tools directory
	@figlet $@ || true
	@$(GOLINT_PATH) > /dev/null 2>&1 || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./.tools
	@$(AIR_PATH) -v > /dev/null 2>&1 || curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b ./.tools

watch: ## 👀 Run the server with reloading
	@figlet $@ || true
	@$(AIR_PATH)

run: ## 🚀 Run the server
	@figlet $@ || true
	go run ./*.go

build: ## 🔨 Build the server
	@figlet $@ || true
	go build -o ./bin/server htmx-go-chat/app

lint: ## 🔍 Lint & format check only, sets exit code on error for CI
	@figlet $@ || true
	$(GOLINT_PATH) run

lint-fix: ## 📝 Lint & format, attempts to fix errors & modify code
	@figlet $@ || true
	$(GOLINT_PATH) run --fix

image: ## 🐳 Build the docker image
	@figlet $@ || true
	docker build . --file build/Dockerfile --tag $(IMAGE_NAME):$(VERSION)

push: ## 📤 Push the docker image to Docker Hub
	@figlet $@ || true
	docker push $(IMAGE_NAME):$(VERSION)

deploy: ## ⛅ Deploy to Azure
	@figlet $@ || true
	@./build/deploy.sh