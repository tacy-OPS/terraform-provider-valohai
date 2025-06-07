package tests

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tacy-ops/terraform-provider-valohai/valohai"
)

var ProviderFactories = map[string]func() (*schema.Provider, error){
	"valohai": func() (*schema.Provider, error) { return valohai.Provider(), nil },
}

func uniqueName(base string) string {
	return fmt.Sprintf("%s-%s", base, uuid.New().String())
}

func deleteTestStateFiles() {
	files := []string{"terraform.tfstate", "terraform.tfstate.backup"}
	for _, f := range files {
		if err := os.Remove(f); err != nil && !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "[deleteTestStateFiles] Failed to delete %s: %v\n", f, err)
		}
	}
}
