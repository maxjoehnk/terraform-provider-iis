package iis

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/maxjoehnk/microsoft-iis-administration"
)

const PathKey = "path"
const PhysicalPathKey = "physical_path"
const WebsiteKey = "website"
const ApplicationPoolKey = "application_pool"

func resourceApplication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationCreate,
		ReadContext:   resourceApplicationRead,
		UpdateContext: resourceApplicationUpdate,
		DeleteContext: resourceApplicationDelete,

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

func resourceApplicationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*iis.Client)
	request := createApplicationRequest(d)
	tflog.Debug(ctx, "Creating application: "+toJSON(request))
	application, err := client.CreateApplication(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Created application: "+toJSON(application))
	d.SetId(application.ID)
	return nil
}

func resourceApplicationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*iis.Client)
	application, err := client.ReadApplication(ctx, d.Id())
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Read application: "+toJSON(application))
	if err = d.Set(WebsiteKey, application.Website.ID); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set(ApplicationPoolKey, application.ApplicationPool.ID); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("location", application.Location); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceApplicationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceApplicationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*iis.Client)
	id := d.Id()
	tflog.Debug(ctx, "Deleting application: "+toJSON(id))
	err := client.DeleteApplication(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Deleted application: "+toJSON(id))
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
