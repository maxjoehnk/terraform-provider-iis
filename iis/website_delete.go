package iis

import (
	"context"
	"fmt"
)

func (client Client) DeleteWebsite(ctx context.Context, id string) error {
	url := fmt.Sprintf("/api/webserver/websites/%s", id)
	return httpDelete(ctx, client, url)
}
