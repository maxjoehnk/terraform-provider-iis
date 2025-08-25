package iis

import (
	"context"
	"encoding/json"
	"fmt"
)

func (r *Application) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ResourceReferences map[string]*ResourceReference

type ResourceReference struct {
	Href string `json:"href"`
}

type Application struct {
	Location         string               `json:"location"`
	Path             string               `json:"path"`
	ID               string               `json:"id"`
	PhysicalPath     string               `json:"physical_path"`
	EnabledProtocols string               `json:"enabled_protocols"`
	Website          ApplicationReference `json:"website"`
	ApplicationPool  ApplicationReference `json:"application_pool"`
	Links            ResourceReferences   `json:"_links,omitempty"`
}

type ApplicationReference struct {
	Name   string `json:"name"`
	ID     string `json:"id"`
	Status string `json:"status"`
}

func (client Client) ReadApplication(ctx context.Context, id string) (*Application, error) {
	url := fmt.Sprintf("/api/webserver/webapps/%s", id)
	var app Application
	if err := getJson(ctx, client, url, &app); err != nil {
		return nil, err
	}
	return &app, nil
}

func (client Client) DeleteApplication(ctx context.Context, id string) error {
	url := fmt.Sprintf("/api/webserver/webapps/%s", id)
	return httpDelete(ctx, client, url)
}
