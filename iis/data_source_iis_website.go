package iis

import (
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/maxjoehnk/microsoft-iis-administration"
)

func dataSourceIisWebsite() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIisWebsiteRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceIisWebsiteRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*iis.Client)

	sites, err := client.ListWebsites()
	if err != nil {
		return err
	}

	siteIds := make([]string, 0)

	for _, site := range sites {
		siteIds = append(siteIds, site.ID)
	}

	d.SetId(resource.UniqueId())
	if err := d.Set("ids", siteIds); err != nil {
		return err
	}

	return nil
}
