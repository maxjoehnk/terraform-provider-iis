package iis

import (
	"context"
	"encoding/json"
	"fmt"
)

type DigestAuthentication struct {
	ID      string `json:"id"`
	Enabled bool   `json:"enabled"`
	Realm   string `json:"realm"`
}

func (digest DigestAuthentication) ToMap() map[string]interface{} {
	digestMap := make(map[string]interface{}, 2)
	digestMap["enabled"] = digest.Enabled
	digestMap["realm"] = digest.Realm

	return digestMap
}

func (client Client) UpdateDigestAuthentication(ctx context.Context, auth *DigestAuthentication) (*DigestAuthentication, error) {
	url := fmt.Sprintf("/api/webserver/authentication/digest-authentication/%s", auth.ID)
	res, err := httpPatch(ctx, client, url, &auth)
	if err != nil {
		return nil, err
	}
	var digest DigestAuthentication
	err = json.Unmarshal(res, &digest)
	if err != nil {
		return nil, err
	}
	return &digest, nil
}
