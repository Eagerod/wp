GO := go
IMAGEMAGICK := convert

MAIN_FILE := main.go

BUILD_DIR := build
EXECUTABLE := wp
BIN_NAME := $(BUILD_DIR)/$(EXECUTABLE)

WP_PACKAGE_DIR := ./cmd/wp
PACKAGE_PATHS := $(WP_PACKAGE_DIR)

AUTOGEN_VERSION_FILENAME=$(WP_PACKAGE_DIR)/version-temp.go

SRC := $(shell find . -iname "*.go" -and -not -name "*_test.go") $(AUTOGEN_VERSION_FILENAME)

TEST_IMAGES_DIR := test_images
TEST_IMAGES := \
	$(TEST_IMAGES_DIR)/square.jpg \
	$(TEST_IMAGES_DIR)/wide.jpg \
	$(TEST_IMAGES_DIR)/tall.jpg


.PHONY: all
all: $(BIN_NAME)

$(BIN_NAME): $(SRC)
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BIN_NAME) $(MAIN_FILE)

.PHONY: install
install: $(BIN_NAME)
	cp $(BIN_NAME) /usr/local/bin/$(EXECUTABLE)

.PHONY: test
test: $(TEST_IMAGES)
	@if [ -z $$T ]; then \
		$(GO) test -v $(PACKAGE_PATHS); \
	else \
		$(GO) test -v $(PACKAGE_PATHS) -run $$T; \
	fi

.PHONY: system-test
system-test: install $(TEST_IMAGES)
	@if [ -z $$T ]; then \
		$(GO) test -v main_test.go; \
	else \
		$(GO) test -v main_test.go -run $$T; \
	fi
	

.PHONY: test-cover
test-cover:
	$(GO) test -v --coverprofile=coverage.out $(PACKAGE_PATHS)

.PHONY: coverage
coverage: test-cover
	$(GO) tool cover -func=coverage.out

.INTERMEDIATE: $(AUTOGEN_VERSION_FILENAME)
$(AUTOGEN_VERSION_FILENAME):
	@version="v$$(cat VERSION)" && \
	build="$$(if [ "$$(git describe)" != "$$version" ]; then echo "-$$(git rev-parse --short HEAD)"; fi)" && \
	dirty="$$(if [ ! -z "$$(git diff)" ]; then echo "-dirty"; fi)" && \
	printf "package cmd\n\nconst VersionBuild = \"%s%s%s\"" $$version $$build $$dirty > $@

.PHONY: pretty-coverage
pretty-coverage: test-cover
	$(GO) tool cover -html=coverage.out

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
	@$(GO) fmt .
	@$(GO) fmt $(WP_PACKAGE_DIR)

.PHONY: clean
clean:
	rm -rf coverage.out $(BUILD_DIR) $(TEST_IMAGES_DIR)
