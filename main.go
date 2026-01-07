package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	timeFormat = "15:04:05"
	dateFormat = "20060102"
)

type Stop struct {
	ID            string  `json:"stop_id"`
	Code          int     `json:"stop_code"`
	Name          string  `json:"stop_name"`
	Description   string  `json:"stop_desc"`
	Latitude      float64 `json:"stop_lat"`
	Longitude     float64 `json:"stop_lon"`
	ZoneID        string  `json:"zone_id"`
	Url           string  `json:"stop_url"`
	LocationType  int     `json:"location_type"`
	ParentStation string  `json:"parent_station"`
	StopTimezone  string  `json:"stop_timezone"`
	Wheelchair    int     `json:"wheelchair_boarding"`
}

type Route struct {
	ID          string `json:"route_id"`
	AgencyID    string `json:"agency_id"`
	ShortName   string `json:"route_short_name"`
	LongName    string `json:"route_long_name"`
	Description string `json:"route_desc"`
	Type        int    `json:"route_type"`
	Url         string `json:"route_url"`
	Color       string `json:"route_color"`
	TextColor   string `json:"route_text_color"`
}

type Trip struct {
	RouteID     string `json:"route_id"`
	ServiceID   string `json:"service_id"`
	TripID      string `json:"trip_id"`
	Headsign    string `json:"trip_headsign"`
	ShortName   string `json:"trip_short_name"`
	DirectionID string `json:"direction_id"`
	BlockID     string `json:"block_id"`
	ShapeID     string `json:"shape_id"`
	Wheelchair  bool   `json:"wheelchair_accessible"`
}

type StopTime struct {
	TripID        string `json:"trip_id"`
	ArrivalTime   string `json:"arrival_time"`
	DepartureTime string `json:"departure_time"`
	StopID        string `json:"stop_id"`
	StopSequence  int    `json:"stop_sequence"`
}

type Calendar struct {
	ServiceID string `json:"service_id"`
	Monday    bool   `json:"monday"`
	Tuesday   bool   `json:"tuesday"`
	Wednesday bool   `json:"wednesday"`
	Thursday  bool   `json:"thursday"`
	Friday    bool   `json:"friday"`
	Saturday  bool   `json:"saturday"`
	Sunday    bool   `json:"sunday"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type CalendarDate struct {
	ServiceID     string `json:"service_id"`
	Date          int    `json:"date"`
	ExceptionType bool   `json:"exception_type"`
}

type Shape struct {
	ID                string  `json:"shape_id"`
	Latitude          float64 `json:"shape_pt_lat"`
	Longitude         float64 `json:"shape_pt_lon"`
	Sequence          int     `json:"shape_pt_sequence"`
	DistanceTravelled float64 `json:"shape_dist_traveled"`
}

//const staticFilesRoot = "web/build"

var stops []Stop

func main() {

	if len(os.Args) > 2 {
		fmt.Println("Usage: madrid-metro-api [port]")
		return
	}
	var port string
	if len(os.Args) == 2 {
		port = os.Args[1]
	} else {
		port = "8080"
	}

	mux := http.NewServeMux()
	//fileServerHandler := http.FileServer(http.Dir(staticFilesRoot))
	//mux.Handle("/app/", http.StripPrefix("/app", fileServerHandler))
	mux.HandleFunc("GET /api/healthz", checkReadiness)

	stops, err := loadStops("./data/stops.txt")
	if err != nil {
		log.Fatal(err)
	}

	mux.HandleFunc("GET /api/stops", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(stops)
	})

	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Println("Listening on port " + port)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func loadStops(path string) ([]Stop, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	header := map[string]int{}
	for i, col := range rows[0] {
		header[col] = i
	}

	var stops []Stop
	for _, row := range rows[1:] {
		stop := Stop{
			ID:            row[header["stop_id"]],
			Code:          toInt(row[header["stop_code"]]),
			Name:          row[header["stop_name"]],
			Description:   row[header["stop_desc"]],
			Latitude:      toFloat(row[header["stop_lat"]]),
			Longitude:     toFloat(row[header["stop_lon"]]),
			ZoneID:        row[header["zone_id"]],
			Url:           row[header["stop_url"]],
			LocationType:  toInt(row[header["location_type"]]),
			ParentStation: row[header["parent_station"]],
			StopTimezone:  row[header["stop_timezone"]],
			Wheelchair:    toInt(row[header["wheelchair_boarding"]]),
		}
		stops = append(stops, stop)
	}
	return stops, nil

}

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

func toFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func toBool(s string) bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func checkReadiness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func getStops(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
