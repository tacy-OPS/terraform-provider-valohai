package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/tacy-ops/terraform-provider-valohai/valohai"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: valohai.Provider,
	})
}
