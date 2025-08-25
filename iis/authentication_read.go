package iis

import (
	"context"
	"fmt"
)

type Authentication struct {
	ID    string              `json:"id"`
	Links AuthenticationLinks `json:"_links"`
}

type AuthenticationLinks struct {
	Anonymous ResourceReference `json:"anonymous"`
	Basic     ResourceReference `json:"basic"`
	Digest    ResourceReference `json:"digest"`
	Windows   ResourceReference `json:"windows"`
}

func (client Client) ReadAuthentication(ctx context.Context, id string) (Authentication, error) {
	url := fmt.Sprintf("/api/webserver/authentication/%s", id)
	var auth Authentication
	if err := getJson(ctx, client, url, &auth); err != nil {
		return auth, err
	}
	return auth, nil
}
func (client Client) ReadAuthenticationFromApplication(ctx context.Context, applicationId string) (Authentication, error) {
	var auth Authentication
	application, err := client.ReadApplication(ctx, applicationId)
	if err != nil {
		return auth, err
	}
	url := application.Links["authentication"].Href
	if err := getJson(ctx, client, url, &auth); err != nil {
		return auth, err
	}
	return auth, nil
}

func (client Client) ReadAnonymousAuthentication(ctx context.Context, auth *Authentication) (AnonymousAuthentication, error) {
	var anonymous AnonymousAuthentication
	if err := getJson(ctx, client, auth.Links.Anonymous.Href, &anonymous); err != nil {
		return anonymous, err
	}
	return anonymous, nil
}

func (client Client) ReadBasicAuthentication(ctx context.Context, auth *Authentication) (BasicAuthentication, error) {
	var basic BasicAuthentication
	if err := getJson(ctx, client, auth.Links.Basic.Href, &basic); err != nil {
		return basic, err
	}
	return basic, nil
}

func (client Client) ReadDigestAuthentication(ctx context.Context, auth *Authentication) (DigestAuthentication, error) {
	var digest DigestAuthentication
	if err := getJson(ctx, client, auth.Links.Digest.Href, &digest); err != nil {
		return digest, err
	}
	return digest, nil
}

func (client Client) ReadWindowsAuthentication(ctx context.Context, auth *Authentication) (WindowsAuthentication, error) {
	var windows WindowsAuthentication
	if err := getJson(ctx, client, auth.Links.Windows.Href, &windows); err != nil {
		return windows, err
	}
	return windows, nil
}
