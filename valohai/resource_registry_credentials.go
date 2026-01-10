package valohai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceRegistryCredentials() *schema.Resource {
	return &schema.Resource{
		Create: resourceRegistryCredentialsCreate,
		Read:   resourceRegistryCredentialsRead,
		Update: resourceRegistryCredentialsUpdate,
		Delete: resourceRegistryCredentialsDelete,

		CustomizeDiff: validateRegistryCredentialsConfiguration(),

		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"docker", "aws-ecr", "aws-ecr-role", "gcp-cr"}, false,
				),
			},
			"image_pattern": {Type: schema.TypeString, Required: true},
			"owner":         {Type: schema.TypeInt, Optional: true},

			"configuration": {
				Type:      schema.TypeMap,
				Optional:  true,
				Sensitive: true,
				Elem:      &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func ResourceRegistryCredentials() *schema.Resource {
	return resourceRegistryCredentials()
}

var allowedConfigurationKeys = map[string]map[string]struct{}{
	"docker": {
		"username": {},
		"password": {},
		"version":  {},
	},
	"aws-ecr": {
		"access_key_id":     {},
		"secret_access_key": {},
		"region":            {},
		"version":           {},
	},
	"aws-ecr-role": {
		"role_name": {},
		"region":    {},
		"version":   {},
	},
	"gcp-cr": {
		"service_account_json": {},
		"version":              {},
	},
}

var requiredConfigurationKeys = map[string][]string{
	"docker":       {"username", "password"},
	"aws-ecr":      {"access_key_id", "secret_access_key", "region"},
	"aws-ecr-role": {"role_name", "region"},
	"gcp-cr":       {"service_account_json"},
}

var defaultConfigurationValues = map[string]map[string]string{
	"docker": {
		"version": "1",
	},
	"aws-ecr": {
		"version": "1",
	},
	"aws-ecr-role": {
		"version": "1",
	},
	"gcp-cr": {
		"version": "1",
	},
}

func cloneMap(in map[string]interface{}) map[string]interface{} {
	if in == nil {
		return map[string]interface{}{}
	}
	out := make(map[string]interface{}, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func validateRegistryCredentialsConfiguration() schema.CustomizeDiffFunc {
	return func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
		typ := d.Get("type").(string)

		allowedKeys, ok := allowedConfigurationKeys[typ]
		if !ok {
			return fmt.Errorf("unsupported type %q", typ)
		}

		conf := map[string]interface{}{}
		if raw, ok := d.GetOk("configuration"); ok && raw != nil {
			conf = cloneMap(raw.(map[string]interface{}))
		}

		conf = normalizeConfiguration(typ, conf)

		for k := range conf {
			if _, ok := allowedKeys[k]; !ok {
				return fmt.Errorf(
					"invalid configuration key %q for type %q (allowed keys: %s)",
					k, typ, strings.Join(keysOfSet(allowedKeys), ", "),
				)
			}
		}

		for _, rk := range requiredConfigurationKeys[typ] {
			v, exists := conf[rk]
			if !exists || strings.TrimSpace(fmt.Sprint(v)) == "" {
				return fmt.Errorf("missing or empty configuration.%s for type %q", rk, typ)
			}
		}

		return nil
	}
}

func normalizeConfiguration(typ string, conf map[string]interface{}) map[string]interface{} {
	conf = cloneMap(conf)

	if defs, ok := defaultConfigurationValues[typ]; ok {
		for k, dv := range defs {
			if v, exists := conf[k]; !exists || strings.TrimSpace(fmt.Sprint(v)) == "" {
				conf[k] = dv
			}
		}
	}
	return conf
}

func keysOfSet(m map[string]struct{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func resourceRegistryCredentialsCreate(d *schema.ResourceData, m interface{}) error {
	apiURL := "https://app.valohai.com/api/v0/registry-credentials/"
	authToken := m.(map[string]interface{})["token"].(string)

	payload := map[string]interface{}{
		"type":          d.Get("type").(string),
		"image_pattern": d.Get("image_pattern").(string),
	}

	if v, ok := d.GetOk("owner"); ok {
		payload["owner"] = v.(int)
	}

	conf := map[string]interface{}{}
	if v, ok := d.GetOk("configuration"); ok && v != nil {
		conf = cloneMap(v.(map[string]interface{}))
	}
	conf = normalizeConfiguration(d.Get("type").(string), conf)
	if len(conf) > 0 {
		payload["configuration"] = conf
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to encode payload: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+authToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return formatAPIError(resp)
	}

	var result struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	d.SetId(result.ID)
	return resourceRegistryCredentialsRead(d, m)
}

func resourceRegistryCredentialsRead(d *schema.ResourceData, m interface{}) error {
	authToken := m.(map[string]interface{})["token"].(string)
	id := d.Id()
	url := fmt.Sprintf("https://app.valohai.com/api/v0/registry-credentials/%s/", id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create GET request: %w", err)
	}
	req.Header.Set("Authorization", "Token "+authToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		return formatAPIError(resp)
	}

	var out map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return fmt.Errorf("failed to decode read response: %w", err)
	}

	_ = d.Set("type", out["type"])
	_ = d.Set("image_pattern", out["image_pattern"])
	_ = d.Set("owner", out["owner"])

	return nil
}

func resourceRegistryCredentialsUpdate(d *schema.ResourceData, m interface{}) error {
	authToken := m.(map[string]interface{})["token"].(string)
	apiURL := fmt.Sprintf("https://app.valohai.com/api/v0/registry-credentials/%s/", d.Id())

	payload := map[string]interface{}{
		"type":          d.Get("type").(string),
		"image_pattern": d.Get("image_pattern").(string),
	}

	if v, ok := d.GetOk("owner"); ok {
		payload["owner"] = v.(int)
	}

	conf := map[string]interface{}{}
	if v, ok := d.GetOk("configuration"); ok && v != nil {
		conf = cloneMap(v.(map[string]interface{}))
	}
	conf = normalizeConfiguration(d.Get("type").(string), conf)
	if len(conf) > 0 {
		payload["configuration"] = conf
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to encode payload: %w", err)
	}

	req, err := http.NewRequest("PUT", apiURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+authToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return formatAPIError(resp)
	}

	return resourceRegistryCredentialsRead(d, m)
}

func resourceRegistryCredentialsDelete(d *schema.ResourceData, m interface{}) error {
	authToken := m.(map[string]interface{})["token"].(string)
	id := d.Id()
	url := fmt.Sprintf("https://app.valohai.com/api/v0/registry-credentials/%s/", id)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create DELETE request: %w", err)
	}
	req.Header.Set("Authorization", "Token "+authToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute DELETE request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return formatAPIError(resp)
	}

	d.SetId("")
	return nil
}

func formatAPIError(resp *http.Response) error {
	var errResp map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&errResp)

	if nonField, ok := errResp["non_field_errors"].([]interface{}); ok && len(nonField) > 0 {
		if first, ok := nonField[0].(map[string]interface{}); ok {
			msg, _ := first["message"].(string)
			code, _ := first["code"].(string)
			if code != "" || msg != "" {
				return fmt.Errorf("API error %d (%s) - %s", resp.StatusCode, code, msg)
			}
		}
	}

	if detail, ok := errResp["detail"].(string); ok && detail != "" {
		return fmt.Errorf("API error %d - %s", resp.StatusCode, detail)
	}

	if len(errResp) > 0 {
		b, _ := json.Marshal(errResp)
		return fmt.Errorf("API error %d - %s", resp.StatusCode, string(b))
	}

	return fmt.Errorf("API error %d", resp.StatusCode)
}
