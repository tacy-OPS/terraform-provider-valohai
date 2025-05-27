# Terraform Provider Valohai

[![Build Status](https://github.com/tacy-OPS/terraform-provider-valohai/actions/workflows/release.yml/badge.svg)](https://github.com/tacy-OPS/terraform-provider-valohai/actions)
[![Go Version](https://img.shields.io/badge/go-1.21%2B-blue.svg)](https://golang.org/doc/go1.21)
[![License](https://img.shields.io/github/license/tacy-OPS/terraform-provider-valohai)](./LICENCE)

A Terraform provider to manage Valohai projects and resources.

## Prerequisites

- Valohai account
- Valohai API token (`VALOHAI_API_TOKEN`)
- Terraform >= 1.0
- Go >= 1.21 (for development)

## Installation

Recommended: Install from the [Terraform Registry](https://registry.terraform.io/providers/tacy-OPS/valohai/latest).

Manual build:
```sh
git clone https://github.com/tacy-OPS/terraform-provider-valohai.git
cd terraform-provider-valohai
go build -o terraform-provider-valohai
```

## Example Usage

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

## Argument Reference

### Provider

- `token` (String, Required): Valohai API token. Can also be set via the `VALOHAI_API_TOKEN` environment variable.

### Resource `valohai_project`

- `name` (String, Required): Project name.
- `owner` (String, Required): Project owner (organization or user).
- `description` (String, Optional): Project description.
- `template_url` (String, Optional): Template URL.
- `default_notifications` (String, Optional): Enable default notifications ("true" or "false").

## Import

```sh
terraform import valohai_project.example <project_id>
```

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