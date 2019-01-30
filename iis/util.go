package iis

import "github.com/hashicorp/terraform/helper/schema"

func getNestedMap(d *schema.ResourceData, key string) map[string]interface{} {
	return d.Get(key).([]interface{})[0].(map[string]interface{})
}
