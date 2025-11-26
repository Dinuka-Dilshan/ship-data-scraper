package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
    functions.HTTP("Handler", Handler)
}

func Handler(w http.ResponseWriter, r *http.Request) {

	var request struct {
		ApiKey       string `json:"apiKey"`
		CarrierCodes string `json:"carrierCodes"`
		StartDate    string `json:"StartDate"`
		EndDate      string `json:"endDate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {

	}

	client := &MaerskClient{
		httpClient: AppHttpClient{
			client: http.Client{Timeout: time.Minute * 2},
			apiKey: request.ApiKey,
		},
		CarrierCodes: request.CarrierCodes,
		StartDate:    request.StartDate,
		EndDate:      request.EndDate,
	}

	vesells, err := client.GetVessels()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	data := make([][]string, 0)
	data = append(data, GetHeaderNames())

	nodes := make(map[string]int)

	type StreamData struct {
		csv   *[][]string
		nodes *[]string
	}

	dataStream := make(chan *StreamData, len(vesells.Vessels))
	var wg sync.WaitGroup

	for _, vessel := range vesells.Vessels {
		wg.Go(func() {
			schedule, err := client.GetScedule(vessel)

			if err != nil {
				slog.Error(err.Error())
				dataStream <- nil
				return
			}

			dataStream <- &StreamData{
				csv:   schedule.toCsvFormat(vessel),
				nodes: GetNodeList(schedule.VesselSchedules),
			}
		})
	}

	go func() {
		wg.Wait()
		close(dataStream)
	}()

	for d := range dataStream {
		if d != nil {
			data = append(data, *d.csv...)
			for _, key := range *(d.nodes) {
				nodes[key]++
			}
		}
	}

	nodestrs := make([][]string, 0)

	for key, value := range nodes {
		nodestrs = append(nodestrs, []string{
			key, fmt.Sprint(value),
		})
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", `attachment; filename="files.zip"`)

	createMultipleCsv(w, map[string][][]string{"nodes.csv": nodestrs, "data.csv": data})

}
