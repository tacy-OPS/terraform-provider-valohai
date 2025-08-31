# Terraform Provider: Valohai  

The **Valohai Terraform Provider** allows you to manage Valohai resources—such as teams, projects, and stores—directly from your Terraform configurations.  
It also includes multiple data sources to simplify querying and integrating existing Valohai entities into your infrastructure code.  

If you have suggestions or encounter issues, feel free to open an **Issue** on GitHub.  

---

## ✨ Features  

- Define and manage Valohai resources as Terraform-managed entities  
- Import and manage existing Valohai projects  
- Access metadata for projects, stores, and teams via data sources  
- Full support for local development and testing environments  

---

## 🔧 Provider Setup  

To begin using the provider, you'll need a valid **API token** from your Valohai account.  
You can generate this token in your **Account Settings** in Valohai.  

![Get your Valohai API token](https://help.valohai.com/hc/article_attachments/4419921059345/get_auth_token.gif)  

**Example configuration:**  

```hcl
provider "valohai" {
  token = "<your_valohai_token>"
}
```

If no token is defined in the provider block, the provider will automatically check for the environment variable `VALOHAI_API_TOKEN`:

```bash
export VALOHAI_API_TOKEN="<your_valohai_token>"
```

📦 Available Resources

Manage Valohai resources directly in your Terraform configuration:

- [valohai_project](ressources/valohai_project.md) - Create and manage Valohai projects
- [valohai_team](ressources/valohai_team.md) - Manage team configurations and memberships
- [valohai_store](ressources/valohai_store.md) - Define and manage Valohai stores

🔍 Data Sources

Query existing Valohai resources within your Terraform plans:

- [valohai_project](data-sources/valohai_project.md) - Access metadata for existing projects
- [valohai_team](data-sources/valohai_team.md) - Retrieve details about teams
- [valohai_store](data-sources/valohai_store.md) - Fetch information about existing stores