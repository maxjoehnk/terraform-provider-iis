package iis

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/maxjoehnk/microsoft-iis-administration"
)

func dataSourceIisWebsite() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIisWebsiteRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceIisWebsiteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*iis.Client)

	sites, err := client.ListWebsites(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	siteIds := make([]string, 0)

	for _, site := range sites {
		siteIds = append(siteIds, site.ID)
	}

	d.SetId(resource.UniqueId())
	if err := d.Set("ids", siteIds); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
