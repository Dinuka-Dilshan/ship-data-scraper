package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type AppHttpClient struct {
	client http.Client
	apiKey string
}

func (appClient *AppHttpClient) Get(url string, jsonRes any) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")
	req.Header.Set("consumer-key", appClient.apiKey)

	req.Header.Set("Origin", "https://www.maersk.com")
	req.Header.Set("Referer", "https://www.maersk.com/")

	// Browser “client hints”
	req.Header.Set("sec-ch-ua", `"Chromium";v="142", "Google Chrome";v="142", "Not_A Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)

	// Fetch metadata headers
	req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")

	// Optional but good to mimic browser
	req.Header.Set("DNT", "1")
	req.Header.Set("Cache-Control", "no-cache")

	res, err := appClient.client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		slog.Error(res.Status)
		return errors.New("failed request")
	}

	defer res.Body.Close()

	return json.NewDecoder(res.Body).Decode(jsonRes)

}
