terraform {
  required_providers {
    valohai = {
      source  = "tacy-ops/valohai"
      version = "0.1.0"
    }
  }
}

provider "valohai" {
  token = "example"
}

resource "valohai_project" "example" {
  name        = "example-terraform-project"
  description = "example terraform project"
  owner = "org-tacy-ops"
}

resource "valohai_team" "example" {
  name         = "example-terraform-team"
  organization = 0
}