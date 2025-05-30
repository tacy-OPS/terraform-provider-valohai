# Resource: valohai_project

Provides a Valohai project resource. This allows you to create and manage projects in Valohai via Terraform.

## Example Usage

```hcl
resource "valohai_project" "example" {
  name        = "my-project"
  owner       = "my-org"
  description = "Managed by Terraform"
}
```

## Argument Reference

- `name` (Required) – The name of the project.
- `owner` (Required) – The owner/organization for the project.
- `description` (Optional) – The project description.
- `template_url` (Optional) – The template URL for the project.
- `default_notifications` (Optional) – Enable default notifications (true/false).

## Attributes Reference

- `id` – The UUID of the project in Valohai.

## Import

Projects can be imported using the UUID:

```sh
terraform import valohai_project.example <project_uuid>
```
