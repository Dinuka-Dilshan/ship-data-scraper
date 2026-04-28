package scraper

import (
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"golang.org/x/sync/semaphore"
)

func init() {
	functions.HTTP("Handler", Handler)
}

func Handler(w http.ResponseWriter, r *http.Request) {

	var request = struct {
		ApiKey       string
		CarrierCodes string
		StartDate    string
		EndDate      string
	}{
		ApiKey:       "uXe7bxTHLY0yY0e8jnS6kotShkLuAAqG",
		CarrierCodes: r.URL.Query().Get("carrierCodes"),
		StartDate:    r.URL.Query().Get("startDate"),
		EndDate:      r.URL.Query().Get("endDate"),
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	semaphore := semaphore.NewWeighted(5)
	var wg sync.WaitGroup

	go func() {
		for _, vessel := range vesells.Vessels {
			if err := semaphore.Acquire(r.Context(), 1); err != nil {
				slog.Error(err.Error())
				continue
			}
			wg.Go(func() {
				defer semaphore.Release(1)
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
