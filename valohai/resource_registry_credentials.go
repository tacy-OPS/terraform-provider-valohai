package valohai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceRegistryCredentials() *schema.Resource {
	return &schema.Resource{
		Create: resourceRegistryCredentialsCreate,
		Read:   resourceRegistryCredentialsRead,
		Update: resourceRegistryCredentialsUpdate,
		Delete: resourceRegistryCredentialsDelete,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"docker", "aws_ecr", "aws_ecr_role"}, false,
				),
			},
			"image_pattern": {
				Type:     schema.TypeString,
				Required: true,
			},
			"owner": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"configuration": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

// ResourceRegistryCredentials returns the valohai registry credentials resource schema.
func ResourceRegistryCredentials() *schema.Resource {
	return resourceRegistryCredentials()
}

func resourceRegistryCredentialsCreate(d *schema.ResourceData, m interface{}) error {
	apiURL := "https://app.valohai.com/api/v0/registry-credentials/"
	authToken := m.(map[string]interface{})["token"].(string)

	payload := map[string]interface{}{
		"type":          d.Get("type").(string),
		"image_pattern": d.Get("image_pattern").(string),
	}
	// Optional fields
	if v, ok := d.GetOk("owner"); ok {
		payload["owner"] = v.(int)
	}
	if v, ok := d.GetOk("configuration"); ok {
		payload["configuration"] = v.(map[string]interface{})
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

func resourceRegistryCredentialsRead(d *schema.ResourceData, m interface{}) error {
	authToken := m.(map[string]interface{})["token"].(string)
	id := d.Id() // UUID du projet Valohai
	url := fmt.Sprintf("https://app.valohai.com/api/v0/registry-credentials/%s/", id)

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
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d while deleting project %s", resp.StatusCode, id)
	}
	return nil
}

func resourceRegistryCredentialsUpdate(d *schema.ResourceData, m interface{}) error {
	apiURL := "https://app.valohai.com/api/v0/registry-credentials/"
	authToken := m.(map[string]interface{})["token"].(string)

	payload := map[string]interface{}{
		"type":          d.Get("type").(string),
		"image_pattern": d.Get("image_pattern").(string),
	}
	// Optional fields
	if v, ok := d.GetOk("owner"); ok {
		payload["owner"] = v.(int)
	}
	if v, ok := d.GetOk("configuration"); ok {
		payload["configuration"] = v.(map[string]interface{})
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
	if resp.StatusCode != http.StatusCreated {
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

func resourceRegistryCredentialsDelete(d *schema.ResourceData, m interface{}) error {
	authToken := m.(map[string]interface{})["token"].(string)
	id := d.Id() // UUID du projet Valohai
	url := fmt.Sprintf("https://app.valohai.com/api/v0/registry-credentials/%s/", id)

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
