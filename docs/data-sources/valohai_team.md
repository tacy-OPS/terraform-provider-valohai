# Data Source: valohai_team

Use this data source to retrieve information about an existing Valohai team by its ID.

## Example Usage

```hcl
provider "valohai" {
  token = "${var.valohai_token}"
}

data "valohai_team" "example" {
  id = "01974bcf-52a2-6e8d-b454-900e415908d4"
}

output "team_name" {
  value = data.valohai_team.example.name
}
```

## Argument Reference

- `id` (String, Required): The ID of the team to look up.

## Attributes Reference

The following attributes are exported:

- `id` – The team ID.
- `name` – The name of the team.
- `url` – The API URL of the team.
- `organization` – A map with the organization's `id` and `username`.
- `members` – A list of team members, each with:
  - `user` (map: `id`, `username`)
  - `ctime`
  - `allow_project_administration`
  - `is_read_only`
- `projects` – A list of projects, each with:
  - `id`, `name`, `description`, `owner` (map: `id`, `username`), `ctime`, `mtime`, `url`, `urls`, `execution_count`, `running_execution_count`, `queued_execution_count`, `enabled_endpoint_count`, `last_execution_ctime`
