package scraper

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
	req.Header.Set("User-Agent", "Mozilla/5.0")

	// REQUIRED by Maersk APIs
	req.Header.Set("Consumer-Key", appClient.apiKey)
	// or, if the API expects this variant instead:
	// req.Header.Set("x-maersk-consumer-key", appClient.apiKey)

	// Optional (only if required by the specific endpoint)
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", "Bearer <token>")
	// req.Header.Set("Maersk-Usi-Authorization", "<value>")
	// req.Header.Set("requestDate", "<RFC3339 timestamp>")
	// req.Header.Set("userId", "<id>")

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
