KIND = AWSParameterStoreSecret
KIND_LOWER = awsparameterstoresecret
API_VERSION = kustomize.daisaru11.dev/v1

CONFIG_DIR = ${HOME}/.config
INSTALL_DIR = $(CONFIG_DIR)/kustomize/plugin/$(API_VERSION)/$(KIND_LOWER)
BUILD_DIR = ./kustomize/plugin/$(API_VERSION)/$(KIND_LOWER)

PLUGIN_BIN = $(BUILD_DIR)/$(KIND)
PLUGIN_SRC = main.go

build: $(PLUGIN_BIN)

$(PLUGIN_BIN): $(PLUGIN_SRC)
	mkdir -p $(BUILD_DIR)/
	CGO_ENABLED=0 go build -o $(PLUGIN_BIN)

install: build
	mkdir -p $(INSTALL_DIR)/
	cp $(PLUGIN_BIN) $(INSTALL_DIR)/

clean:
	rm -f $(BUILD_DIR)/*

test:
	go test ./... -v

lint:
	golangci-lint run

.PHONY: test lint clean

.SUFFIXES:
