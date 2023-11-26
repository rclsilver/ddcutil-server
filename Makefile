GOCMD=go
SOURCES := $(shell find . -name '*.go')
BINARY_NAME=ddcutil-server
BUILD_DIR=build
INSTALL_DIR=/usr/local/bin

all: build

build: $(BUILD_DIR)/$(BINARY_NAME)

$(BUILD_DIR)/$(BINARY_NAME): $(SOURCES)
	$(GOCMD) build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/ddcutil-server/main.go

$(INSTALL_DIR)/$(BINARY_NAME): $(BUILD_DIR)/$(BINARY_NAME)
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

install: $(INSTALL_DIR)/$(BINARY_NAME)

clean:
	rm -f $(BUILD_DIR)/$(BINARY_NAME)
	sudo rm -f $(INSTALL_DIR)/$(BINARY_NAME)

.PHONY: build install clean
