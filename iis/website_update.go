package iis

import (
	"context"
	"encoding/json"
	"fmt"
)

func (client Client) UpdateWebsite(ctx context.Context, update Website) (*Website, error) {
	url := fmt.Sprintf("/api/webserver/websites/%s", update.ID)
	res, err := httpPatch(ctx, client, url, update)
	if err != nil {
		return nil, err
	}
	var site Website
	err = json.Unmarshal(res, &site)
	if err != nil {
		return nil, err
	}
	return &site, nil
}
