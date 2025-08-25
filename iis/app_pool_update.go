package iis

import (
	"context"
	"encoding/json"
	"fmt"
)

func (client Client) UpdateAppPool(ctx context.Context, id, name string) (*ApplicationPool, error) {
	reqBody := struct {
		Name string `json:"name"`
	}{name}
	url := fmt.Sprintf("/api/webserver/application-pools/%s", id)
	res, err := httpPatch(ctx, client, url, reqBody)
	if err != nil {
		return nil, err
	}
	var pool ApplicationPool
	err = json.Unmarshal(res, &pool)
	if err != nil {
		return nil, err
	}
	return &pool, nil
}
