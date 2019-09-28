ENV_PREFIX := GOPATH=$$(pwd)/src
PREFIX := $(ENV_PREFIX)

SOURCES := src/main.go

BUILD_DIR := build
EXECUTABLE := wp
BIN_NAME := $(BUILD_DIR)/$(EXECUTABLE)

.PHONY: all
all: $(BIN_NAME)

$(BIN_NAME): $(SOURCES)
	@mkdir -p $(BUILD_DIR)
	$(PREFIX) go build -o $(BIN_NAME) $(SOURCES)

.PHONY: install
install: $(BIN_NAME)
	cp $(BIN_NAME) /usr/local/bin/$(EXECUTABLE)

.PHONY: test
test:
	$(PREFIX) go test -v 'wpservice'

.PHONY: system-test
system-test: install
	$(PREFIX) go test -v src/main_test.go 

.PHONY: test-cover
test-cover: 
	$(PREFIX) go test -v --coverprofile=coverage.out 'wpservice'

.PHONY: coverage
cover: test-cover
	$(PREFIX) go tool cover -func=coverage.out

.PHONY: pretty-coverage
pretty-coverage: test-cover
	$(PREFIX) go tool cover -html=coverage.out

.PHONY: clean
clean:
	rm coverage.out || true
	rm -rf build || true
	rm -rf $(DEPS_DIR) || true
	rm $(BIN_NAME) || true
