# Terraform Provider Valohai

A Terraform provider to manage Valohai projects and resources, including project creation, import, and data access. This provider allows you to automate and manage your Valohai resources directly from your Terraform configuration.

- Manage Valohai projects as Terraform resources
- Import existing Valohai projects
- Access project details via data sources
- Full support for local development and testing

Welcome to the Valohai Terraform provider documentation.

This provider allows you to manage Valohai resources using Terraform.

## Resources

- [valohai_project](resources/valohai_project.md)
- [valohai_team](resources/valohai_team.md)

## Data Sources

### valohai_project

Provides access to a Valohai project by its ID.

#### Example Usage

```hcl
data "valohai_project" "my_project" {
  id = "<project_id>"
}

output "project_name" {
  value = data.valohai_project.my_project.name
}
```

#### Argument Reference

- `id` (Required) – The UUID of the project in Valohai.

#### Attributes Reference

- `id` – The UUID of the project.
- `name` – The name of the project.
- `description` – The project description.
- `owner` – A map with `id` and `username` of the owner.
- `ctime` – Creation time (ISO8601).
- `mtime` – Last modification time (ISO8601).
- `url` – Project URL.
- `urls` – Map of related URLs.
- `execution_count` – Number of executions.
- `running_execution_count` – Number of running executions.
- `queued_execution_count` – Number of queued executions.
- `enabled_endpoint_count` – Number of enabled endpoints.
- `last_execution_ctime` – Last execution creation time.
- `environment_variables` – Map of environment variables.
- `execution_summary` – Map of execution summary counts.
- `repository` – Map with `id`, `url`, and `ref`.
- `tags` – List of tags (each with `project`, `name`, `color`).
- `upload_store_id` – Upload store UUID.
- `read_only` – Whether the project is read-only.
- `yaml_path` – Path to the YAML config.
