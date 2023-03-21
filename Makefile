

SHELL := /bin/bash

NAME = dockerfile-json
BUILD_DIR ?= build

GO := $(shell command -v go)
GIT := $(shell command -v git)
TAR := $(shell command -v tar)

GOARCH ?= $(shell $(GO) env GOARCH)
GOOS ?= $(shell $(GO) env GOOS)
VERSION := $(shell $(GIT) describe --always --tags --dirty --first-parent)

LDFLAGS := -s -w -X main.version=$(VERSION)

BUILD_NAME ?= $(NAME)_$(GOOS)_$(GOARCH)
BIN_NAME := $(BUILD_DIR)/$(BUILD_NAME)

ARCHIVE_EXT := txz
ARCHIVE_FLAGS := cJvf

ARCHIVE_NAME ?= $(BUILD_NAME)-$(VERSION).$(ARCHIVE_EXT)
ARCHIVE_TARGET ?=  $(BUILD_DIR)/$(ARCHIVE_NAME)

# ensure that compilation doesn't link against libc
export CGO_ENABLED := 0

.DEFAULT_TARGET := $(BIN_NAME)
.PHONY: build compress test

$(BIN_NAME):
	$(GO) build \
		-a \
		-ldflags "$(LDFLAGS)" \
		-o $(BIN_NAME) \
		-trimpath \
		./

build: $(BIN_NAME)

compress: $(ARCHIVE_TARGET)

$(ARCHIVE_TARGET): $(BIN_NAME)
	$(TAR) -$(ARCHIVE_FLAGS) $(ARCHIVE_TARGET) -C $(BUILD_DIR) $(BUILD_NAME)

test:
	$(GO) test -v ./...

clean:
	@$(RM) -v $(BIN_NAME) $(ARCHIVE_TARGET)

cleaner: clean
	@$(RM) -rv $(BUILD_DIR)
	@$(GO) clean -cache -modcache

cleanest: cleaner
	@$(GIT) clean -fdx

rebuild: clean build
