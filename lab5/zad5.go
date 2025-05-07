package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Stop struct {
	StopName string `json:"stopDesc"`
	StopID   int    `json:"stopId"`
}

type Stops struct {
	Stops []Stop `json:"stops"`
}

type Departure struct {
	EstimatedTime string `json:"estimatedTime"`
}

type Departures struct {
	Departures []Departure `json:"departures"`
}

type Route struct {
	ArrivalTime string `json:"arrivalTime"`
}

type StopTimes struct {
	Times []Route `json:"stopTimes"`
}

func getStopByName(stops []Stop, name string) (*Stop, error) {
	for _, stop := range stops {
		if strings.ToLower(stop.StopName) == strings.ToLower(name) {
			return &stop, nil
		}
	}
	return nil, fmt.Errorf("stop with name %s not found", name)
}

func getDeparturesByStopId(id int) (*Departures, error) {
	url := "https://ckan2.multimediagdansk.pl/departures"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var departures map[int]Departures
	if err := json.Unmarshal(body, &departures); err != nil {
		return nil, err
	}

	if dep, exists := departures[id]; exists {
		return &dep, nil
	}
	return nil, fmt.Errorf("no departure found for ID")
}

func estimateTime(lineId string) {
	for true {
		currentTime := time.Now()
		formattedDate := currentTime.Format("2006-01-02")

		url := fmt.Sprintf("https://ckan2.multimediagdansk.pl/stopTimes?date=%s&routeId=%s", formattedDate, lineId)

		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}

		var stopTimes StopTimes
		if err := json.Unmarshal(body, &stopTimes); err != nil {
			fmt.Println("JSON unmarshal error:", err)
			return
		}

		counter := 0
		loopHandler := true
		for loopHandler {
			arrivalTimeStr := stopTimes.Times[counter].ArrivalTime
			layout := "2006-01-02T15:04:05"
			arrivalTime, err := time.Parse(layout, arrivalTimeStr)
			if err != nil {
				fmt.Println("Error parsing time:", err)
				return
			}
			currentTime := time.Now()

			if currentTime.After(arrivalTime) {
				duration := currentTime.Sub(arrivalTime).Minutes()
				fmt.Printf("Czas do następnego przystanku dla linii %s: %.2f\n", lineId, duration)
				loopHandler = false
			}
		}
		time.Sleep(2 * time.Second)
	}
}

func main() {
	url := "https://ckan.multimediagdansk.pl/dataset/c24aa637-3619-4dc2-a171-a23eec8f2172/resource/d3e96eb6-25ad-4d6c-8651-b1eb39155945/download/stopsingdansk.json"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var stops Stops
	if err := json.Unmarshal(body, &stops); err != nil {
		fmt.Println("JSON unmarshal error:", err)
		return
	}

	var stopName string
	fmt.Println("Podaj nazwę przystanku: ")
	fmt.Scanln(&stopName)

	stop, err := getStopByName(stops.Stops, stopName)
	if err != nil {
		fmt.Println(err)
		return
	}

	url = "https://ckan2.multimediagdansk.pl/departures"
	resp, err = http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}
	departures, err := getDeparturesByStopId(stop.StopID)

	fmt.Println("Estymowane wyjazdy")
	for _, departure := range departures.Departures {
		fmt.Println(departure.EstimatedTime)
	}

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		estimateTime("3")
	}()

	go func() {
		defer wg.Done()
		estimateTime("5")
	}()

	wg.Wait()
}
