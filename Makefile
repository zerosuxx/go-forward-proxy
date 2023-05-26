default: help

.PHONY: build

help: ## Show this help
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' -e 's/:.*#/ #/'

install: ## Install the binary
	go get -d ./...
	go install honnef.co/go/tools/cmd/staticcheck@latest

build: ## Build the application
	go build -o build/forward-proxy-bin forward-proxy.go

run: ## Run the application
	go run forward-proxy.go -v

lint: ## Check lint errors
	staticcheck ./...
