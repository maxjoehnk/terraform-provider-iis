package iis

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/maxjoehnk/microsoft-iis-administration"
)

func Provider() *schema.Provider {
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
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	client := &iis.Client{
		Host:      d.Get("host").(string),
		AccessKey: d.Get("access_key").(string),
	}

	return client, nil
}
