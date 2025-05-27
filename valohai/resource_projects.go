package valohai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_notifications": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

// ResourceProject returns the valohai project resource schema.
func ResourceProject() *schema.Resource {
	return resourceProject()
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	apiURL := "https://app.valohai.com/api/v0/projects/"
	authToken := m.(map[string]interface{})["token"].(string)

	payload := map[string]interface{}{
		"name":  d.Get("name").(string),
		"owner": d.Get("owner").(string),
	}

	// Optional fields
	if v, ok := d.GetOk("description"); ok {
		payload["description"] = v.(string)
	}
	if v, ok := d.GetOk("template_url"); ok {
		payload["template"] = v.(string)
	}
	if v, ok := d.GetOk("default_notifications"); ok {
		payload["default_notifications"] = v.(string) == "true"
	}

	// JSON encoding
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to encode payload: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+authToken)

	// Send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check HTTP status code
	if resp.StatusCode != http.StatusCreated {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("API error %d: %v", resp.StatusCode, errResp)
	}

	// Decode response
	var result struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	d.SetId(result.ID)

	return nil
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	// Implement logic to read a project.
	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	// Implement logic to update a project.
	return nil
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	// Implement logic to delete a project.
	return nil
}

// GetOptionalString returns the string value for a key if set, otherwise an empty string.
func GetOptionalString(d *schema.ResourceData, key string) string {
	if v, ok := d.GetOk(key); ok {
		return v.(string)
	}
	return ""
}
