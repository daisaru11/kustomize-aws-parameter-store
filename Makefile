KIND = AWSParameterStoreSecret
KIND_LOWER = awsparameterstoresecret
API_VERSION = kustomize.daisaru11.dev/v1

CONFIG_DIR = ${HOME}/.config
INSTALL_DIR = $(CONFIG_DIR)/kustomize/plugin/$(API_VERSION)/$(KIND_LOWER)

PLUGIN_BIN = ./build/$(KIND).so
PLUGIN_SRC = ./$(KIND).go

build: $(PLUGIN_BIN)

$(PLUGIN_BIN): $(PLUGIN_SRC)
	go build -buildmode plugin -o $(PLUGIN_BIN) $(PLUGIN_SRC)

install: build
	mkdir -p $(INSTALL_DIR)/
	cp $(PLUGIN_BIN) $(INSTALL_DIR)/

clean:
	rm -f build/*

test:
	go test ./... -v

lint:
	golangci-lint run

.PHONY: test lint clean

.SUFFIXES:
