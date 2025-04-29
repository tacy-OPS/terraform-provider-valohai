package valohai

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Fonction pour configurer le provider
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	// Récupérer le token depuis la configuration
	authToken := d.Get("token").(string)

	// Retourner un objet contenant le token pour l'utiliser dans les ressources
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
			"valohai_project": resourceProject(),
		},

		DataSourcesMap: map[string]*schema.Resource{},

		// Ajout du ConfigureFunc pour initialiser la configuration du provider
		ConfigureFunc: providerConfigure,
	}
}
