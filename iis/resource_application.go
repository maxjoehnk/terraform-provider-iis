package iis

import (
	"github.com/hashicorp/terraform/helper/schema"
	iis "github.com/maxjoehnk/microsoft-iis-administration"
)

const PathKey = "path"
const PhysicalPathKey = "physical_path"
const WebsiteKey = "website"
const ApplicationPoolKey = "application_pool"

func resourceApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceApplicationCreate,
		Read:   resourceApplicationRead,
		Update: resourceApplicationUpdate,
		Delete: resourceApplicationDelete,

		Schema: map[string]*schema.Schema{
			PathKey: {
				Type:     schema.TypeString,
				Required: true,
			},
			PhysicalPathKey: {
				Type:     schema.TypeString,
				Required: true,
			},
			WebsiteKey: {
				Type:     schema.TypeString,
				Required: true,
			},
			ApplicationPoolKey: {
				Type:     schema.TypeString,
				Optional: true,
			},
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceApplicationCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*iis.Client)
	request := createApplicationRequest(d)
	application, err := client.CreateApplication(request)
	if err != nil {
		return err
	}
	d.SetId(application.ID)
	return nil
}

func resourceApplicationRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*iis.Client)
	application, err := client.ReadApplication(d.Id())
	if err != nil {
		d.SetId("")
		return err
	}
	if err = d.Set(WebsiteKey, application.Website.ID); err != nil {
		return err
	}
	if err = d.Set(ApplicationPoolKey, application.ApplicationPool.ID); err != nil {
		return err
	}
	if err = d.Set("location", application.Location); err != nil {
		return err
	}
	return nil
}

func resourceApplicationUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceApplicationDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*iis.Client)
	id := d.Id()
	err := client.DeleteApplication(id)
	if err != nil {
		return err
	}
	return nil
}

func createApplicationRequest(d *schema.ResourceData) iis.CreateApplicationRequest {
	path := d.Get(PathKey).(string)
	physicalPath := d.Get(PhysicalPathKey).(string)
	websiteId := d.Get(WebsiteKey).(string)
	appPoolId := d.Get(ApplicationPoolKey)
	website := iis.Reference{ID: websiteId}
	var appPool iis.Reference
	if appPoolId != nil {
		appPool = iis.Reference{ID: appPoolId.(string)}
	}
	request := iis.CreateApplicationRequest{
		Path:            path,
		PhysicalPath:    physicalPath,
		Website:         website,
		ApplicationPool: appPool,
	}
	return request
}
