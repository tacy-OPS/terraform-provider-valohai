PROVIDER_NAME := valohai
NAMESPACE := tacy-ops
VERSION := 0.1.0
OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH := $(shell uname -m | sed 's/x86_64/amd64/' | sed 's/aarch64/arm64/')
BINARY := terraform-provider-$(PROVIDER_NAME)

# Utilise HOME pour Linux/macOS, USERPROFILE pour Windows
ifeq ($(OS),windows_nt)
    PLUGIN_BASE := $(USERPROFILE)
    BINARY_EXT := .exe
else
    PLUGIN_BASE := $(HOME)
    BINARY_EXT :=
endif

PLUGIN_DIR ?= $(PLUGIN_BASE)/.terraform.d/plugins/registry.terraform.io/$(NAMESPACE)/$(PROVIDER_NAME)/$(VERSION)/$(OS)_$(ARCH)

.PHONY: all build install-local clean tfinit tfplan dev check-binary

all: build install-local tfinit

build:
	@echo "🔨 Compilation du provider..."
	go build -o $(BINARY)$(BINARY_EXT)

install-local:
	@echo "📦 Installation du provider localement..."
	mkdir -p $(PLUGIN_DIR)
	@echo "📦 Copie du binaire dans le répertoire de plugins Terraform..."
	cp $(BINARY)$(BINARY_EXT) $(PLUGIN_DIR)/$(BINARY)$(BINARY_EXT)
	@echo "📦 Attribution des permissions d'exécution au binaire..."
	chmod +x $(PLUGIN_DIR)/$(BINARY)$(BINARY_EXT)

tfinit:
	@echo "🚀 Initialisation Terraform..."
	cd example && terraform init

tfplan:
	@echo "🚀 Planification Terraform..."
	cd example && terraform plan

dev: clean build install-local check-binary tfinit tfplan

clean:
	@echo "🧹 Nettoyage..."
	rm -f $(BINARY) $(BINARY).exe
	@echo "🧹 Clean terraform cache..."
	rm -rf ./example/.terraform ./example/.terraform.lock.hcl

check-binary:
	@echo "🔍 Vérification du binaire :"
	cmp -l $(BINARY)$(BINARY_EXT) $(PLUGIN_DIR)/$(BINARY)$(BINARY_EXT) || echo "Les binaires sont différents !"

