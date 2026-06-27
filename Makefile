BIN_DIR=bin
CLIENT_DIR=./cmd/client
SERVER_DIR=./cmd/server

build-server:
	@echo "Building server..."
	@mkdir -p ${BIN_DIR}
	@go build -o ${BIN_DIR}/server ${SERVER_DIR}/

build-client:
	@echo "Building client..."
	@mkdir -p ${BIN_DIR}
	@go build -o ${BIN_DIR}/client ${CLIENT_DIR}/

build: build-server build-client

run-server: build-server
	@echo "Starting server..."
	@./bin/server

run-client: build-client
	@echo "Starting client..."
	@./bin/client

test:
	@echo "Running tests..."
	@go test -v ./...