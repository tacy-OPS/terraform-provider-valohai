PROVIDER_NAME := valohai
NAMESPACE := hashicorp
OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH := $(shell uname -m | sed 's/x86_64/amd64/' | sed 's/aarch64/arm64/')
PLUGIN_DIR := $(HOME)/.terraform.d/plugins/$(NAMESPACE)/$(PROVIDER_NAME)/$(OS)_$(ARCH)
BINARY := terraform-provider-$(PROVIDER_NAME)

.PHONY: all build install clean tfinit tfplan dev

all: build install tfinit

build:
	@echo "🔨 Compilation du provider..."
	go build -o $(BINARY)

install: build
	@echo "📁 Création du répertoire $(PLUGIN_DIR)..."
	mkdir -p $(PLUGIN_DIR)
	@echo "📦 Copie du binaire dans $(PLUGIN_DIR)..."
	cp $(BINARY) $(PLUGIN_DIR)/

tfinit:
	@echo "🚀 Initialisation Terraform..."
	cd example && terraform init

tfplan:
	@echo "🚀 Planification Terraform..."
	cd example && terraform plan

dev: build install tfinit tfplan

clean:
	@echo "🧹 Nettoyage..."
	rm -f $(BINARY)
