package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func createValohaiTeam(apiToken, organization, name string) (string, error) {
	url := "https://app.valohai.com/api/v0/teams/"
	body := fmt.Sprintf(`{"name": "%s", "organization": "%s"}`, name, organization)
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Token "+apiToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		return "", fmt.Errorf("failed to create team: %s", resp.Status)
	}
	var result struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.ID, nil
}

func deleteValohaiTeam(apiToken, teamID string) error {
	url := fmt.Sprintf("https://app.valohai.com/api/v0/teams/%s/", teamID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Token "+apiToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 204 {
		return fmt.Errorf("failed to delete team: %s", resp.Status)
	}
	return nil
}

func testAccCheckValohaiStoreDestroy(s *terraform.State) error {
	return nil
}

func TestAccValohaiStoreOwnerOrg(t *testing.T) {
	// Skip test if required environment variables are not set
	if os.Getenv("VALOHAI_API_TOKEN") == "" || os.Getenv("VALOHAI_ORGANIZATION") == "" {
		t.Skip("VALOHAI_API_TOKEN or VALOHAI_ORGANIZATION is not set; skipping acceptance test.")
	}
	valohaiOrganization := getValohaiOrganization()
	name := fmt.Sprintf("tf-acc-test-ds-org-%s", uuid.New().String())
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
		CheckDestroy:      testAccCheckValohaiStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "valohai_store" "test" {
  name        = "%s"
  type        = "s3"
  access_mode = "owner_organization"
  allow_read  = true
  allow_write = true
  allow_uri_download = false
  configuration = {
	bucket = "example-bucket"
	region = "eu-west-1"
	access_key_id = "AKIAIOSFODNN7EXAMPLE1"
	secret_access_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	test_configuration = false
  }
  owner_id = %s
  paths = {
	"input"  = "data/input"
	"output" = "data/output"
  }
}
`, name, valohaiOrganization),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("valohai_store.test", "name", name),
					resource.TestCheckResourceAttr("valohai_store.test", "access_mode", "owner_organization"),
				),
			},
		},
	})
}

func TestAccValohaiStoreTeams(t *testing.T) {
	if os.Getenv("VALOHAI_API_TOKEN") == "" || os.Getenv("VALOHAI_ORGANIZATION") == "" {
		t.Skip("VALOHAI_API_TOKEN or VALOHAI_ORGANIZATION is not set; skipping acceptance test.")
	}
	apiToken := os.Getenv("VALOHAI_API_TOKEN")
	valohaiOrganization := getValohaiOrganization()
	teamName := fmt.Sprintf("tf-acc-team-%s", uuid.New().String())
	teamID, err := createValohaiTeam(apiToken, valohaiOrganization, teamName)
	if err != nil {
		t.Fatalf("Failed to create Valohai team: %v", err)
	}
	defer func() {
		if err := deleteValohaiTeam(apiToken, teamID); err != nil {
			t.Logf("Warning: failed to delete test team: %v", err)
		}
	}()
	name := fmt.Sprintf("tf-acc-test-ds-teams-%s", uuid.New().String())
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
		CheckDestroy:      testAccCheckValohaiStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "valohai_store" "test" {
  name        = "%s"
  type        = "s3"
  access_mode = "teams"
  allow_read  = true
  allow_write = true
  allow_uri_download = false
  configuration = {
	bucket = "example-bucket"
	region = "eu-west-1"
	access_key_id = "AKIAIOSFODNN7EXAMPLE2"
	secret_access_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	test_configuration = false
  }
  owner_id = %s
  teams = ["%s"]
  paths = {
	"input"  = "data/input"
	"output" = "data/output"
  }
}
`, name, valohaiOrganization, teamID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("valohai_store.test", "name", name),
					resource.TestCheckResourceAttr("valohai_store.test", "access_mode", "teams"),
				),
			},
		},
	})
}

func TestAccValohaiStoreSingleProject(t *testing.T) {
	// Skip test if required environment variables are not set
	if os.Getenv("VALOHAI_API_TOKEN") == "" || os.Getenv("VALOHAI_ORGANIZATION") == "" {
		t.Skip("VALOHAI_API_TOKEN or VALOHAI_ORGANIZATION is not set; skipping acceptance test.")
	}
	valohaiOrganization := getValohaiOrganization()
	projectID := os.Getenv("VALOHAI_PROJECT_ID")
	if projectID == "" {
		t.Skip("VALOHAI_PROJECT_ID is not set; skipping acceptance test.")
	}
	name := fmt.Sprintf("tf-acc-test-ds-proj-%s", uuid.New().String())
	defer deleteTestStateFiles()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if os.Getenv("VALOHAI_API_TOKEN") == "" {
				t.Fatal("VALOHAI_API_TOKEN must be set for acceptance tests")
			}
			if os.Getenv("VALOHAI_ORGANIZATION") == "" {
				t.Fatal("VALOHAI_ORGANIZATION must be set for acceptance tests")
			}
			if os.Getenv("VALOHAI_PROJECT_ID") == "" {
				t.Fatal("VALOHAI_PROJECT_ID must be set for acceptance tests")
			}
		},
		ProviderFactories: ProviderFactories,
		CheckDestroy:      testAccCheckValohaiStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "valohai_store" "test" {
  name        = "%s"
  type        = "s3"
  access_mode = "single_project"
  allow_read  = true
  allow_write = true
  allow_uri_download = false
  configuration = {
	bucket = "example-bucket"
	region = "eu-west-1"
	access_key_id = "AKIAIOSFODNN7EXAMPLE3"
	secret_access_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	test_configuration = false
  }
  owner_id = %s
  project = "%s"
  paths = {
	"input"  = "data/input"
	"output" = "data/output"
  }
}
`, name, valohaiOrganization, projectID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("valohai_store.test", "name", name),
					resource.TestCheckResourceAttr("valohai_store.test", "access_mode", "single_project"),
				),
			},
		},
	})
}
