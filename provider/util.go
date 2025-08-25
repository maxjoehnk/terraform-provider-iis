package provider

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getList(d *schema.ResourceData, key string) []interface{} {
	return d.Get(key).([]interface{})
}

func getNestedMap(d *schema.ResourceData, key string) map[string]interface{} {
	return getList(d, key)[0].(map[string]interface{})
}

func hasNestedMap(d *schema.ResourceData, key string) bool {
	return len(getList(d, key)) == 1
}

func toJSON(obj interface{}) string {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}
