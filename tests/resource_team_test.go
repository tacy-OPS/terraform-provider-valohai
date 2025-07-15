package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testAccCheckValohaiTeamDestroy(s *terraform.State) error {
	return nil
}

func TestAccValohaiTeamBasic(t *testing.T) {
	if os.Getenv("VALOHAI_API_TOKEN") == "" || os.Getenv("VALOHAI_ORGANIZATION") == "" {
		t.Skip("VALOHAI_API_TOKEN or VALOHAI_ORGANIZATION is not set; skipping acceptance test.")
	}
	valohaiOrganization := getValohaiOrganization()
	name := fmt.Sprintf("tf-acc-test-team-%s", uuid.New().String())
	defer deleteTestStateFiles()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if os.Getenv("VALOHAI_API_TOKEN") == "" {
				t.Fatal("VALOHAI_API_TOKEN must be set for acceptance tests")
			}
			if os.Getenv("VALOHAI_ORGANIZATION") == "" {
				t.Fatal("VALOHAI_ORGANIZATION must be set for acceptance tests")
			}
		},
		ProviderFactories: ProviderFactories,
		CheckDestroy:      testAccCheckValohaiTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: `
resource "valohai_team" "test" {
  name         = "` + name + `"
  organization = "` + valohaiOrganization + `"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("valohai_team.test", "name", name),
				),
			},
		},
	})
}

func TestAccValohaiTeamUpdate(t *testing.T) {
	if os.Getenv("VALOHAI_API_TOKEN") == "" || os.Getenv("VALOHAI_ORGANIZATION") == "" {
		t.Skip("VALOHAI_API_TOKEN or VALOHAI_ORGANIZATION is not set; skipping acceptance test.")
	}
	valohaiOrganization := getValohaiOrganization()
	name := fmt.Sprintf("tf-acc-test-team-update-%s", uuid.New().String())
	defer deleteTestStateFiles()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if os.Getenv("VALOHAI_API_TOKEN") == "" {
				t.Fatal("VALOHAI_API_TOKEN must be set for acceptance tests")
			}
			if os.Getenv("VALOHAI_ORGANIZATION") == "" {
				t.Fatal("VALOHAI_ORGANIZATION must be set for acceptance tests")
			}
		},
		ProviderFactories: ProviderFactories,
		CheckDestroy:      testAccCheckValohaiTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: `
resource "valohai_team" "test" {
  name  = "` + name + `"
  organization = "` + valohaiOrganization + `"
}
`,
				Check: resource.TestCheckResourceAttr("valohai_team.test", "name", name),
			},
			{
				Config: `
resource "valohai_team" "test" {
  name  = "` + name + `"
  organization = "` + valohaiOrganization + `"
}
`,
				Check: resource.TestCheckResourceAttr("valohai_team.test", "name", name),
			},
		},
	})
}

func TestAccValohaiTeamDelete(t *testing.T) {
	if os.Getenv("VALOHAI_API_TOKEN") == "" || os.Getenv("VALOHAI_ORGANIZATION") == "" {
		t.Skip("VALOHAI_API_TOKEN or VALOHAI_ORGANIZATION is not set; skipping acceptance test.")
	}
	valohaiOrganization := getValohaiOrganization()
	name := uniqueName("tf-acc-test-team-delete")
	defer deleteTestStateFiles()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if os.Getenv("VALOHAI_API_TOKEN") == "" {
				t.Fatal("VALOHAI_API_TOKEN must be set for acceptance tests")
			}
			if os.Getenv("VALOHAI_ORGANIZATION") == "" {
				t.Fatal("VALOHAI_ORGANIZATION must be set for acceptance tests")
			}
		},
		ProviderFactories: ProviderFactories,
		CheckDestroy:      testAccCheckValohaiTeamDestroy,
		Steps: []resource.TestStep{
			{
				Config: `
resource "valohai_team" "test" {
  name         = "` + name + `"
  organization = "` + valohaiOrganization + `"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("valohai_team.test", "name", name),
				),
			},
			{
				Config: "# Empty config to trigger resource deletion\n",
			},
		},
	})
}
