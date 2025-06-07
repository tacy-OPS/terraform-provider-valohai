package tests

import (
	"testing"

	"github.com/tacy-ops/terraform-provider-valohai/valohai"
)

func TestProvider(t *testing.T) {
	if err := valohai.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
