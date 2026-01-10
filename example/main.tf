terraform {
  required_providers {
    valohai = {
      source  = "tacy-ops/valohai"
      version = "0.1.0"
    }
  }
}

# provider "valohai" {
#   token = "example"
# }

resource "valohai_project" "example" {
  name        = "example-terraform-project"
  description = "example terraform project"
  owner = "org-tacy-ops"
}

resource "valohai_team" "example" {
  name         = "example-terraform-team"
  organization = 0 # Using 0 for the default organization
}

resource "valohai_store" "example" {
  name        = "example-terraform-store"
  type       = "s3"
  access_mode = "single_project"
  allow_read  = true
  allow_write = true
  allow_uri_download = false
  configuration = {
    bucket = "example-bucket"
    region = "eu-west-1"
    # Fake credentials for example purposes
    # In a real scenario, you would use environment variables or a secure vault to manage credentials
    access_key_id = "AKIAIOSFODNN7EXAMPLE"
    secret_access_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEB"
    test_configuration = false
    # multipart_upload_iam_role = "ABC"
    # endpoint_url = "http://s3.test.com"
    # role_arn = "arn:aws:iam::123456789012:role/MyExampleRole"
    # kms_key_arn = "arn:aws:kms:us-west-2:123456789012:key/1234abcd-12ab-34cd-56ef-1234567890ab"
    # use_presigned_put_object = false
    # insecure = false
    # skip_upload_file_name_check = false
  }
  owner_id = 9506
  project   = valohai_project.example.id
  paths     = {
    "input"  = "data/input"
    "output" = "data/output"
  }
  # teams     = [valohai_team.example.id]
}

resource "valohai_registry_credentials" "docker" {
  type = "docker"
  image_pattern = "docker.io/*"
  configuration = {
    password = "test"
    username = "test"
  }
}

resource "valohai_registry_credentials" "ecr" {
  type = "aws-ecr"
  image_pattern = "ecr.io/*"
  configuration = {
    access_key_id = "AKIAIOSFODNN7EXAMPLE"
    secret_access_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEB"
    region = "eu-west-1"
  }
}

resource "valohai_registry_credentials" "ecr_role" {
  type = "aws-ecr-role"
  image_pattern = "ecr-with-role.io/*"
  configuration = {
    region = "eu-west-1"
    role_name = "arn:aws:iam::123456789012:role/MyExampleRole"
  }
}

resource "valohai_registry_credentials" "gcp_cr" {
  type = "gcp-cr"
  image_pattern = "gcp.io/*"
  configuration = {
    service_account_json = "heOhRHt3MyQYfXOD054uqnW7XSlzwBXRjsU45bBD"
  }
}