# Data Source: valohai_datastore

The `valohai_datastore` data source allows you to retrieve information about an existing Valohai datastore by its ID.

## Example Usage

```hcl
data "valohai_datastore" "example" {
  id = "<datastore_id>"
}

output "datastore_name" {
  value = data.valohai_datastore.example.name
}
```

## Argument Reference

- `id` (String, Required): The UUID of the datastore to look up.

## Attributes Reference

All attributes are computed and match those of the `valohai_datastore` resource:

- `name` (String): Name of the datastore.
- `type` (String): Type of the datastore (`s3`, `swift`, `azure`, `google`).
- `access_mode` (String): Access mode (`public`, `single_project`, `teams`, `owner_organization`).
- `allow_read` (Bool): Whether read access is allowed.
- `allow_write` (Bool): Whether write access is allowed.
- `allow_uri_download` (Bool): Whether URI download is allowed.
- `configuration` (Map(String)): Provider-specific configuration.
- `owner_id` (Int): Organization ID.
- `project` (String): Associated project ID.
- `paths` (Map(String)): Named paths for the datastore.
- `teams` (List(String)): List of team IDs with access.
- `url` (String): Datastore URL.

## Import

This data source does not support import (use the resource for import).

## See Also
- [Valohai Documentation](https://valohai.com/docs/)
- [Terraform Provider Registry](https://registry.terraform.io/providers/tacy-ops/valohai/latest)
