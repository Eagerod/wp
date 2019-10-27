MAIN_FILE := main.go

BUILD_DIR := build
EXECUTABLE := wp
BIN_NAME := $(BUILD_DIR)/$(EXECUTABLE)

SRC := $(shell find . -iname "*.go" -and -not -name "*_test.go") cmd/wpservice/version.go
TEST_IMAGES := test_images/square.jpg test_images/wide.jpg test_images/tall.jpg

.PHONY: all
all: $(BIN_NAME)

$(BIN_NAME): $(SRC)
	@mkdir -p $(BUILD_DIR)
	go build -o $(BIN_NAME) $(MAIN_FILE)

.PHONY: install
install: $(BIN_NAME)
	cp $(BIN_NAME) /usr/local/bin/$(EXECUTABLE)

.PHONY: test
test: $(TEST_IMAGES)
	go test -v ./cmd/wpservice

.PHONY: system-test
system-test: install $(TEST_IMAGES)
	go test -v main_test.go 

.PHONY: test-cover
test-cover: 
	go test -v --coverprofile=coverage.out ./cmd/wpservice

.PHONY: coverage
coverage: test-cover
	go tool cover -func=coverage.out

.INTERMEDIATE: cmd/wpservice/version.go
cmd/wpservice/version.go:
	version=$$(cat VERSION) && \
	build=$$(git rev-parse --short HEAD && if [ ! -z "$$(git diff)" ]; then echo "- dirty"; fi) && \
	printf \
		"%s\n\n%s\n%s\n\n%s\n" \
		"package wpservice" \
		"const Version string = \"v$$(printf '%s' $$version)\"" \
		"const Build string = \"$$(printf '%s' $$build)\"" \
		"const VersionBuild string = Version + \"-\" + Build" > $@

.PHONY: pretty-coverage
pretty-coverage: test-cover
	go tool cover -html=coverage.out

test_images/square.jpg:
	mkdir -p test_images
	convert -size 128x128 xc:black $@

test_images/wide.jpg:
	mkdir -p test_images
	convert -size 256x128 xc:black $@

test_images/tall.jpg:
	mkdir -p test_images
	convert -size 128x256 xc:black $@

.PHONY: fmt
fmt:
	@go fmt .
	@go fmt ./cmd/wpservice

.PHONY: clean
clean:
	rm -rf coverage.out $(BUILD_DIR)
