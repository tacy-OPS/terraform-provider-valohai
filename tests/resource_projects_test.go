package tests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var _ = ProviderFactories

func getValohaiOwner() string {
	owner := os.Getenv("VALOHAI_OWNER")
	if owner == "" {
		panic("The VALOHAI_OWNER environment variable must be set for acceptance tests.")
	}
	return owner
}

func testAccCheckValohaiProjectDestroy(s *terraform.State) error {
	return nil
}

func TestAccValohaiProjectBasic(t *testing.T) {
	if os.Getenv("VALOHAI_API_TOKEN") == "" {
		t.Skip("VALOHAI_API_TOKEN is not set; skipping acceptance test.")
	}
	valohaiOwner := getValohaiOwner()
	name := uniqueName("tf-acc-test-project")
	defer deleteTestStateFiles()
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
		CheckDestroy:      testAccCheckValohaiProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: `
resource "valohai_project" "test" {
  name  = "` + name + `"
  owner = "` + valohaiOwner + `"
  description = "Project created by acceptance test"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("valohai_project.test", "name", name),
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
	name := uniqueName("tf-acc-test-project-update")
	defer deleteTestStateFiles()
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
		CheckDestroy:      testAccCheckValohaiProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: `
resource "valohai_project" "test" {
  name  = "` + name + `"
  owner = "` + valohaiOwner + `"
  description = "Initial description"
}
`,
				Check: resource.TestCheckResourceAttr("valohai_project.test", "description", "Initial description"),
			},
			{
				Config: `
resource "valohai_project" "test" {
  name  = "` + name + `"
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
	name := uniqueName("tf-acc-test-project-delete")
	defer deleteTestStateFiles()
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
		CheckDestroy:      testAccCheckValohaiProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: `
resource "valohai_project" "test" {
  name  = "` + name + `"
  owner = "` + valohaiOwner + `"
  description = "Project to delete"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("valohai_project.test", "name", name),
				),
			},
			{
				Config: "# Empty config to trigger resource deletion\n",
			},
		},
	})
}
