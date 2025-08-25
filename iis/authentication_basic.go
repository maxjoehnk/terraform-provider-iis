package iis

import (
	"context"
	"encoding/json"
	"fmt"
)

type BasicAuthentication struct {
	ID                 string `json:"id"`
	Enabled            bool   `json:"enabled"`
	DefaultLogonDomain string `json:"default_logon_domain"`
	Realm              string `json:"realm"`
}

func (basic BasicAuthentication) ToMap() map[string]interface{} {
	basicMap := make(map[string]interface{}, 3)
	basicMap["enabled"] = basic.Enabled
	basicMap["default_domain"] = basic.DefaultLogonDomain
	basicMap["realm"] = basic.Realm

	return basicMap
}

func (client Client) UpdateBasicAuthentication(ctx context.Context, auth *BasicAuthentication) (BasicAuthentication, error) {
	url := fmt.Sprintf("/api/webserver/authentication/basic-authentication/%s", auth.ID)
	var basic BasicAuthentication
	res, err := httpPatch(ctx, client, url, &auth)
	if err != nil {
		return basic, err
	}
	err = json.Unmarshal(res, &basic)
	if err != nil {
		return basic, err
	}
	return basic, nil
}
