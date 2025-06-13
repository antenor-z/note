COLOR_RESET = \033[0m
COLOR_GREEN = \033[92m
COLOR_YELLOW = \033[93m
COLOR_BLUE = \033[94m
COLOR_BLUE = \033[36m

.DEFAULT_GOAL := help


run: ## Start development server
	go run main.go

admin: ## User administration
	go run cmd/userAdministration/main.go

help:  ## Show this exact message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  ${COLOR_BLUE}%-15s${COLOR_RESET} %s\n", $$1, $$2}' $(MAKEFILE_LIST)