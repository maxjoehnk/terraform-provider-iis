package iis

import (
	"context"
	"encoding/json"
	"fmt"
)

func (client Client) CreateApplication(ctx context.Context, application CreateApplicationRequest) (*Application, error) {
	res, err := httpPost(ctx, client, "/api/webserver/webapps", application)
	if err != nil {
		return nil, err
	}
	var app Application
	err = json.Unmarshal(res, &app)
	if err != nil {
		return nil, err
	}
	fmt.Println(app)
	return &app, nil
}

type Reference struct {
	ID string `json:"id"`
}

type CreateApplicationRequest struct {
	Path            string    `json:"path"`
	PhysicalPath    string    `json:"physical_path"`
	Website         Reference `json:"website"`
	ApplicationPool Reference `json:"application_pool"`
}
