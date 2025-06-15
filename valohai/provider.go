package valohai

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// configureProvider configures the provider.
func configureProvider(d *schema.ResourceData) (interface{}, error) {
	// Retrieve the token from the configuration
	authToken := d.Get("token").(string)

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
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VALOHAI_API_TOKEN", nil),
				Description: "Valohai API token.",
				Sensitive:   true,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"valohai_project":   resourceProject(),
			"valohai_team":      resourceTeam(),
			"valohai_datastore": resourceDatastore(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"valohai_project": dataSourceProject(),
			"valohai_team":    dataSourceTeam(),
		},

		// Add ConfigureFunc to initialize the provider configuration
		ConfigureFunc: configureProvider,
	}
}

// (getOrganizationID and UserResponse moved to utils.go)
