default: help

.PHONY: build

help: ## Show this help
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' -e 's/:.*#/ #/'

install: ## Install the binary
	go get -d ./...
	go get -u golang.org/x/lint/golint

build: ## Build the application
	go build -o build/forward-proxy-bin forward-proxy.go

run: ## Run the application
	go run proxy.go -v

lint: ## Check lint errors
	golint -set_exit_status=1 -min_confidence=1.1 ./...