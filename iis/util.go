package iis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func getJson(ctx context.Context, client Client, path string, r interface{}) error {
	data, err := httpGet(ctx, client, path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &r)
}

func httpGet(ctx context.Context, client Client, path string) ([]byte, error) {
	response, err := request(ctx, client, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	return fetchBody(response)
}

func httpPost(ctx context.Context, client Client, path string, body interface{}) ([]byte, error) {
	response, err := request(ctx, client, "POST", path, body)
	if err != nil {
		return nil, err
	}
	return fetchBody(response)
}

func httpPatch(ctx context.Context, client Client, path string, body interface{}) ([]byte, error) {
	response, err := request(ctx, client, "PATCH", path, body)
	if err != nil {
		return nil, err
	}
	return fetchBody(response)
}

func httpDelete(ctx context.Context, client Client, path string) error {
	if _, err := request(ctx, client, "DELETE", path, nil); err != nil {
		return err
	}
	return nil
}

func buildRequest(ctx context.Context, client Client, method, path string, body interface{}) (*http.Request, error) {
	b := new(bytes.Buffer)
	if body != nil {
		if err := json.NewEncoder(b).Encode(body); err != nil {
			return nil, err
		}
	}

	requestUrl := fmt.Sprintf("%s%s", client.Host, path)
	req, err := http.NewRequestWithContext(ctx, method, requestUrl, b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Access-Token", fmt.Sprintf("Bearer %s", client.AccessKey))
	req.Header.Set("Accept", "application/hal+json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

func executeRequest(client Client, req *http.Request) (*http.Response, error) {
	response, err := client.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if err := guardStatusCode(req.Method, req.URL, response); err != nil {
		return nil, err
	}
	return response, err
}

func request(ctx context.Context, client Client, method, path string, body interface{}) (*http.Response, error) {
	req, err := buildRequest(ctx, client, method, path, body)
	if err != nil {
		return nil, err
	}
	return executeRequest(client, req)
}

func fetchBody(res *http.Response) ([]byte, error) {
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if err = res.Body.Close(); err != nil {
		return nil, err
	}
	return resBody, nil
}

func guardStatusCode(method string, url *url.URL, response *http.Response) error {
	if response.StatusCode < 200 || response.StatusCode > 400 {
		var body string
		if buffer, err := fetchBody(response); err == nil {
			body = string(buffer[:])
		}
		return fmt.Errorf("%s %s returned invalid status code: %s\n%s", method, url, response.Status, body)
	}
	return nil
}
