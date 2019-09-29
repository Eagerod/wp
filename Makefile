SOURCES := main.go

BUILD_DIR := build
EXECUTABLE := wp
BIN_NAME := $(BUILD_DIR)/$(EXECUTABLE)

.PHONY: all
all: $(BIN_NAME)

$(BIN_NAME): $(SOURCES)
	@mkdir -p $(BUILD_DIR)
	go build -o $(BIN_NAME) $(SOURCES)

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
	go test -v --coverprofile=coverage.out 'wpservice'

.PHONY: coverage
cover: test-cover
	go tool cover -func=coverage.out

.PHONY: pretty-coverage
pretty-coverage: test-cover
	go tool cover -html=coverage.out

.PHONY: clean
clean:
	rm -rf coverage.out $(BUILD_DIR)
