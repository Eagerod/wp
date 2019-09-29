MAIN_FILE := main.go

BUILD_DIR := build
EXECUTABLE := wp
BIN_NAME := $(BUILD_DIR)/$(EXECUTABLE)

DEPS := $(shell find . -iname "*.go" -and -not -name "*_test.go")

.PHONY: all
all: $(BIN_NAME)

$(BIN_NAME): $(DEPS)
	@mkdir -p $(BUILD_DIR)
	go build -o $(BIN_NAME) $(MAIN_FILE)

.PHONY: install
install: $(BIN_NAME)
	cp $(BIN_NAME) /usr/local/bin/$(EXECUTABLE)

.PHONY: test
test:
	go test -v ./cmd/wpservice

.PHONY: system-test
system-test: install
	go test -v main_test.go 

.PHONY: test-cover
test-cover: 
	go test -v --coverprofile=coverage.out ./cmd/wpservice

.PHONY: coverage
coverage: test-cover
	go tool cover -func=coverage.out

.PHONY: pretty-coverage
pretty-coverage: test-cover
	go tool cover -html=coverage.out

.PHONY: fmt
fmt:
	@go fmt .
	@go fmt ./cmd/wpservice

.PHONY: clean
clean:
	rm -rf coverage.out $(BUILD_DIR)
