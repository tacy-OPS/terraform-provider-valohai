# Terraform Provider Valohai

A Terraform provider to manage Valohai projects and resources, including project creation, import, and data access. This provider allows you to automate and manage your Valohai resources directly from your Terraform configuration.

- Manage Valohai projects as Terraform resources
- Import existing Valohai projects
- Access project details via data sources
- Full support for local development and testing

Welcome to the Valohai Terraform provider documentation.

This provider allows you to manage Valohai resources using Terraform.

## Provider Configuration

To use the Valohai provider, you must supply a valid API token. You can create this token in your Valohai account (see the official Valohai documentation for more details).

![Get your Valohai API token](https://help.valohai.com/hc/article_attachments/4419921059345/get_auth_token.gif)

Example configuration:


```hcl
provider "valohai" {
  token = "<your_valohai_token>"
}
```

## Resources

- [valohai_project](resources/valohai_project.md): Manage Valohai projects
- [valohai_team](resources/valohai_team.md): Manage Valohai teams

## Data Sources

- [valohai_project](data-sources/valohai_project.md): Retrieve information about an existing Valohai project
- [valohai_team](data-sources/valohai_team.md): Retrieve information about an existing Valohai team
