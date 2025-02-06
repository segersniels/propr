BINARY_NAME := propr
VERSION := $(shell node -p "require('./package.json').version")
BUILD_DIR := bin

GOFLAGS=CGO_ENABLED=0
TARGETS := darwin-arm64 darwin-amd64 linux-arm64 linux-amd64
LDFLAGS := -w -s -X main.AppVersion=$(VERSION) -X main.AppName=$(BINARY_NAME)

build: $(TARGETS)

$(TARGETS):
	$(GOFLAGS) GOOS=$(word 1,$(subst -, ,$@)) GOARCH=$(word 2,$(subst -, ,$@)) go build -o $(BUILD_DIR)/$(BINARY_NAME)-$@ -ldflags "$(LDFLAGS)"

clean:
	rm -rf $(BUILD_DIR)

dev:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) -ldflags "$(LDFLAGS)"

local: dev
	cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

version:
	@echo $(VERSION)

demo:
	@vhs demo.tape

.PHONY: build clean dev version demo $(TARGETS)
