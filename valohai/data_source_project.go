package valohai

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceProject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceProjectRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
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
			"environment_variables": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"execution_summary": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"repository": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project": {Type: schema.TypeString, Computed: true},
						"name":    {Type: schema.TypeString, Computed: true},
						"color":   {Type: schema.TypeString, Computed: true},
					},
				},
			},
			"upload_store_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"read_only": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"yaml_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceProjectRead(d *schema.ResourceData, m interface{}) error {
	token := m.(map[string]interface{})["token"].(string)
	id := d.Get("id").(string)
	url := fmt.Sprintf("https://app.valohai.com/api/v0/projects/%s/", id)

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

	type tag struct {
		Project string `json:"project"`
		Name    string `json:"name"`
		Color   string `json:"color"`
	}
	type owner struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
	}
	type repository struct {
		ID  int    `json:"id"`
		URL string `json:"url"`
		Ref string `json:"ref"`
	}
	var result struct {
		ID                    string                 `json:"id"`
		Name                  string                 `json:"name"`
		Description           string                 `json:"description"`
		Owner                 owner                  `json:"owner"`
		Ctime                 string                 `json:"ctime"`
		Mtime                 string                 `json:"mtime"`
		Url                   string                 `json:"url"`
		Urls                  map[string]string      `json:"urls"`
		ExecutionCount        int                    `json:"execution_count"`
		RunningExecutionCount int                    `json:"running_execution_count"`
		QueuedExecutionCount  int                    `json:"queued_execution_count"`
		EnabledEndpointCount  int                    `json:"enabled_endpoint_count"`
		LastExecutionCtime    string                 `json:"last_execution_ctime"`
		EnvironmentVariables  map[string]interface{} `json:"environment_variables"`
		ExecutionSummary      map[string]int         `json:"execution_summary"`
		Repository            repository             `json:"repository"`
		Tags                  []tag                  `json:"tags"`
		UploadStoreID         string                 `json:"upload_store_id"`
		ReadOnly              bool                   `json:"read_only"`
		YamlPath              string                 `json:"yaml_path"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	d.SetId(result.ID)
	// Set all attributes, converting IDs to string as required by Terraform schema
	if err := d.Set("owner", map[string]interface{}{
		"id":       fmt.Sprintf("%v", result.Owner.ID),
		"username": result.Owner.Username,
	}); err != nil {
		return fmt.Errorf("failed to set owner: %w", err)
	}
	if err := d.Set("ctime", result.Ctime); err != nil {
		return fmt.Errorf("failed to set ctime: %w", err)
	}
	if err := d.Set("mtime", result.Mtime); err != nil {
		return fmt.Errorf("failed to set mtime: %w", err)
	}
	if err := d.Set("url", result.Url); err != nil {
		return fmt.Errorf("failed to set url: %w", err)
	}
	if err := d.Set("urls", result.Urls); err != nil {
		return fmt.Errorf("failed to set urls: %w", err)
	}
	if err := d.Set("execution_count", result.ExecutionCount); err != nil {
		return fmt.Errorf("failed to set execution_count: %w", err)
	}
	if err := d.Set("running_execution_count", result.RunningExecutionCount); err != nil {
		return fmt.Errorf("failed to set running_execution_count: %w", err)
	}
	if err := d.Set("queued_execution_count", result.QueuedExecutionCount); err != nil {
		return fmt.Errorf("failed to set queued_execution_count: %w", err)
	}
	if err := d.Set("enabled_endpoint_count", result.EnabledEndpointCount); err != nil {
		return fmt.Errorf("failed to set enabled_endpoint_count: %w", err)
	}
	if err := d.Set("last_execution_ctime", result.LastExecutionCtime); err != nil {
		return fmt.Errorf("failed to set last_execution_ctime: %w", err)
	}
	// Convert environment_variables to map[string]string (if possible)
	envVars := map[string]string{}
	for k, v := range result.EnvironmentVariables {
		if v == nil {
			envVars[k] = ""
		} else {
			envVars[k] = fmt.Sprintf("%v", v)
		}
	}
	if err := d.Set("environment_variables", envVars); err != nil {
		return fmt.Errorf("failed to set environment_variables: %w", err)
	}
	if err := d.Set("execution_summary", result.ExecutionSummary); err != nil {
		return fmt.Errorf("failed to set execution_summary: %w", err)
	}
	if err := d.Set("repository", map[string]interface{}{
		"id":  fmt.Sprintf("%v", result.Repository.ID),
		"url": result.Repository.URL,
		"ref": result.Repository.Ref,
	}); err != nil {
		return fmt.Errorf("failed to set repository: %w", err)
	}
	tags := make([]map[string]interface{}, len(result.Tags))
	for i, t := range result.Tags {
		tags[i] = map[string]interface{}{"project": t.Project, "name": t.Name, "color": t.Color}
	}
	if err := d.Set("tags", tags); err != nil {
		return fmt.Errorf("failed to set tags: %w", err)
	}
	if err := d.Set("upload_store_id", result.UploadStoreID); err != nil {
		return fmt.Errorf("failed to set upload_store_id: %w", err)
	}
	if err := d.Set("read_only", result.ReadOnly); err != nil {
		return fmt.Errorf("failed to set read_only: %w", err)
	}
	if err := d.Set("yaml_path", result.YamlPath); err != nil {
		return fmt.Errorf("failed to set yaml_path: %w", err)
	}
	if err := d.Set("name", result.Name); err != nil {
		return fmt.Errorf("failed to set name: %w", err)
	}
	if err := d.Set("description", result.Description); err != nil {
		return fmt.Errorf("failed to set description: %w", err)
	}
	return nil
}
