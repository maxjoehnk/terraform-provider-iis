package iis

import (
	"context"
	"fmt"
)

func (client Client) ReadWebsite(ctx context.Context, id string) (*Website, error) {
	url := fmt.Sprintf("/api/webserver/websites/%s", id)
	var site Website
	if err := getJson(ctx, client, url, &site); err != nil {
		return nil, err
	}
	return &site, nil
}
