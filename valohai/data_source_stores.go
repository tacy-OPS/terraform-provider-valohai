package valohai

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceStore() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceStoreRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the store (UUID)",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the store",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the store (s3, swift, azure, google)",
			},
			"access_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Access mode (public, single_project, teams, owner_organization)",
			},
			"allow_read": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"allow_write": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"allow_uri_download": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"configuration": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"owner_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"project": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"paths": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"teams": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceStoreRead(d *schema.ResourceData, m interface{}) error {
	token := m.(map[string]interface{})["token"].(string)
	id := d.Get("id").(string)
	url := fmt.Sprintf("https://app.valohai.com/api/v0/stores/%s/", id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create GET request: %w", err)
	}
	req.Header.Set("Authorization", "Token "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute GET request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error %d", resp.StatusCode)
	}

	// Decode API response
	var result struct {
		ID               string                 `json:"id"`
		Name             string                 `json:"name"`
		Type             string                 `json:"type"`
		AccessMode       string                 `json:"access_mode"`
		AllowRead        bool                   `json:"allow_read"`
		AllowWrite       bool                   `json:"allow_write"`
		AllowURIDownload bool                   `json:"allow_uri_download"`
		Configuration    map[string]interface{} `json:"configuration"`
		Owner            int                    `json:"owner"`
		Project          interface{}            `json:"project"`
		Paths            map[string]interface{} `json:"paths"`
		Teams            []string               `json:"teams"`
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

	// Convert configuration values to string for Terraform state
	conf := map[string]string{}
	for k, v := range result.Configuration {
		conf[k] = fmt.Sprintf("%v", v)
	}
	d.Set("configuration", conf)
	d.Set("owner_id", result.Owner)

	// Project: handle string or object (API may return either)
	if result.Project != nil {
		switch v := result.Project.(type) {
		case string:
			d.Set("project", v)
		case map[string]interface{}:
			if id, ok := v["id"].(string); ok {
				d.Set("project", id)
			} else if id, ok := v["id"].(float64); ok {
				d.Set("project", fmt.Sprintf("%.0f", id))
			}
		}
	}

	// Convert paths values to string for Terraform state
	paths := map[string]string{}
	for k, v := range result.Paths {
		paths[k] = fmt.Sprintf("%v", v)
	}
	d.Set("paths", paths)
	d.Set("teams", result.Teams)
	d.Set("url", result.URL)
	return nil
}
