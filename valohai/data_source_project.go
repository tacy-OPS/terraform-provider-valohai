package valohai

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func dataSourceProject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceProjectRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceProjectRead(d *schema.ResourceData, m interface{}) error {
	id := d.Get("id").(string)
	d.SetId(id)
	d.Set("name", "Example Project")

	return nil
}
