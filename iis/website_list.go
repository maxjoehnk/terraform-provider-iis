package iis

import "context"

type WebsiteListItem struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type WebsiteListResponse struct {
	Websites []WebsiteListItem `json:"websites"`
}

func (client Client) ListWebsites(ctx context.Context) ([]WebsiteListItem, error) {
	var res WebsiteListResponse
	err := getJson(ctx, client, "/api/webserver/websites", &res)
	if err != nil {
		return nil, err
	}
	return res.Websites, nil
}
