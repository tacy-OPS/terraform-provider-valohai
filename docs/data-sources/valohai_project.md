# Data Source: valohai_project

Use this data source to retrieve information about an existing Valohai project by its ID.

## Example Usage

```hcl
provider "valohai" {
  token = "${var.valohai_token}"
}

data "valohai_project" "example" {
  id = "0197420b-164b-983d-0e23-528b6253a2b5"
}

output "project_name" {
  value = data.valohai_project.example.name
}
```

## Argument Reference

- `id` (String, Required): The ID of the project to look up.

## Attributes Reference

The following attributes are exported:

- `id` – The project ID.
- `name` – The name of the project.
- `description` – The project description.
- `owner` – A map with the owner's `id` and `username`.
- `ctime` – Creation time.
- `mtime` – Last modification time.
- `url` – The API URL of the project.
- `urls` – A map of related URLs.
- `execution_count`, `running_execution_count`, `queued_execution_count`, `enabled_endpoint_count`, `last_execution_ctime` – Various project stats.
