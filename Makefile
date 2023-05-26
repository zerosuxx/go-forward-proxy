default: help

.PHONY: build

version=`git describe --tags`

help: ## Show this help
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' -e 's/:.*#/ #/'

install: ## Install the binary
	go install
	go install honnef.co/go/tools/cmd/staticcheck@latest

build: ## Build the application
	CGO_ENABLED=0 go build -ldflags="-X 'main.Version=${version}'" -o build/app forward-proxy.go

build-all: ## Build application for supported architectures
	@echo "version: ${version}"
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-X 'main.Version=${version}'" -o build/${BINARY_NAME}-linux-x86_64 forward-proxy.go
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-X 'main.Version=${version}'" -o build/${BINARY_NAME}-linux-aarch64 forward-proxy.go
	GOOS=linux GOARCH=arm CGO_ENABLED=0 go build -ldflags="-X 'main.Version=${version}'" -o build/${BINARY_NAME}-linux-armv7l forward-proxy.go
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-X 'main.Version=${version}'" -o build/${BINARY_NAME}-darwin-x86_64 forward-proxy.go
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-X 'main.Version=${version}'" -o build/${BINARY_NAME}-darwin-aarch64 forward-proxy.go
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-X 'main.Version=${version}'" -o build/${BINARY_NAME}-windows-x86_64.exe forward-proxy.go

run: ## Run the application
	go run forward-proxy.go

lint: ## Check lint errors
	staticcheck ./...