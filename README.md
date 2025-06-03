# Terraform Provider Valohai

[![Build Status](https://github.com/tacy-ops/terraform-provider-valohai/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/tacy-ops/terraform-provider-valohai/actions)
[![Go Version](https://img.shields.io/badge/go-1.21%2B-blue.svg)](https://golang.org/doc/go1.21)
[![License](https://img.shields.io/github/license/tacy-ops/terraform-provider-valohai)](./LICENCE)

A Terraform provider to manage Valohai projects and registry credentials.

## Prerequisites

- Valohai account
- Valohai API token (`VALOHAI_API_TOKEN`)
- Terraform >= 1.0
- Go >= 1.21 (for development)

## Local Development Setup

**Important:** To use your local provider build, you must add the following to your `~/.terraformrc` (Linux/macOS) or `%APPDATA%/terraform.rc` (Windows):

```hcl
provider_installation {
  filesystem_mirror {
    path    = "/absolute/path/to/your/.terraform.d/plugins"
  }
  direct {
    exclude = ["tacy-ops/valohai"]
  }
}
```

This ensures Terraform uses your local provider binary instead of downloading from the registry.

## Installation (Local Development)

1. Clone the repository:
   ```sh
   git clone https://github.com/tacy-ops/terraform-provider-valohai.git
   cd terraform-provider-valohai
   ```
2. Build and install the provider locally:
   ```sh
   make install-local
   ```
   This will build the provider and copy the binary to the correct Terraform plugin directory for local use.

## Usage Example

```hcl
provider "valohai" {
  token = var.valohai_api_token # or via VALOHAI_API_TOKEN
}

resource "valohai_project" "example" {
  name        = "my-project"
  owner       = "my-organization"
  description = "Project managed by Terraform"
}
```

## Importing Existing Resources

You can import existing Valohai resources into Terraform state:

```sh
terraform import valohai_project.example <project_id>
```

## Makefile Commands

- `make build`          – Build the provider binary
- `make install-local`  – Install the provider binary locally for Terraform
- `make tfinit`         – Run `terraform init` in the example directory
- `make tfplan`         – Run `terraform plan` in the example directory
- `make dev`            – Clean, build, install, and run init/plan
- `make clean`          – Remove binaries and Terraform cache
- `make check-binary`   – Compare the built and installed binaries

## Testing

Unit tests:
```sh
go test -v ./tests/resource_projects_unit_test.go
```

Acceptance tests (require a real Valohai token):
```sh
export TF_ACC=1
export VALOHAI_API_TOKEN=your_token
export VALOHAI_OWNER=your_org
go test -v ./tests/...
```

## License

See [LICENCE](./LICENCE).