package valohai

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTeam() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTeamRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"members": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"ctime": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"allow_project_administration": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_read_only": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"projects": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"ctime": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mtime": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"urls": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"execution_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"running_execution_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"queued_execution_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"enabled_endpoint_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"last_execution_ctime": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"organization": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceTeamRead(d *schema.ResourceData, m interface{}) error {
	token := m.(map[string]interface{})["token"].(string)
	id := d.Get("id").(string)
	url := fmt.Sprintf("https://app.valohai.com/api/v0/teams/%s/", id)

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
		return fmt.Errorf("valohai_team: team with id %s not found", id)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error %d", resp.StatusCode)
	}

	var result struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		URL          string `json:"url"`
		Organization struct {
			ID       int    `json:"id"`
			Username string `json:"username"`
		} `json:"organization"`
		Members []struct {
			User struct {
				ID       int    `json:"id"`
				Username string `json:"username"`
			} `json:"user"`
			Ctime                      string `json:"ctime"`
			AllowProjectAdministration bool   `json:"allow_project_administration"`
			IsReadOnly                 bool   `json:"is_read_only"`
		} `json:"members"`
		Projects []struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Owner       struct {
				ID       int    `json:"id"`
				Username string `json:"username"`
			} `json:"owner"`
			Ctime                 string            `json:"ctime"`
			Mtime                 string            `json:"mtime"`
			URL                   string            `json:"url"`
			Urls                  map[string]string `json:"urls"`
			ExecutionCount        int               `json:"execution_count"`
			RunningExecutionCount int               `json:"running_execution_count"`
			QueuedExecutionCount  int               `json:"queued_execution_count"`
			EnabledEndpointCount  int               `json:"enabled_endpoint_count"`
			LastExecutionCtime    string            `json:"last_execution_ctime"`
		} `json:"projects"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	d.SetId(result.ID)
	if err := d.Set("name", result.Name); err != nil {
		return fmt.Errorf("failed to set name: %w", err)
	}
	if err := d.Set("url", result.URL); err != nil {
		return fmt.Errorf("failed to set url: %w", err)
	}
	if err := d.Set("organization", map[string]interface{}{
		"id":       fmt.Sprintf("%v", result.Organization.ID),
		"username": result.Organization.Username,
	}); err != nil {
		return fmt.Errorf("failed to set organization: %w", err)
	}
	members := make([]map[string]interface{}, len(result.Members))
	for i, m := range result.Members {
		members[i] = map[string]interface{}{
			"user": map[string]interface{}{
				"id":       fmt.Sprintf("%v", m.User.ID),
				"username": m.User.Username,
			},
			"ctime":                        m.Ctime,
			"allow_project_administration": m.AllowProjectAdministration,
			"is_read_only":                 m.IsReadOnly,
		}
	}
	if err := d.Set("members", members); err != nil {
		return fmt.Errorf("failed to set members: %w", err)
	}
	projects := make([]map[string]interface{}, len(result.Projects))
	for i, p := range result.Projects {
		projects[i] = map[string]interface{}{
			"id":          p.ID,
			"name":        p.Name,
			"description": p.Description,
			"owner": map[string]interface{}{
				"id":       fmt.Sprintf("%v", p.Owner.ID),
				"username": p.Owner.Username,
			},
			"ctime":                   p.Ctime,
			"mtime":                   p.Mtime,
			"url":                     p.URL,
			"urls":                    p.Urls,
			"execution_count":         p.ExecutionCount,
			"running_execution_count": p.RunningExecutionCount,
			"queued_execution_count":  p.QueuedExecutionCount,
			"enabled_endpoint_count":  p.EnabledEndpointCount,
			"last_execution_ctime":    p.LastExecutionCtime,
		}
	}
	if err := d.Set("projects", projects); err != nil {
		return fmt.Errorf("failed to set projects: %w", err)
	}
	return nil
}
