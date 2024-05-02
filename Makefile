PROJECT_NAME = resticky
PROJECT_PATH = github.com/osousa/${PROJECT_NAME}
SERVER = server
CLIENT = client
ENV ?= local
LOWER_ENV = `(echo ${ENV} | tr 'A-Z' 'a-z')`

echo-env:
	@echo ""
	@echo "=== Makefile enviroment ==="
	@echo "PROJECT_PATH ${PROJECT_PATH}"
	@echo "PROJECT_NAME ${PROJECT_NAME}"
	@echo "COMPONENT ${COMPONENT}"
	@echo "WRT_METRICS_PUSHGW_URL ${WRT_METRICS_PUSHGW_URL}"
	@echo "ENV  ${ENV}"
	@echo "LOWER_ENV  ${LOWER_ENV}"
	@echo ""

install:
	@echo "=== Installing defining/dependencies ==="
	@export GOPATH="$$HOME/go" && PATH="$$GOPATH/bin:$$PATH"
	@go mod tidy
	@echo "Done"

wire:
ifndef app
	$(error The 'app' variable is not set.)
endif
	@if ! command -v wire >/dev/null 2>&1; then \
		echo "Wire is not installed. Installing dependencies first." && \
		go get -u github.com/google/wire/cmd/wire && \
		make install; \
	fi
	@echo "=== Running wire for ${app} ==="
	@echo $$PATH	
	@echo $$GOPATH
	@echo $$HOME
	@ls $$HOME/go/bin
	@cd ./cmd/${app} && wire

run-server:
	make wire app=server
	cd cmd/server && go run .

build-server:
	@make echo-env
	@make install
	@make wire app=server
	@rm -rf build/server
	@echo "=== Building ${PROJECT_NAME} Server ==="
	@CGO_ENABLED=1 go build -a -o build/server/app "./cmd/server"
	@echo "=== ${PROJECT_NAME} ${COMPONENT} build done ==="
	@echo ""

run-client:
	cd cmd/client && go run .


build-client:
	@make echo-env
	@make install
	@make wire app=client
	@rm -rf build/client
	@echo "=== Building ${PROJECT_NAME} Client ==="
	@CGO_ENABLED=0 go build -a -o build/client/app "./cmd/client"
	@echo "=== ${PROJECT_NAME} ${COMPONENT} build done ==="
	@echo ""


build-all:
	@make build-server
	@make build-client


