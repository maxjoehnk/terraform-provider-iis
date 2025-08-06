package iis

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/maxjoehnk/microsoft-iis-administration"
)

const NameKey = "name"
const StatusKey = "status"

func resourceApplicationPool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationPoolCreate,
		ReadContext:   resourceApplicationPoolRead,
		UpdateContext: resourceApplicationPoolUpdate,
		DeleteContext: resourceApplicationPoolDelete,

		Schema: map[string]*schema.Schema{
			NameKey: {
				Type:     schema.TypeString,
				Required: true,
			},
			StatusKey: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "started",
			},
		},
	}
}

func resourceApplicationPoolCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*iis.Client)
	name := d.Get(NameKey).(string)
	tflog.Debug(ctx, "Creating application pool", map[string]interface{}{
		"name": name,
	})
	pool, err := client.CreateAppPool(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Created application pool", map[string]interface{}{
		"pool": pool,
	})
	d.SetId(pool.ID)
	return nil
}

func resourceApplicationPoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*iis.Client)
	id := d.Id()
	appPool, err := client.ReadAppPool(ctx, id)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Read application pool", map[string]interface{}{
		"appPool": appPool,
	})

	if err = d.Set(NameKey, appPool.Name); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set(StatusKey, appPool.Status); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceApplicationPoolUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*iis.Client)
	if d.HasChange(NameKey) {
		name := d.Get(NameKey).(string)
		tflog.Debug(ctx, "Updating application pool", map[string]interface{}{
			"id":   d.Id(),
			"name": name,
		})
		applicationPool, err := client.UpdateAppPool(ctx, d.Id(), name)
		if err != nil {
			return diag.FromErr(err)
		}
		tflog.Debug(ctx, "Updated application pool", map[string]interface{}{
			"applicationPool": applicationPool,
		})
		d.SetId(applicationPool.ID)
	}
	return nil
}

func resourceApplicationPoolDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*iis.Client)
	id := d.Id()
	tflog.Debug(ctx, "Deleting application pool", map[string]interface{}{
		"id": id,
	})
	err := client.DeleteAppPool(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Deleted application pool", map[string]interface{}{
		"id": id,
	})
	return nil
}
