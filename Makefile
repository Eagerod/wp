GO := go
IMAGEMAGICK := convert

MAIN_FILE := main.go

BUILD_DIR := build
EXECUTABLE := wp
BIN_NAME := $(BUILD_DIR)/$(EXECUTABLE)
INSTALLED_NAME := /usr/local/bin/$(EXECUTABLE)

WP_PACKAGE_DIR := ./cmd/wp
PACKAGE_PATHS := $(WP_PACKAGE_DIR)

AUTOGEN_VERSION_FILENAME=$(WP_PACKAGE_DIR)/version-temp.go

SRC := $(shell find . -iname "*.go" -and -not -name "*_test.go") $(AUTOGEN_VERSION_FILENAME)

TEST_IMAGES_DIR := test_images
TEST_IMAGES := \
	$(TEST_IMAGES_DIR)/square.jpg \
	$(TEST_IMAGES_DIR)/wide.jpg \
	$(TEST_IMAGES_DIR)/tall.jpg

COVERAGE_FILE=coverage.out


.PHONY: all
all: $(BIN_NAME)

$(BIN_NAME): $(SRC)
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BIN_NAME) $(MAIN_FILE)

.PHONY: install isntall
install isntall: $(BIN_NAME)
	cp $(BIN_NAME) $(INSTALLED_NAME)

.PHONY: test
test: $(TEST_IMAGES) $(AUTOGEN_VERSION_FILENAME) $(BIN_NAME)
	@if [ -z $$T ]; then \
		$(GO) test -v ./...; \
	else \
		$(GO) test -v ./... -run $$T; \
	fi


$(COVERAGE_FILE): $(TEST_IMAGES) $(AUTOGEN_VERSION_FILENAME) $(BIN_NAME)
	$(GO) test -v --coverprofile=$(COVERAGE_FILE) ./...

.PHONY: coverage
coverage: $(COVERAGE_FILE)
	$(GO) tool cover -func=$(COVERAGE_FILE)

.INTERMEDIATE: $(AUTOGEN_VERSION_FILENAME)
$(AUTOGEN_VERSION_FILENAME):
	@version="v$$(cat VERSION)" && \
	build="$$(if [ "$$(git describe)" != "$$version" ]; then echo "-$$(git rev-parse --short HEAD)"; fi)" && \
	dirty="$$(if [ ! -z "$$(git diff)" ]; then echo "-dirty"; fi)" && \
	printf "package cmd\n\nconst VersionBuild = \"%s%s%s\"" $$version $$build $$dirty > $@

.PHONY: pretty-coverage
pretty-coverage: $(COVERAGE_FILE)
	$(GO) tool cover -html=$(COVERAGE_FILE)

$(TEST_IMAGES_DIR)/square.jpg:
	mkdir -p $(TEST_IMAGES_DIR)
	$(IMAGEMAGICK) -size 128x128 xc:black $@

$(TEST_IMAGES_DIR)/wide.jpg:
	mkdir -p $(TEST_IMAGES_DIR)
	$(IMAGEMAGICK) -size 256x128 xc:black $@

$(TEST_IMAGES_DIR)/tall.jpg:
	mkdir -p $(TEST_IMAGES_DIR)
	$(IMAGEMAGICK) -size 128x256 xc:black $@

.PHONY: fmt
fmt:
	@$(GO) fmt ./...

.PHONY: clean
clean:
	rm -rf $(COVERAGE_FILE) $(BUILD_DIR) $(TEST_IMAGES_DIR)
