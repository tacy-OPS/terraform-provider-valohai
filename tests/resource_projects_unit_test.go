package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tacy-OPS/terraform-provider-valohai/valohai"
)

func TestGetOptionalString(t *testing.T) {
	d := schema.TestResourceDataRaw(t, valohai.ResourceProject().Schema, map[string]interface{}{
		"description": "foo",
	})
	if got := valohai.GetOptionalString(d, "description"); got != "foo" {
		t.Errorf("expected 'foo', got '%s'", got)
	}

	d = schema.TestResourceDataRaw(t, valohai.ResourceProject().Schema, map[string]interface{}{})
	if got := valohai.GetOptionalString(d, "description"); got != "" {
		t.Errorf("expected empty string, got '%s'", got)
	}
}
