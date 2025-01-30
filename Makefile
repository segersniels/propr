.PHONY: build clean run version

BINARY_NAME=propr
VERSION=2.10.0
BUILD_DIR=bin

GOFLAGS=CGO_ENABLED=0
TARGETS=darwin-arm64 darwin-amd64 linux-arm64 linux-amd64
LDFLAGS="-w -s -X main.AppVersion=$(VERSION) -X main.AppName=$(BINARY_NAME)"

build: $(TARGETS)

darwin-arm64:
	$(GOFLAGS) GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 -ldflags $(LDFLAGS)

darwin-amd64:
	$(GOFLAGS) GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 -ldflags $(LDFLAGS)

linux-arm64:
	$(GOFLAGS) GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 -ldflags $(LDFLAGS)

linux-amd64:
	$(GOFLAGS) GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 -ldflags $(LDFLAGS)

clean:
	rm -rf $(BUILD_DIR)

local:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) -ldflags $(LDFLAGS)
	cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

version:
	@echo $(VERSION)
