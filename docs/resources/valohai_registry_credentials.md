# Resource: valohai_registry_credentials

Manages container registry credentials in Valohai.

This resource allows you to configure authentication for different container
registries (Docker Hub, AWS ECR, AWS ECR via IAM role, Google Container Registry)
used by Valohai.

The provider automatically injects required internal defaults (such as schema
versions) and hides implementation details from the user.

---

## Example Usage

### Docker Registry

```hcl
resource "valohai_registry_credentials" "docker" {
  type          = "docker"
  image_pattern = "docker.io/*"
  owner         = 9506

  configuration = {
    username = "myuser"
    password = "mypassword"
  }
}
```

### AWS ECR (Access Key)

```hcl
resource "valohai_registry_credentials" "ecr" {
  type          = "aws-ecr"
  image_pattern = "123456789012.dkr.ecr.eu-west-1.amazonaws.com/*"
  owner         = 9506

  configuration = {
    access_key_id     = "AKIA..."
    secret_access_key = "secret"
    region            = "eu-west-1"
  }
}
```

### AWS ECR (IAM Role)

```hcl
resource "valohai_registry_credentials" "ecr_role" {
  type          = "aws-ecr-role"
  image_pattern = "123456789012.dkr.ecr.eu-west-1.amazonaws.com/*"
  owner         = 9506

  configuration = {
    role_name = "valohai-ecr-role"
    region    = "eu-west-1"
  }
}
```

### Google Container Registry (GCP)

```hcl

resource "valohai_registry_credentials" "gcp" {
  type          = "gcp-cr"
  image_pattern = "gcr.io/*"
  owner         = 9506

  configuration = {
    service_account_json = file("service-account.json")
  }
}

```


