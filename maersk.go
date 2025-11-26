package scraper

import (
	"fmt"
	"time"
)

type MaerskClient struct {
	httpClient   AppHttpClient
	StartDate    string
	EndDate      string
	CarrierCodes string
}

type Vessel struct {
	CallSign           string `json:"callSign"`
	FlagISOCountryCode string `json:"flagISOCountryCode"`
	VesselIMONumber    string `json:"vesselIMONumber"`
	VesselMaerskCode   string `json:"vesselMaerskCode"`
	VesselName         string `json:"vesselName"`
}

type VesselsResponse struct {
	Vessels []Vessel `json:"vessels"`
}

type Schedule struct {
	IsoCountryCode                  string    `json:"isoCountryCode"`
	CountryName                     string    `json:"countryName"`
	UnLocationCode                  string    `json:"unLocationCode"`
	CityName                        string    `json:"cityName"`
	PortName                        string    `json:"portName"`
	PortCode                        string    `json:"portCode"`
	RegionCode                      string    `json:"regionCode,omitempty"`
	MarineContainerTerminalName     string    `json:"marineContainerTerminalName"`
	MarineContainerTerminalRKSTCode string    `json:"marineContainerTerminalRKSTCode"`
	MarineContainerTerminalGeoCode  string    `json:"marineContainerTerminalGeoCode"`
	ArrivalTime                     time.Time `json:"arrivalTime"`
	ArrivalTimingClassifier         string    `json:"arrivalTimingClassifier"`
	DepartureTime                   time.Time `json:"departureTime"`
	DepartureTimingClassifier       string    `json:"departureTimingClassifier"`
	ArrivalVoyageNumber             string    `json:"arrivalVoyageNumber"`
	DepartureVoyageNumber           string    `json:"departureVoyageNumber"`
	ArrivalServiceName              string    `json:"arrivalServiceName"`
	ArrivalServiceCode              string    `json:"arrivalServiceCode"`
	DepartureServiceName            string    `json:"departureServiceName"`
	DepartureServiceCode            string    `json:"departureServiceCode"`
}

type VesselSchedulesResponse struct {
	VesselSchedules []Schedule `json:"vesselSchedules"`
}

func (vr *VesselSchedulesResponse) toCsvFormat(vessel Vessel) *[][]string {
	res := make([][]string, 0, len(vr.VesselSchedules)+1)

	for _, schedule := range vr.VesselSchedules {
		res = append(res, []string{
			vessel.VesselName,
			vessel.CallSign,
			vessel.FlagISOCountryCode,
			vessel.VesselIMONumber,
			vessel.VesselMaerskCode,
			schedule.IsoCountryCode,
			schedule.CountryName,
			schedule.UnLocationCode,
			schedule.CityName,
			schedule.PortName,
			schedule.PortCode,
			schedule.RegionCode,
			schedule.MarineContainerTerminalName,
			schedule.MarineContainerTerminalRKSTCode,
			schedule.MarineContainerTerminalGeoCode,
			schedule.ArrivalTime.Format(time.RFC3339),
			schedule.ArrivalTimingClassifier,
			schedule.DepartureTime.Format(time.RFC3339),
			schedule.DepartureTimingClassifier,
			schedule.ArrivalVoyageNumber,
			schedule.DepartureVoyageNumber,
			schedule.ArrivalServiceName,
			schedule.ArrivalServiceCode,
			schedule.DepartureServiceName,
			schedule.DepartureServiceCode,
		})
	}

	return &res
}

func (m *MaerskClient) GetVessels() (*VesselsResponse, error) {
	url := fmt.Sprintf("https://api.maersk.com/synergy/schedules/active-vessels?carrierCodes=%v", m.CarrierCodes)
	var vesselsRes VesselsResponse
	err := m.httpClient.Get(url, &vesselsRes)
	if err != nil {
		return nil, err
	}
	return &vesselsRes, nil
}

func (m *MaerskClient) GetScedule(vessel Vessel) (*VesselSchedulesResponse, error) {
	url := fmt.Sprintf("https://api.maersk.com/synergy/schedules/vessel-schedules?vesselMaerskCode=%v&fromDate=%v&toDate=%v&carrierCodes=%v", vessel.VesselMaerskCode, m.StartDate, m.EndDate, m.CarrierCodes)
	var scheduleRes VesselSchedulesResponse
	err := m.httpClient.Get(url, &scheduleRes)
	if err != nil {
		return nil, err
	}
	return &scheduleRes, err
}
