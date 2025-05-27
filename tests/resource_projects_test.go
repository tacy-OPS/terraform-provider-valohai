package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var _ = ProviderFactories

// getValohaiOwner retrieves the Valohai owner from the environment variable.
func getValohaiOwner() string {
	owner := os.Getenv("VALOHAI_OWNER")
	if owner == "" {
		panic("The VALOHAI_OWNER environment variable must be set for acceptance tests.")
	}
	return owner
}

func TestAccValohaiProjectBasic(t *testing.T) {
	if os.Getenv("VALOHAI_API_TOKEN") == "" {
		t.Skip("VALOHAI_API_TOKEN is not set; skipping acceptance test.")
	}
	valohaiOwner := getValohaiOwner()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if os.Getenv("VALOHAI_API_TOKEN") == "" {
				t.Fatal("VALOHAI_API_TOKEN must be set for acceptance tests")
			}
			if os.Getenv("VALOHAI_OWNER") == "" {
				t.Fatal("VALOHAI_OWNER must be set for acceptance tests")
			}
		},
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "valohai_project" "test" {
  name  = "tf-acc-test-project"
  owner = "` + valohaiOwner + `"
  description = "Project created by acceptance test"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("valohai_project.test", "name", "tf-acc-test-project"),
					resource.TestCheckResourceAttr("valohai_project.test", "description", "Project created by acceptance test"),
				),
			},
		},
	})
}

func TestAccValohaiProjectUpdate(t *testing.T) {
	if os.Getenv("VALOHAI_API_TOKEN") == "" {
		t.Skip("VALOHAI_API_TOKEN is not set; skipping acceptance test.")
	}
	valohaiOwner := getValohaiOwner()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if os.Getenv("VALOHAI_API_TOKEN") == "" {
				t.Fatal("VALOHAI_API_TOKEN must be set for acceptance tests")
			}
			if os.Getenv("VALOHAI_OWNER") == "" {
				t.Fatal("VALOHAI_OWNER must be set for acceptance tests")
			}
		},
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "valohai_project" "test" {
  name  = "tf-acc-test-project-update"
  owner = "` + valohaiOwner + `"
  description = "Initial description"
}
`,
				Check: resource.TestCheckResourceAttr("valohai_project.test", "description", "Initial description"),
			},
			{
				Config: `
resource "valohai_project" "test" {
  name  = "tf-acc-test-project-new-update"
  owner = "` + valohaiOwner + `"
  description = "Updated description"
}
`,
				Check: resource.TestCheckResourceAttr("valohai_project.test", "description", "Updated description"),
			},
		},
	})
}

func TestAccValohaiProjectDelete(t *testing.T) {
	if os.Getenv("VALOHAI_API_TOKEN") == "" {
		t.Skip("VALOHAI_API_TOKEN is not set; skipping acceptance test.")
	}
	valohaiOwner := getValohaiOwner()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if os.Getenv("VALOHAI_API_TOKEN") == "" {
				t.Fatal("VALOHAI_API_TOKEN must be set for acceptance tests")
			}
			if os.Getenv("VALOHAI_OWNER") == "" {
				t.Fatal("VALOHAI_OWNER must be set for acceptance tests")
			}
		},
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "valohai_project" "test" {
  name  = "tf-acc-test-project-delete"
  owner = "` + valohaiOwner + `"
  description = "Project to delete"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("valohai_project.test", "name", "tf-acc-test-project-delete"),
				),
			},
			{
				Config:       "",
				ImportState:  false,
				RefreshState: true,
			},
		},
	})
}
