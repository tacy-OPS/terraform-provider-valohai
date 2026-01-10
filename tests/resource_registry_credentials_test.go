package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testAccCheckValohaiRegistryCredentialsDestroy(s *terraform.State) error {
	return nil
}

func TestAccValohaiRegistryCredentialsDocker_NoDrift(t *testing.T) {
	if os.Getenv("VALOHAI_API_TOKEN") == "" {
		t.Skip("VALOHAI_API_TOKEN is not set; skipping acceptance test.")
	}
	defer deleteTestStateFiles()

	name := fmt.Sprintf("tf-acc-regcred-docker-%s", uuid.New().String())

	cfg := fmt.Sprintf(`
resource "valohai_registry_credentials" "test" {
  type          = "docker"
  image_pattern = "docker.io/*"
  owner         = 9506

  configuration = {
    username = "test-%s"
    password = "test-%s"
  }
}
`, name, name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if os.Getenv("VALOHAI_API_TOKEN") == "" {
				t.Fatal("VALOHAI_API_TOKEN must be set for acceptance tests")
			}
		},
		ProviderFactories: ProviderFactories,
		CheckDestroy:      testAccCheckValohaiRegistryCredentialsDestroy,
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("valohai_registry_credentials.test", "id"),
					resource.TestCheckResourceAttr("valohai_registry_credentials.test", "type", "docker"),
					resource.TestCheckResourceAttr("valohai_registry_credentials.test", "image_pattern", "docker.io/*"),
				),
			},
			{
				Config:             cfg,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccValohaiRegistryCredentialsECR_NoDrift(t *testing.T) {
	if os.Getenv("VALOHAI_API_TOKEN") == "" {
		t.Skip("VALOHAI_API_TOKEN is not set; skipping acceptance test.")
	}
	defer deleteTestStateFiles()

	name := fmt.Sprintf("tf-acc-regcred-ecr-%s", uuid.New().String())
	cfg := fmt.Sprintf(`
resource "valohai_registry_credentials" "test" {
  type          = "aws-ecr"
  image_pattern = "123456789012.dkr.ecr.eu-west-1.amazonaws.com/*"
  owner         = 9506

  configuration = {
    access_key_id     = "AKIA%s"
    secret_access_key = "secret-%s"
    region            = "eu-west-1"
  }
}
`, name[:8], name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if os.Getenv("VALOHAI_API_TOKEN") == "" {
				t.Fatal("VALOHAI_API_TOKEN must be set for acceptance tests")
			}
		},
		ProviderFactories: ProviderFactories,
		CheckDestroy:      testAccCheckValohaiRegistryCredentialsDestroy,
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("valohai_registry_credentials.test", "id"),
					resource.TestCheckResourceAttr("valohai_registry_credentials.test", "type", "aws-ecr"),
				),
			},
			{
				Config:             cfg,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccValohaiRegistryCredentialsECRRole_NoDrift(t *testing.T) {
	if os.Getenv("VALOHAI_API_TOKEN") == "" {
		t.Skip("VALOHAI_API_TOKEN is not set; skipping acceptance test.")
	}
	defer deleteTestStateFiles()

	name := fmt.Sprintf("tf-acc-regcred-ecrrole-%s", uuid.New().String())

	cfg := fmt.Sprintf(`
resource "valohai_registry_credentials" "test" {
  type          = "aws-ecr-role"
  image_pattern = "123456789012.dkr.ecr.eu-west-1.amazonaws.com/*"
  owner         = 9506

  configuration = {
    role_name = "role-%s"
    region    = "eu-west-1"
  }
}
`, name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if os.Getenv("VALOHAI_API_TOKEN") == "" {
				t.Fatal("VALOHAI_API_TOKEN must be set for acceptance tests")
			}
		},
		ProviderFactories: ProviderFactories,
		CheckDestroy:      testAccCheckValohaiRegistryCredentialsDestroy,
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("valohai_registry_credentials.test", "id"),
					resource.TestCheckResourceAttr("valohai_registry_credentials.test", "type", "aws-ecr-role"),
				),
			},
			{
				Config:             cfg,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccValohaiRegistryCredentialsGCPCR_NoDrift(t *testing.T) {
	if os.Getenv("VALOHAI_API_TOKEN") == "" {
		t.Skip("VALOHAI_API_TOKEN is not set; skipping acceptance test.")
	}
	defer deleteTestStateFiles()

	sa := uuid.New().String()

	cfg := fmt.Sprintf(`
resource "valohai_registry_credentials" "test" {
  type          = "gcp-cr"
  image_pattern = "gcr.io/*"
  owner         = 9506

  configuration = {
    service_account_json = "%s"
  }
}
`, sa)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if os.Getenv("VALOHAI_API_TOKEN") == "" {
				t.Fatal("VALOHAI_API_TOKEN must be set for acceptance tests")
			}
		},
		ProviderFactories: ProviderFactories,
		CheckDestroy:      testAccCheckValohaiRegistryCredentialsDestroy,
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("valohai_registry_credentials.test", "id"),
					resource.TestCheckResourceAttr("valohai_registry_credentials.test", "type", "gcp-cr"),
				),
			},
			{
				Config:             cfg,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
