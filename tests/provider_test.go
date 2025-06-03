package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tacy-ops/terraform-provider-valohai/valohai"
)

// ProviderFactories is required by the acceptance testing framework
var ProviderFactories = map[string]func() (*schema.Provider, error){
	"valohai": func() (*schema.Provider, error) { return valohai.Provider(), nil },
}

func TestProvider(t *testing.T) {
	if err := valohai.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
