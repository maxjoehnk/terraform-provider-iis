package iis

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	iis "github.com/maxjoehnk/microsoft-iis-administration"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"iis_application_pool": resourceApplicationPool(),
			"iis_application":      resourceApplication(),
			"iis_authentication":   resourceAuthentication(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"iis_website": dataSourceIisWebsite(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client := &iis.Client{
		Host:      d.Get("host").(string),
		AccessKey: d.Get("access_key").(string),
	}

	return client, nil
}
