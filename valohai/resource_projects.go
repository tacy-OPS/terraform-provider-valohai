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

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

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
		// Extraction du message et du code d'erreur s'ils existent
		if nonField, ok := errResp["non_field_errors"].([]interface{}); ok && len(nonField) > 0 {
			if first, ok := nonField[0].(map[string]interface{}); ok {
				msg, _ := first["message"].(string)
				code, _ := first["code"].(string)
				if code != "" || msg != "" {
					return fmt.Errorf("API error %d (%s) - %s", resp.StatusCode, code, msg)
				}
			}
		}
		return fmt.Errorf("API error %d", resp.StatusCode)
	}

	// Decode response
	var result struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	d.SetId(result.ID) // Stocke l'UUID Valohai comme ID de la ressource

	return nil
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	authToken := m.(map[string]interface{})["token"].(string)
	id := d.Id() // UUID du projet Valohai
	url := fmt.Sprintf("https://app.valohai.com/api/v0/projects/%s/", id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create GET request: %w", err)
	}
	req.Header.Set("Authorization", "Token "+authToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute GET request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		d.SetId("") // Not found, remove from state
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		if nonField, ok := errResp["non_field_errors"].([]interface{}); ok && len(nonField) > 0 {
			if first, ok := nonField[0].(map[string]interface{}); ok {
				msg, _ := first["message"].(string)
				code, _ := first["code"].(string)
				if code != "" || msg != "" {
					return fmt.Errorf("API error %d: %s %s", resp.StatusCode, code, msg)
				}
			}
		}
		return fmt.Errorf("API error %d", resp.StatusCode)
	}

	// Decode response
	var result struct {
		ID                   string      `json:"id"`
		Name                 string      `json:"name"`
		Owner                interface{} `json:"owner"`
		Description          string      `json:"description"`
		Template             string      `json:"template"`
		DefaultNotifications bool        `json:"default_notifications"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	d.SetId(result.ID)
	d.Set("name", result.Name)
	// Gestion owner string ou map
	if ownerStr, ok := result.Owner.(string); ok {
		d.Set("owner", ownerStr)
	} else if ownerMap, ok := result.Owner.(map[string]interface{}); ok {
		if slug, ok := ownerMap["slug"].(string); ok {
			d.Set("owner", slug)
		}
	}
	d.Set("description", result.Description)
	d.Set("template_url", result.Template)
	// Ne set que si true, sinon laisse null
	if result.DefaultNotifications {
		d.Set("default_notifications", "true")
	}
	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	id := d.Id() // UUID du projet Valohai
	apiURL := fmt.Sprintf("https://app.valohai.com/api/v0/projects/%s/", id)
	authToken := m.(map[string]interface{})["token"].(string)

	payload := map[string]interface{}{
		"name": d.Get("name").(string),
	}
	// Optional fields
	if v, ok := d.GetOk("description"); ok {
		payload["description"] = v.(string)
	}

	// JSON encoding
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to encode payload: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("PUT", apiURL, bytes.NewBuffer(body))
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
	if resp.StatusCode != http.StatusOK {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		// Extraction du message et du code d'erreur s'ils existent
		if nonField, ok := errResp["non_field_errors"].([]interface{}); ok && len(nonField) > 0 {
			if first, ok := nonField[0].(map[string]interface{}); ok {
				msg, _ := first["message"].(string)
				code, _ := first["code"].(string)
				if code != "" || msg != "" {
					return fmt.Errorf("API error %d: %s %s", resp.StatusCode, code, msg)
				}
			}
		}
		return fmt.Errorf("API error %d", resp.StatusCode)
	}

	// Decode response
	var result struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	d.SetId(result.ID) // Stocke l'UUID Valohai comme ID de la ressource
	return nil
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	authToken := m.(map[string]interface{})["token"].(string)
	id := d.Id() // UUID du projet Valohai
	url := fmt.Sprintf("https://app.valohai.com/api/v0/projects/%s/", id)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create DELETE request: %w", err)
	}
	req.Header.Set("Authorization", "Token "+authToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute DELETE request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d while deleting project %s", resp.StatusCode, id)
	}
	return nil
}

// GetOptionalString returns the string value for a key if set, otherwise an empty string.
func GetOptionalString(d *schema.ResourceData, key string) string {
	if v, ok := d.GetOk(key); ok {
		return v.(string)
	}
	return ""
}
