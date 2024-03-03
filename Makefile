PROJECTNAME := $(shell basename "$(PWD)")
include .env
export $(shell sed 's/=.*//' .env)

# Dockerfile
## gen-images: Generate serivces' image
.PHONY: gen-images
gen-images:
	@docker build --tag auth-svc:$(shell git rev-parse --short HEAD) -f ./build/Dockerfile .

## service-up: Run the all components by deployment/compose.yaml
.PHONY: service-up
service-up:
	@docker-compose  -f ./deployment/compose.yaml --project-directory . up

## service-down: Docker-compose down
.PHONY: service-down
service-down:
	@docker-compose -f ./deployment/compose.yaml --project-directory . down

## dynamodb-up: create dynamodb table
.PHONY: dynamodb-up
dynamodb-up:
	@aws dynamodb create-table --cli-input-json file://deployment/create-table.json --endpoint-url http://localhost:8000

## help: Print usage information
.PHONY: help
help: Makefile
	@echo
	@echo "Choose a command to run in $(PROJECTNAME)"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo