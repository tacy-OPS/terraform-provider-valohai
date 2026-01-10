package tests

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tacy-ops/terraform-provider-valohai/valohai"
)

func TestProvider(t *testing.T) {
	if err := valohai.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProviderConfigureMissingToken(t *testing.T) {
	t.Setenv("VALOHAI_API_TOKEN", "")

	provider := valohai.Provider()
	data := schema.TestResourceDataRaw(t, provider.Schema, map[string]interface{}{})
	_, diags := provider.ConfigureContextFunc(context.Background(), data)
	if !diags.HasError() {
		t.Fatal("expected error when token is missing")
	}

	if len(diags) == 0 {
		t.Fatal("expected diagnostics to include error")
	}
	if !strings.Contains(diags[0].Summary, "token is required") {
		t.Fatalf("unexpected error: %s", diags[0].Summary)
	}
}
