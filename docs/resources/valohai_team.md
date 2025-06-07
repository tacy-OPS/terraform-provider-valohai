# Resource: valohai_team

Provides a Valohai team resource. This allows you to create and manage teams in Valohai via Terraform.

## Example Usage

```hcl
resource "valohai_team" "example" {
  name         = "my-team"
  organization = 12345
}
```

## Argument Reference

- `name` (Required) – The name of the team.
- `organization` (Required) – The organization ID for the team.

## Attributes Reference

- `id` – The UUID of the team in Valohai.
- `url` – The URL of the team in Valohai.

## Import

Teams can be imported using the UUID:

```sh
terraform import valohai_team.example <team_uuid>
```
