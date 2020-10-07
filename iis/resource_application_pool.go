package iis

import (
	"context"
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
			},
		},
	}
}

func resourceApplicationPoolCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*iis.Client)
	name := d.Get(NameKey).(string)
	pool, err := client.CreateAppPool(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}
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
		applicationPool, err := client.UpdateAppPool(ctx, d.Id(), d.Get(NameKey).(string))
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(applicationPool.ID)
	}
	return nil
}

func resourceApplicationPoolDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*iis.Client)
	id := d.Id()
	err := client.DeleteAppPool(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
