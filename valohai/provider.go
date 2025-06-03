package valohai

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// configureProvider configures the provider.
func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	_ = ctx
	// Retrieve the token from the configuration
	authToken := d.Get("token").(string)
	if authToken == "" {
		// Fallback to environment variable if not set in provider config
		authToken = os.Getenv("VALOHAI_API_TOKEN")
	}
	if authToken == "" {
		return nil, diag.Errorf("valohai provider token is required: set token in provider config or VALOHAI_API_TOKEN env var")
	}

	// Return an object containing the token for use in resources
	return map[string]interface{}{
		"token": authToken,
	}, nil
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VALOHAI_API_TOKEN", nil),
				Description: "Valohai API token.",
				Sensitive:   true,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"valohai_project": resourceProject(),
			"valohai_team":    resourceTeam(),
			"valohai_store":   resourceStore(),
			"valohai_registry_credentials": resourceRegistryCredentials(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"valohai_project": dataSourceProject(),
			"valohai_team":    dataSourceTeam(),
			"valohai_store":   dataSourceStore(),
		},

		ConfigureContextFunc: configureProvider,
	}
}
