package main

import (
	"fmt"
	"strings"
)

func GetHeaderNames() []string {
	return []string{
		"Vessel Name",
		"Vessel Call Sign",
		"Vessel Flag ISO Country Code",
		"Vessel IMO Number",
		"Vessel Maersk Code",
		"Iso Country Code",
		"Country Name",
		"UnLocation Code",
		"City Name",
		"Port Name",
		"Port Code",
		"Region Code",
		"Marine Container Terminal Name",
		"Marine Container Terminal RKST Code",
		"Marine Container Terminal Geo Code",
		"Arrival Time",
		"Arrival Timing Classifier",
		"Departure Time",
		"Departure Timing Classifier",
		"Arrival Voyage Number",
		"Departure Voyage Number",
		"Arrival Service Name",
		"Arrival Service Code",
		"Departure Service Name",
		"Departure Service Code",
	}
}

func GetNodeList(schedules []Schedule) *[]string {
	nodes := make([]string, 0, len(schedules))

	for x := 0; x < len(schedules)-1; x++ {
		nodes = append(nodes, fmt.Sprintf("%v - %v", strings.ToLower(schedules[x].PortName), strings.ToLower(schedules[x+1].PortName)))
	}

	return &nodes
}
