GO := go
IMAGEMAGICK := convert

MAIN_FILE := main.go

BUILD_DIR := build
EXECUTABLE := wp
BIN_NAME := $(BUILD_DIR)/$(EXECUTABLE)

WP_PACKAGE_DIR := ./cmd/wp
PACKAGE_PATHS := $(WP_PACKAGE_DIR)

SRC := $(shell find . -iname "*.go" -and -not -name "*_test.go") $(WP_PACKAGE_DIR)/version.go
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
	$(GO) test -v $(PACKAGE_PATHS)

.PHONY: system-test
system-test: install $(TEST_IMAGES)
	$(GO) test -v main_test.go

.PHONY: test-cover
test-cover:
	$(GO) test -v --coverprofile=coverage.out $(PACKAGE_PATHS)

.PHONY: coverage
coverage: test-cover
	$(GO) tool cover -func=coverage.out

.INTERMEDIATE: $(WP_PACKAGE_DIR)/version.go
$(WP_PACKAGE_DIR)/version.go:
	@version=$$(cat VERSION) && \
	build=$$(git rev-parse --short HEAD && if [ ! -z "$$(git diff)" ]; then echo "- dirty"; fi) && \
	printf \
		"%s\n\n%s\n%s\n\n%s\n" \
		"package wp" \
		"const Version string = \"v$$(printf '%s' $$version)\"" \
		"const Build string = \"$$(printf '%s' $$build)\"" \
		"const VersionBuild string = Version + \"-\" + Build" > $@

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
