BUILD_DIR=bin

build-server:
	@mkdir -p $(BUILD_DIR)
	go build -o ./$(BUILD_DIR)/server ./cmd/server

build-client:
	@mkdir -p $(BUILD_DIR)
	go build -o ./$(BUILD_DIR)/client ./cmd/client

build: build-server build-client

run-server: build-server
	@./$(BUILD_DIR)/server

run-client: build-client
	@./$(BUILD_DIR)/client