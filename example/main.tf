terraform {
  required_providers {
    valohai = {
      source  = "tacy-ops/valohai"
      version = "0.1.0"
    }
  }
}

provider "valohai" {}

resource "valohai_project" "example" {
  name        = "example"
  description = "example terraform project"
  owner = "org-tacy-ops"
}
