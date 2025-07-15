# Terraform Provider Valohai

The Valohai Terraform Provider enables you to manage and automate your Valohai resources—such as projects, stores, and teams—directly through your Terraform configuration.
This integration streamlines infrastructure-as-code workflows for machine learning operations using Valohai.

## ✨ Features
- Define and manage Valohai resources as Terraform-managed entities
- Import and work with existing Valohai projects
- Access project, store, and team metadata via Terraform data sources
- Fully supports local development and testing environments
Welcome to the official documentation for the Valohai Terraform provider

Welcome to the Valohai Terraform provider documentation.


## 🔧 Provider setup

To begin using the provider, you'll need a valid API token from your Valohai account.
Visit your account settings in Valohai to generate an authentication token.

![Get your Valohai API token](https://help.valohai.com/hc/article_attachments/4419921059345/get_auth_token.gif)

Example configuration:


```hcl
provider "valohai" {
  token = "<your_valohai_token>"
}
```

## 📦 Available Resources
Manage Valohai infrastructure directly from your Terraform files:
- [valohai_project](resources/valohai_project.md) – Create and update Valohai projects
- [valohai_team](resources/valohai_team.md) – Manage team configurations and membership
- [valohai_team](resources/valohai_team.md) – Define and control Valohai store

## 🔍 Data Sources
Reference and retrieve existing resources within your Terraform plans:
- [valohai_project](data-sources/valohai_project.md) – Access metadata from existing projects
- [valohai_team](data-sources/valohai_team.md)– Query details about your teams
- [valohai_store](data-sources/valohai_store.md) – Fetch information about existing stores
