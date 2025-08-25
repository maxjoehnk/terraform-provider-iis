package iis

import (
	"context"
	"encoding/json"
	"fmt"
)

type WindowsAuthenticationProvider struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

type WindowsAuthentication struct {
	ID        string                          `json:"id"`
	Enabled   bool                            `json:"enabled"`
	Providers []WindowsAuthenticationProvider `json:"providers"`
}

func (windows WindowsAuthentication) ToMap() map[string]interface{} {
	providers := make([]string, 0)
	for _, provider := range windows.Providers {
		if provider.Enabled {
			providers = append(providers, provider.Name)
		}
	}
	windowsMap := make(map[string]interface{}, 1)
	windowsMap["enabled"] = windows.Enabled
	windowsMap["providers"] = providers

	return windowsMap
}

func (client Client) UpdateWindowsAuthentication(ctx context.Context, auth *WindowsAuthentication) (*WindowsAuthentication, error) {
	url := fmt.Sprintf("/api/webserver/authentication/windows-authentication/%s", auth.ID)
	res, err := httpPatch(ctx, client, url, &auth)
	if err != nil {
		return nil, err
	}
	var windows WindowsAuthentication
	err = json.Unmarshal(res, &windows)
	if err != nil {
		return nil, err
	}
	return &windows, nil
}
