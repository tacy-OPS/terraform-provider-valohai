
# valohai_store Resource


The `valohai_store` resource allows you to manage stores in Valohai via Terraform.

## Example Usage

```hcl
resource "valohai_store" "example" {
  name        = "example-terraform-store"           # Unique name for the store
  type        = "s3"                                    # Type: s3, swift, azure, google
  access_mode = "teams"                                 # Access mode: public, single_project, teams, owner_organization
  allow_read  = true                                     # Allow read access
  allow_write = true                                     # Allow write access
  allow_uri_download = false                             # Allow URI download

  configuration = {
    bucket = "example-bucket"                            # S3 bucket name
    region = "eu-west-1"                                 # AWS region
    access_key_id = "AKIAIOSFODNN7EXAMPLE"              # Example only, use env vars or secrets manager in production
    secret_access_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEB" # Example only
    test_configuration = false
    # multipart_upload_iam_role = "ABC"
    # endpoint_url = "http://s3.test.com"
    # role_arn = "arn:aws:iam::123456789012:role/MyExampleRole"
    # kms_key_arn = "arn:aws:kms:us-west-2:123456789012:key/1234abcd-12ab-34cd-56ef-1234567890ab"
    # use_presigned_put_object = false
    # insecure = false
    # skip_upload_file_name_check = false
  }

  owner_id = 9506                                         # Organization ID
  project  = valohai_project.example.id                   # Associated project
  paths = {
    "input"  = "data/input"
    "output" = "data/output"
  }
  teams = [valohai_team.example.id]                       # List of team IDs
}
```

## Argument Reference

- `name` (String, Required): Name of the store (unique per organization).
- `type` (String, Required): Type of the store. One of `s3`, `swift`, `azure`, `google`.
- `access_mode` (String, Optional): Access mode. One of `public`, `single_project`, `teams`, `owner_organization`.
- `allow_read` (Bool, Optional): Allow read access. Default: `true`.
- `allow_write` (Bool, Optional): Allow write access. Default: `true`.
- `allow_uri_download` (Bool, Optional): Allow URI download. Default: `false`.
- `configuration` (Map(String), Optional): Provider-specific configuration. See above for S3 example.
- `owner_id` (Int, Optional): Organization ID.
- `project` (String, Optional): Associated project ID.
- `paths` (Map(String), Optional): Named paths for the store.
- `teams` (List(String), Optional): List of team IDs with access.

## Security Best Practices

- **Never commit real credentials in version control.**
- Use environment variables or a secrets manager to inject sensitive values.
- Restrict access to your state files if they contain secrets.

## Import

You can import an existing store by its ID:

```
terraform import valohai_store.example <store_id>
```

## See Also
- [Valohai Documentation](https://valohai.com/docs/)
- [Terraform Provider Registry](https://registry.terraform.io/providers/tacy-ops/valohai/latest)
