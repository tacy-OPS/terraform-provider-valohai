package valohai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTeam() *schema.Resource {
	return &schema.Resource{
		Create: resourceTeamCreate,
		Read:   resourceTeamRead,
		Update: resourceTeamUpdate,
		Delete: resourceTeamDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"organization": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// ResourceTeam returns the valohai team resource schema.
func ResourceTeam() *schema.Resource {
	return resourceTeam()
}

func resourceTeamCreate(d *schema.ResourceData, m interface{}) error {
	apiURL := "https://app.valohai.com/api/v0/teams/"
	authToken := m.(map[string]interface{})["token"].(string)

	orgID := d.Get("organization").(int)
	payload := map[string]interface{}{
		"name":         d.Get("name").(string),
		"organization": orgID,
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
		ID           string `json:"id"`
		Name         string `json:"name"`
		URL          string `json:"url"`
		Organization int    `json:"organization"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	d.SetId(result.ID)
	d.Set("name", result.Name)
	d.Set("organization", result.Organization)
	d.Set("url", result.URL)
	return nil
}

func resourceTeamRead(d *schema.ResourceData, m interface{}) error {
	authToken := m.(map[string]interface{})["token"].(string)
	id := d.Id() // UUID du projet Valohai
	url := fmt.Sprintf("https://app.valohai.com/api/v0/teams/%s/", id)

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
					return fmt.Errorf("API error %d (%s) - %s", resp.StatusCode, code, msg)
				}
			}
		}
		return fmt.Errorf("API error %d", resp.StatusCode)
	}

	// Decode response
	type organization struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	var result struct {
		ID           string       `json:"id"`
		Name         string       `json:"name"`
		URL          string       `json:"url"`
		Organization organization `json:"organization"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	d.SetId(result.ID)
	d.Set("name", result.Name)
	d.Set("organization", result.Organization.ID)
	d.Set("url", result.URL)

	return nil
}

func resourceTeamUpdate(d *schema.ResourceData, m interface{}) error {
	id := d.Id() // UUID du projet Valohai
	apiURL := fmt.Sprintf("https://app.valohai.com/api/v0/teams/%s/", id)
	authToken := m.(map[string]interface{})["token"].(string)

	payload := map[string]interface{}{
		"name": d.Get("name").(string),
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

func resourceTeamDelete(d *schema.ResourceData, m interface{}) error {
	authToken := m.(map[string]interface{})["token"].(string)
	id := d.Id() // UUID du projet Valohai
	url := fmt.Sprintf("https://app.valohai.com/api/v0/teams/%s/", id)

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
	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusOK {
		// 404 = already deleted, 204/200 = deleted OK
		return nil
	}
	return fmt.Errorf("unexpected status code: %d while deleting team %s", resp.StatusCode, id)
}
