package iis

import (
	"context"
	"encoding/json"
)

type CreateWebsiteRequest struct {
	Name            string               `json:"name"`
	PhysicalPath    string               `json:"physical_path"`
	Bindings        []WebsiteBinding     `json:"bindings"`
	ApplicationPool ApplicationReference `json:"application_pool"`
}

func (client Client) CreateWebsite(ctx context.Context, req CreateWebsiteRequest) (*Website, error) {
	res, err := httpPost(ctx, client, "/api/webserver/websites", req)
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
