package iis

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/maxjoehnk/microsoft-iis-administration"
)

const NameKey = "name"
const StatusKey = "status"

func resourceApplicationPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceApplicationPoolCreate,
		Read:   resourceApplicationPoolRead,
		Update: resourceApplicationPoolUpdate,
		Delete: resourceApplicationPoolDelete,

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

func resourceApplicationPoolCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*iis.Client)
	name := d.Get(NameKey).(string)
	pool, err := client.CreateAppPool(name)
	if err != nil {
		return err
	}
	d.SetId(pool.ID)
	return nil
}

func resourceApplicationPoolRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*iis.Client)
	id := d.Id()
	appPool, err := client.ReadAppPool(id)
	if err != nil {
		d.SetId("")
		return err
	}

	if err = d.Set(NameKey, appPool.Name); err != nil {
		return err
	}
	if err = d.Set(StatusKey, appPool.Status); err != nil {
		return err
	}
	return nil
}

func resourceApplicationPoolUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*iis.Client)
	if d.HasChange(NameKey) {
		applicationPool, err := client.UpdateAppPool(d.Id(), d.Get(NameKey).(string))
		if err != nil {
			return err
		}
		d.SetId(applicationPool.ID)
	}
	return nil
}

func resourceApplicationPoolDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*iis.Client)
	id := d.Id()
	err := client.DeleteAppPool(id)
	if err != nil {
		return err
	}
	return nil
}
