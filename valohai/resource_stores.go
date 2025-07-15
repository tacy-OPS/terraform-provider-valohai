package valohai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceStore() *schema.Resource {
	return &schema.Resource{
		Create: resourceStoreCreate,
		Read:   resourceStoreRead,
		Update: resourceStoreUpdate,
		Delete: resourceStoreDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the store (max 64 characters)",
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
				Description: "Type of the store (s3, swift, azure, google)",
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
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return fmt.Sprintf("%v", old) == fmt.Sprintf("%v", new)
				},
			},
			"owner_id": {
				Type:     schema.TypeInt,
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
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceStore() *schema.Resource {
	return resourceStore()
}

func resourceStoreCreate(d *schema.ResourceData, m interface{}) error {
	accessMode := ""
	if v, ok := d.GetOk("access_mode"); ok {
		accessMode = v.(string)
	}
	hasTeams := d.Get("teams") != nil && len(d.Get("teams").([]interface{})) > 0
	hasProject := d.Get("project") != nil && d.Get("project").(string) != ""

	switch accessMode {
	case "owner_organization":
		if hasTeams || hasProject {
			return fmt.Errorf("with access_mode 'owner_organization', 'teams' and 'project' must not be set")
		}
	case "teams":
		if hasProject {
			return fmt.Errorf("with access_mode 'teams', 'project' must not be set")
		}
	case "single_project":
		if hasTeams {
			return fmt.Errorf("with access_mode 'single_project', 'teams' must not be set")
		}
	}

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
		conf := map[string]interface{}{}
		if m, ok := v.(map[string]interface{}); ok {
			if val, ok := m["bucket"]; ok {
				conf["bucket"] = val.(string)
			}
			if val, ok := m["access_key_id"]; ok {
				conf["access_key_id"] = val.(string)
			}
			if val, ok := m["secret_access_key"]; ok {
				conf["secret_access_key"] = val.(string)
			}
			if val, ok := m["region"]; ok {
				conf["region"] = val.(string)
			}
			if val, ok := m["multipart_upload_iam_role"]; ok {
				conf["multipart_upload_iam_role"] = val.(string)
			}
			if val, ok := m["endpoint_url"]; ok {
				conf["endpoint_url"] = val.(string)
			}
			if val, ok := m["role_arn"]; ok {
				conf["role_arn"] = val.(string)
			}
			if val, ok := m["kms_key_arn"]; ok {
				conf["kms_key_arn"] = val.(string)
			}
			if val, ok := m["use_presigned_put_object"]; ok {
				switch v := val.(type) {
				case bool:
					conf["use_presigned_put_object"] = v
				case string:
					conf["use_presigned_put_object"] = v == "true"
				}
			}
			if val, ok := m["insecure"]; ok {
				switch v := val.(type) {
				case bool:
					conf["insecure"] = v
				case string:
					conf["insecure"] = v == "true"
				}
			}
			if val, ok := m["skip_upload_file_name_check"]; ok {
				switch v := val.(type) {
				case bool:
					conf["skip_upload_file_name_check"] = v
				case string:
					conf["skip_upload_file_name_check"] = v == "true"
				}
			}
			if val, ok := m["test_configuration"]; ok {
				switch v := val.(type) {
				case bool:
					conf["test_configuration"] = v
				case string:
					conf["test_configuration"] = v == "true"
				}
			}
		}
		payload["configuration"] = conf
	}
	if v, ok := d.GetOk("owner_id"); ok {
		payload["owner"] = v.(int)
	}
	if v, ok := d.GetOk("project"); ok {
		payload["project"] = v.(string)
	}
	if v, ok := d.GetOk("paths"); ok {
		paths := map[string]string{}
		if m, ok := v.(map[string]interface{}); ok {
			for k, val := range m {
				paths[k] = fmt.Sprintf("%v", val)
			}
		}
		payload["paths"] = paths
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
		msg := ""
		if m, ok := errResp["message"].(string); ok && m != "" {
			msg = m
		} else if m, ok := errResp["detail"].(string); ok && m != "" {
			msg = m
		} else if nonField, ok := errResp["non_field_errors"].([]interface{}); ok && len(nonField) > 0 {
			if s, ok := nonField[0].(string); ok {
				msg = s
			}
		}
		raw, _ := json.MarshalIndent(errResp, "", "  ")
		return fmt.Errorf("API error %d: %s\nFull response: %s", resp.StatusCode, msg, string(raw))
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
		Owner            int                    `json:"owner"`
		Paths            map[string]interface{} `json:"paths"`
		Teams            []string               `json:"teams"`
		Project          interface{}            `json:"project"`
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

	expectedConfKeys := map[string]bool{
		"bucket": true, "region": true, "access_key_id": true, "secret_access_key": true,
		"multipart_upload_iam_role": true, "endpoint_url": true, "role_arn": true, "kms_key_arn": true,
		"use_presigned_put_object": true, "insecure": true, "skip_upload_file_name_check": true, "test_configuration": true,
	}
	conf := map[string]string{}
	if tfConf, ok := d.GetOk("configuration"); ok {
		for k := range tfConf.(map[string]interface{}) {
			if !expectedConfKeys[k] {
				continue
			}
			if v, ok := result.Configuration[k]; ok {
				conf[k] = fmt.Sprintf("%v", v)
			} else {
				conf[k] = fmt.Sprintf("%v", tfConf.(map[string]interface{})[k])
			}
		}
	}
	d.Set("configuration", conf)

	d.Set("owner_id", result.Owner)

	// Harmonisation paths (merge config Terraform + API, string)
	paths := map[string]string{}
	if tfPaths, ok := d.GetOk("paths"); ok {
		for k := range tfPaths.(map[string]interface{}) {
			if v, ok := result.Paths[k]; ok {
				paths[k] = fmt.Sprintf("%v", v)
			} else {
				paths[k] = fmt.Sprintf("%v", tfPaths.(map[string]interface{})[k])
			}
		}
	}
	d.Set("paths", paths)
	d.Set("teams", result.Teams)
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
	d.Set("url", result.URL)
	return nil
}

func resourceStoreRead(d *schema.ResourceData, m interface{}) error {
	authToken := m.(map[string]interface{})["token"].(string)
	id := d.Id() // UUID store Valohai
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
		Owner            int                    `json:"owner"`
		Project          interface{}            `json:"project"` // <- interface{}
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
	d.Set("type", result.Type)
	d.Set("name", result.Name)
	d.Set("owner_id", result.Owner)
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

	expectedConfKeys := map[string]bool{
		"bucket": true, "region": true, "access_key_id": true, "secret_access_key": true,
		"multipart_upload_iam_role": true, "endpoint_url": true, "role_arn": true, "kms_key_arn": true,
		"use_presigned_put_object": true, "insecure": true, "skip_upload_file_name_check": true, "test_configuration": true,
	}
	conf := map[string]string{}
	if tfConf, ok := d.GetOk("configuration"); ok {
		for k := range tfConf.(map[string]interface{}) {
			if !expectedConfKeys[k] {
				continue
			}
			if v, ok := result.Configuration[k]; ok {
				conf[k] = fmt.Sprintf("%v", v)
			} else {
				conf[k] = fmt.Sprintf("%v", tfConf.(map[string]interface{})[k])
			}
		}
	}
	d.Set("configuration", conf)
	paths := map[string]string{}
	if tfPaths, ok := d.GetOk("paths"); ok {
		for k := range tfPaths.(map[string]interface{}) {
			if v, ok := result.Paths[k]; ok {
				paths[k] = fmt.Sprintf("%v", v)
			} else {
				paths[k] = fmt.Sprintf("%v", tfPaths.(map[string]interface{})[k])
			}
		}
	}
	d.Set("paths", paths)
	d.Set("teams", result.Teams)
	d.Set("allow_read", result.AllowRead)
	d.Set("allow_write", result.AllowWrite)
	d.Set("access_mode", result.AccessMode)
	d.Set("allow_uri_download", result.AllowURIDownload)

	return nil
}

func resourceStoreUpdate(d *schema.ResourceData, m interface{}) error {
	id := d.Id() // UUID store Valohai
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
		conf := map[string]interface{}{}
		if m, ok := v.(map[string]interface{}); ok {
			if val, ok := m["bucket"]; ok {
				conf["bucket"] = val.(string)
			}
			if val, ok := m["access_key_id"]; ok {
				conf["access_key_id"] = val.(string)
			}
			if val, ok := m["secret_access_key"]; ok {
				conf["secret_access_key"] = val.(string)
			}
			if val, ok := m["region"]; ok {
				conf["region"] = val.(string)
			}
			if val, ok := m["multipart_upload_iam_role"]; ok {
				conf["multipart_upload_iam_role"] = val.(string)
			}
			if val, ok := m["endpoint_url"]; ok {
				conf["endpoint_url"] = val.(string)
			}
			if val, ok := m["role_arn"]; ok {
				conf["role_arn"] = val.(string)
			}
			if val, ok := m["kms_key_arn"]; ok {
				conf["kms_key_arn"] = val.(string)
			}
			if val, ok := m["use_presigned_put_object"]; ok {
				switch v := val.(type) {
				case bool:
					conf["use_presigned_put_object"] = v
				case string:
					conf["use_presigned_put_object"] = v == "true"
				}
			}
			if val, ok := m["insecure"]; ok {
				switch v := val.(type) {
				case bool:
					conf["insecure"] = v
				case string:
					conf["insecure"] = v == "true"
				}
			}
			if val, ok := m["skip_upload_file_name_check"]; ok {
				switch v := val.(type) {
				case bool:
					conf["skip_upload_file_name_check"] = v
				case string:
					conf["skip_upload_file_name_check"] = v == "true"
				}
			}
			if val, ok := m["test_configuration"]; ok {
				switch v := val.(type) {
				case bool:
					conf["test_configuration"] = v
				case string:
					conf["test_configuration"] = v == "true"
				}
			}
		}
		payload["configuration"] = conf
	}
	if v, ok := d.GetOk("owner_id"); ok {
		payload["owner"] = v.(int)
	}
	if v, ok := d.GetOk("paths"); ok {
		paths := map[string]string{}
		if m, ok := v.(map[string]interface{}); ok {
			for k, val := range m {
				paths[k] = fmt.Sprintf("%v", val)
			}
		}
		payload["paths"] = paths
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
		msg := ""
		if m, ok := errResp["message"].(string); ok && m != "" {
			msg = m
		} else if m, ok := errResp["detail"].(string); ok && m != "" {
			msg = m
		} else if nonField, ok := errResp["non_field_errors"].([]interface{}); ok && len(nonField) > 0 {
			if s, ok := nonField[0].(string); ok {
				msg = s
			}
		}
		raw, _ := json.MarshalIndent(errResp, "", "  ")
		return fmt.Errorf("API error %d: %s\nFull response: %s", resp.StatusCode, msg, string(raw))
	}

	// Decode response
	var result struct {
		AccessMode       string                 `json:"access_mode"`
		AllowRead        bool                   `json:"allow_read"`
		AllowURIDownload bool                   `json:"allow_uri_download"`
		AllowWrite       bool                   `json:"allow_write"`
		Configuration    map[string]interface{} `json:"configuration"`
		Name             string                 `json:"name"`
		Owner            int                    `json:"owner"`
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

	expectedConfKeys := map[string]bool{
		"bucket": true, "region": true, "access_key_id": true, "secret_access_key": true,
		"multipart_upload_iam_role": true, "endpoint_url": true, "role_arn": true, "kms_key_arn": true,
		"use_presigned_put_object": true, "insecure": true, "skip_upload_file_name_check": true, "test_configuration": true,
	}
	conf := map[string]string{}
	if tfConf, ok := d.GetOk("configuration"); ok {
		for k := range tfConf.(map[string]interface{}) {
			if !expectedConfKeys[k] {
				continue
			}
			if v, ok := result.Configuration[k]; ok {
				conf[k] = fmt.Sprintf("%v", v)
			} else {
				conf[k] = fmt.Sprintf("%v", tfConf.(map[string]interface{})[k])
			}
		}
	}
	d.Set("configuration", conf)
	d.Set("name", result.Name)
	d.Set("owner_id", result.Owner)

	paths := map[string]string{}
	if tfPaths, ok := d.GetOk("paths"); ok {
		for k := range tfPaths.(map[string]interface{}) {
			if v, ok := result.Paths[k]; ok {
				paths[k] = fmt.Sprintf("%v", v)
			} else {
				paths[k] = fmt.Sprintf("%v", tfPaths.(map[string]interface{})[k])
			}
		}
	}
	d.Set("paths", paths)
	d.Set("teams", result.Teams)
	return nil
}

func resourceStoreDelete(d *schema.ResourceData, m interface{}) error {
	authToken := m.(map[string]interface{})["token"].(string)
	id := d.Id() // UUID Store Valohai
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
	return fmt.Errorf("unexpected status code: %d while deleting store %s", resp.StatusCode, id)
}
