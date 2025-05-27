# Resource: valohai_project

Manages a Valohai project.

## Example Usage

```hcl
provider "valohai" {
  token = var.valohai_api_token
}

resource "valohai_project" "example" {
  name        = "example-project"
  owner       = var.valohai_owner
  description = "Example project managed by Terraform"
}
```

## Argument Reference

- `name` (String, Required): The name of the project.
- `owner` (String, Required): The owner (organization or user) of the project.
- `description` (String, Optional): The description of the project.
- `template_url` (String, Optional): The template URL for the project.
- `default_notifications` (String, Optional): Enable default notifications ("true" or "false").

## Import

Projects can be imported using the project ID:

```sh
terraform import valohai_project.example <project_id>
```

## Environment Variables for Acceptance Tests

- `TF_ACC=1` (enables acceptance tests)
- `VALOHAI_API_TOKEN` (your Valohai API token)
- `VALOHAI_OWNER` (the owner/organization for the project)

## Notes
- Ensure your API token has sufficient permissions to create, update, and delete projects.
- The provider supports full CRUD operations for Valohai projects.
