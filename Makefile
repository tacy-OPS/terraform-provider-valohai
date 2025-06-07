# Note: For local development, you must add the following to your ~/.terraformrc (Linux/macOS) or %APPDATA%/terraform.rc (Windows):
#
# provider_installation {
#   filesystem_mirror {
#     path    = "<absolute path to your ~/.terraform.d/plugins>"
#   }
#   direct {
#     exclude = ["tacy-ops/valohai"]
#   }
# }
#
# This ensures Terraform uses your local provider binary instead of downloading from the registry.
#
# Example for Linux:
# provider_installation {
#   filesystem_mirror {
#     path    = "/home/<user>/.terraform.d/plugins"
#   }
#   direct {
#     exclude = ["tacy-ops/valohai"]
#   }
# }
#
# Example for Windows:
# provider_installation {
#   filesystem_mirror {
#     path    = "C:\\Users\\<user>\\.terraform.d\\plugins"
#   }
#   direct {
#     exclude = ["tacy-ops/valohai"]
#   }
# }

# Provider build and install variables
PROVIDER_NAME := valohai
NAMESPACE := tacy-ops
VERSION := 0.1.0
OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH := $(shell uname -m | sed 's/x86_64/amd64/' | sed 's/aarch64/arm64/')
BINARY := terraform-provider-$(PROVIDER_NAME)

# Use HOME for Linux/macOS, USERPROFILE for Windows
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
	@echo "🔨 Building the provider..."
	go build -o $(BINARY)$(BINARY_EXT)

install-local:
	@echo "📦 Installing the provider locally..."
	mkdir -p $(PLUGIN_DIR)
	@echo "📦 Copying the binary to the Terraform plugins directory..."
	cp $(BINARY)$(BINARY_EXT) $(PLUGIN_DIR)/$(BINARY)$(BINARY_EXT)
	@echo "📦 Setting execution permissions on the binary..."
	chmod +x $(PLUGIN_DIR)/$(BINARY)$(BINARY_EXT)

tfinit:
	@echo "🚀 Running terraform init in ./example..."
	cd example && terraform init

tfplan:
	@echo "🚀 Running terraform plan in ./example..."
	cd example && terraform plan

dev: clean build install-local check-binary tfinit tfplan

test:
	@echo "🧪 Running tests..."
	go test -v ./...

clean:
	@echo "🧹 Cleaning binaries..."
	rm -f $(BINARY) $(BINARY).exe
	@echo "🧹 Cleaning terraform cache..."
	rm -rf ./example/.terraform ./example/.terraform.lock.hcl

check-binary:
	@echo "🔍 Checking the binary:"
	cmp -l $(BINARY)$(BINARY_EXT) $(PLUGIN_DIR)/$(BINARY)$(BINARY_EXT) || echo "Binaries are different!"

