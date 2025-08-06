package iis

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/maxjoehnk/microsoft-iis-administration"
)

const nameKey = "name"
const physicalPathKey = "physical_path"
const bindingsKey = "binding"
const appPoolKey = "application_pool"

const bindingProtocolKey = "protocol"
const bindingPortKey = "port"
const bindingAddressKey = "ip_address"
const bindingHostKey = "hostname"

func resourceWebsite() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWebsiteCreate,
		ReadContext:   resourceWebsiteRead,
		UpdateContext: resourceWebsiteUpdate,
		DeleteContext: resourceWebsiteDelete,

		Schema: map[string]*schema.Schema{
			nameKey: {
				Type:     schema.TypeString,
				Required: true,
			},
			physicalPathKey: {
				Type:     schema.TypeString,
				Required: true,
			},
			appPoolKey: {
				Type:     schema.TypeString,
				Optional: true,
			},
			bindingsKey: {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem:     bindingSchema,
			},
		},
	}
}

var bindingSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		bindingProtocolKey: {
			Type:     schema.TypeString,
			Default:  "http",
			Optional: true,
		},
		bindingPortKey: {
			Type:     schema.TypeInt,
			Default:  80,
			Optional: true,
		},
		bindingAddressKey: {
			Type:     schema.TypeString,
			Default:  "*",
			Optional: true,
		},
		bindingHostKey: {
			Type:     schema.TypeString,
			Optional: true,
		},
	},
}

func resourceWebsiteCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*iis.Client)
	request := createWebsiteRequest(d)
	tflog.Debug(ctx, "Creating website: "+toJSON(request))
	site, err := client.CreateWebsite(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Created website: "+toJSON(site))
	d.SetId(site.ID)
	return nil
}

func resourceWebsiteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*iis.Client)
	site, err := client.ReadWebsite(ctx, d.Id())
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Read website:"+toJSON(site))
	if err = d.Set(nameKey, site.Name); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set(physicalPathKey, site.PhysicalPath); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set(appPoolKey, site.ApplicationPool.ID); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set(bindingsKey, mapBindingsToSet(site)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceWebsiteUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceWebsiteDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*iis.Client)
	id := d.Id()
	err := client.DeleteWebsite(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func createWebsiteRequest(d *schema.ResourceData) iis.CreateWebsiteRequest {
	name := d.Get(nameKey).(string)
	physicalPath := d.Get(physicalPathKey).(string)
	bindings := d.Get(bindingsKey).(*schema.Set)
	request := iis.CreateWebsiteRequest{
		Name:         name,
		PhysicalPath: physicalPath,
		Bindings:     getBindings(bindings),
	}
	appPool := d.Get(appPoolKey)
	if appPool != nil {
		request.ApplicationPool = iis.ApplicationReference{
			ID: appPool.(string),
		}
	}
	return request
}

func getBindings(b *schema.Set) []iis.WebsiteBinding {
	bindings := make([]iis.WebsiteBinding, b.Len())
	for i, entry := range b.List() {
		binding := entry.(map[string]interface{})
		protocol := binding[bindingProtocolKey].(string)
		port := binding[bindingPortKey].(int)
		ipAddress := binding[bindingAddressKey].(string)
		hostname := binding[bindingHostKey].(string)

		bindings[i] = iis.WebsiteBinding{
			Port:      port,
			IPAddress: ipAddress,
			Hostname:  hostname,
			Protocol:  protocol,
		}
	}

	return bindings
}

func mapBindingsToSet(site *iis.Website) *schema.Set {
	var bindings []interface{}
	for _, binding := range site.Bindings {
		bindings = append(bindings, map[string]interface{}{
			bindingProtocolKey: binding.Protocol,
			bindingAddressKey:  binding.IPAddress,
			bindingPortKey:     binding.Port,
			bindingHostKey:     binding.Hostname,
		})
	}
	set := schema.NewSet(hashBinding, bindings)
	return set
}

func hashBinding(v interface{}) int {
	bindingMap := v.(map[string]interface{})
	address := schema.HashString(bindingMap[bindingAddressKey].(string))
	protocol := schema.HashString(bindingMap[bindingProtocolKey].(string))
	port := schema.HashInt(bindingMap[bindingPortKey].(int))
	hostname := schema.HashString(bindingMap[bindingHostKey].(string))

	return address + protocol + port + hostname
}
