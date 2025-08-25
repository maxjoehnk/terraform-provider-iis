package iis

import "net/http"

type Client struct {
	HttpClient http.Client
	Host       string
	AccessKey  string
}
