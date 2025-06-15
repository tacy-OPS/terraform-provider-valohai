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

resource "valohai_datastore" "example" {
  name        = "example-terraform-datastore"
  type       = "s3"
  access_mode = "public"
  allow_read  = true
  allow_write = true
  allow_uri_download = false
  configuration = {
    bucket = "example-bucket"
    region = "eu-west-1"
  }
  owner     = "org-tacy-ops"
  project   = valohai_project.example.id
  paths     = {
    "input"  = "data/input"
    "output" = "data/output"
  }
  teams     = [valohai_team.example.id]
}