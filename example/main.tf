terraform {
  required_providers {
    valohai = {
      source  = "github.com/tacy-OPS/valohai"
      version = "0.1.0"
    }
  }
}

provider "valohai" {}

resource "valohai_project" "example" {
  name        = "example"
  description = "example terraform project"
  owner = "example"
}
