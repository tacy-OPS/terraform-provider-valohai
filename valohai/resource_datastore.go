package valohai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDatastore() *schema.Resource {
	return &schema.Resource{
		Create: resourceDatastoreCreate,
		Read:   resourceDatastoreRead,
		Update: resourceDatastoreUpdate,
		Delete: resourceDatastoreDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the datastore (max 64 characters)",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if len(v) > 64 {
						errs = append(errs, fmt.Errorf("%q cannot be longer than 64 characters", key))
					}
					return
				},
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the datastore (s3, swift, azure, google)",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					allowed := map[string]bool{
						"s3":     true,
						"swift":  true,
						"azure":  true,
						"google": true,
					}
					if !allowed[v] {
						errs = append(errs, fmt.Errorf("%q must be one of [s3, swift, azure, google], got %q", key, v))
					}
					return
				},
			},
			"access_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Access mode (public, single_project, teams, owner_organization)",
			},
			"allow_read": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"allow_write": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"allow_uri_download": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"configuration": {
				Type:     schema.TypeMap,
				Optional: true,
				Default:  map[string]interface{}{},
			},
			"owner": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"paths": {
				Type:     schema.TypeMap,
				Optional: true,
				Default:  map[string]interface{}{},
			},
			"teams": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

// ResourceDatastore returns the valohai datastore resource schema.
func ResourceDatastore() *schema.Resource {
	return resourceDatastore()
}

func resourceDatastoreCreate(d *schema.ResourceData, m interface{}) error {
	apiURL := "https://app.valohai.com/api/v0/stores/"
	authToken := m.(map[string]interface{})["token"].(string)

	// Build payload based on defined schema, only including set fields
	payload := map[string]interface{}{
		"name": d.Get("name").(string),
		"type": d.Get("type").(string),
	}

	if v, ok := d.GetOk("access_mode"); ok {
		payload["access_mode"] = v.(string)
	}
	if v, ok := d.GetOk("allow_read"); ok {
		payload["allow_read"] = v.(bool)
	}
	if v, ok := d.GetOk("allow_write"); ok {
		payload["allow_write"] = v.(bool)
	}
	if v, ok := d.GetOk("allow_uri_download"); ok {
		payload["allow_uri_download"] = v.(bool)
	}
	if v, ok := d.GetOk("configuration"); ok {
		payload["configuration"] = v
	}
	if v, ok := d.GetOk("owner"); ok {
		payload["owner"] = v.(string)
	}
	if v, ok := d.GetOk("project"); ok {
		payload["project"] = v.(string)
	}
	if v, ok := d.GetOk("paths"); ok {
		payload["paths"] = v
	}
	if v, ok := d.GetOk("teams"); ok {
		teams := []string{}
		for _, t := range v.([]interface{}) {
			teams = append(teams, t.(string))
		}
		payload["teams"] = teams
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
		ID               string                 `json:"id"`
		Name             string                 `json:"name"`
		Type             string                 `json:"type"`
		AccessMode       string                 `json:"access_mode"`
		AllowRead        bool                   `json:"allow_read"`
		AllowWrite       bool                   `json:"allow_write"`
		AllowURIDownload bool                   `json:"allow_uri_download"`
		Configuration    map[string]interface{} `json:"configuration"`
		Owner            interface{}            `json:"owner"`
		Paths            map[string]interface{} `json:"paths"`
		Teams            []string               `json:"teams"`
		Project          string                 `json:"project"`
		URL              string                 `json:"url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	d.SetId(result.ID)
	d.Set("name", result.Name)
	d.Set("type", result.Type)
	d.Set("access_mode", result.AccessMode)
	d.Set("allow_read", result.AllowRead)
	d.Set("allow_write", result.AllowWrite)
	d.Set("allow_uri_download", result.AllowURIDownload)
	d.Set("configuration", result.Configuration)
	d.Set("owner", result.Owner)
	d.Set("paths", result.Paths)
	d.Set("teams", result.Teams)
	d.Set("project", result.Project)
	d.Set("url", result.URL)
	return nil
}

func resourceDatastoreRead(d *schema.ResourceData, m interface{}) error {
	authToken := m.(map[string]interface{})["token"].(string)
	id := d.Id() // UUID datastore Valohai
	url := fmt.Sprintf("https://app.valohai.com/api/v0/stores/%s/", id)

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
	var result struct {
		ID               string                 `json:"id"`
		Ctime            string                 `json:"ctime"`
		Mtime            string                 `json:"mtime"`
		Type             string                 `json:"type"`
		Deleted          bool                   `json:"deleted"`
		Name             string                 `json:"name"`
		Owner            interface{}            `json:"owner"`
		Project          map[string]interface{} `json:"project"`
		Configuration    map[string]interface{} `json:"configuration"`
		Paths            map[string]interface{} `json:"paths"`
		Teams            []string               `json:"teams"`
		UriPrefix        interface{}            `json:"uri_prefix"`
		AllowRead        bool                   `json:"allow_read"`
		AllowWrite       bool                   `json:"allow_write"`
		AllowAdopt       bool                   `json:"allow_adopt"`
		AccessMode       string                 `json:"access_mode"`
		AllowURIDownload bool                   `json:"allow_uri_download"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	d.SetId(result.ID)
	d.Set("ctime", result.Ctime)
	d.Set("mtime", result.Mtime)
	d.Set("type", result.Type)
	d.Set("deleted", result.Deleted)
	d.Set("name", result.Name)
	d.Set("owner", result.Owner)
	d.Set("project", result.Project)
	d.Set("configuration", result.Configuration)
	d.Set("paths", result.Paths)
	d.Set("teams", result.Teams)
	d.Set("uri_prefix", result.UriPrefix)
	d.Set("allow_read", result.AllowRead)
	d.Set("allow_write", result.AllowWrite)
	d.Set("allow_adopt", result.AllowAdopt)
	d.Set("access_mode", result.AccessMode)
	d.Set("allow_uri_download", result.AllowURIDownload)

	return nil
}

func resourceDatastoreUpdate(d *schema.ResourceData, m interface{}) error {
	id := d.Id() // UUID datastore Valohai
	apiURL := fmt.Sprintf("https://app.valohai.com/api/v0/stores/%s/", id)
	authToken := m.(map[string]interface{})["token"].(string)

	payload := map[string]interface{}{}

	if v, ok := d.GetOk("name"); ok {
		payload["name"] = v.(string)
	}
	if v, ok := d.GetOk("access_mode"); ok {
		payload["access_mode"] = v.(string)
	}
	if v, ok := d.GetOk("allow_read"); ok {
		payload["allow_read"] = v.(bool)
	}
	if v, ok := d.GetOk("allow_write"); ok {
		payload["allow_write"] = v.(bool)
	}
	if v, ok := d.GetOk("allow_uri_download"); ok {
		payload["allow_uri_download"] = v.(bool)
	}
	if v, ok := d.GetOk("configuration"); ok {
		payload["configuration"] = v
	}
	if v, ok := d.GetOk("owner"); ok {
		payload["owner"] = v.(string)
	}
	if v, ok := d.GetOk("paths"); ok {
		payload["paths"] = v
	}
	if v, ok := d.GetOk("teams"); ok {
		teams := []string{}
		for _, t := range v.([]interface{}) {
			teams = append(teams, t.(string))
		}
		payload["teams"] = teams
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
		AccessMode       string                 `json:"access_mode"`
		AllowRead        bool                   `json:"allow_read"`
		AllowURIDownload bool                   `json:"allow_uri_download"`
		AllowWrite       bool                   `json:"allow_write"`
		Configuration    map[string]interface{} `json:"configuration"`
		Name             string                 `json:"name"`
		Owner            interface{}            `json:"owner"`
		Paths            map[string]interface{} `json:"paths"`
		Teams            []string               `json:"teams"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode update response: %w", err)
	}
	d.Set("access_mode", result.AccessMode)
	d.Set("allow_read", result.AllowRead)
	d.Set("allow_uri_download", result.AllowURIDownload)
	d.Set("allow_write", result.AllowWrite)
	d.Set("configuration", result.Configuration)
	d.Set("name", result.Name)
	d.Set("owner", result.Owner)
	d.Set("paths", result.Paths)
	d.Set("teams", result.Teams)
	return nil
}

func resourceDatastoreDelete(d *schema.ResourceData, m interface{}) error {
	authToken := m.(map[string]interface{})["token"].(string)
	id := d.Id() // UUID Datastore Valoha
	url := fmt.Sprintf("https://app.valohai.com/api/v0/stores/%s/", id)

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
