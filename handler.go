package nflrushing

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

// Handler implements http.Handler interface
// Its a placeholder to maintain config or details
// related to the server in general
type Handler struct {
	tmpl       *template.Template
	csvHeaders []string
}

// NewHandler initializes re-usable variables across requests
func NewHandler() Handler {
	var h Handler
	// templatePath holds the relative path of template (relative to main.go)
	templatePath := "../../page.html"
	h.tmpl = template.Must(template.ParseFiles(templatePath))
	log.Println("Done parsing template")
	h.csvHeaders = append(h.csvHeaders, "Player", "Team", "Pos", "Att", "Att/G", "Yds", "Yds/G", "TD", "Lng", "1st", "1st%", "20+", "40+", "FUM")
	return h
}

// ServeHTTP handles HTTP requests for this app
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only support HTTP GET and POST
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	log.Println("Received new request")
	var params Request
	var ss Stats
	var err error
	// For GET, fetch all stats from DB
	// For POST, take filter params from request and fetch matching stats
	if r.Method == http.MethodGet {
		if ss, err = FetchStats(); err != nil {
			log.Println("Error fetching stats from DB: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		params.ResponseFilter = "SearchAll"
	} else if r.Method == http.MethodPost {
		params.PlayerFilter = r.FormValue("player")
		params.SortFilter = r.FormValue("sort")
		params.OrderFilter = r.FormValue("order")
		params.ResponseFilter = r.FormValue("response")
		setSelectedOptions(&params)
		log.Printf("Form request: %#v\n", params)
		// Check if valid values are received from user
		if err = validateOptions(params); err != nil {
			log.Println("Invalid params: ", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if ss, err = FetchFilteredStats(params); err != nil {
			log.Println("Error fetching stats from DB: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	log.Println("Result set size: ", len(ss))
	// Handle the file download response first
	if params.ResponseFilter == "Download" {
		filename := fmt.Sprintf("rushing-stats-%d.csv", time.Now().Unix())
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment;filename="+filename)
		csvWriter := csv.NewWriter(w)
		if err = createCSVExport(csvWriter, h.csvHeaders, ss); err != nil {
			log.Println("Error writing CSV file: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	// If not file download, respond with template and data
	data := TemplateData{Request: params, Stats: ss}
	h.tmpl.Execute(w, data)
}

// createCSVExport populates the csv file for export with headers and data passed
func createCSVExport(csvWriter *csv.Writer, csvHeaders []string, ss Stats) (err error) {
	if err = csvWriter.Write(csvHeaders); err != nil {
		return err
	}
	for _, stat := range ss {
		rec := []string{
			stat.Player, stat.Team, stat.Position,
			strconv.Itoa(stat.Attempts),
			strconv.FormatFloat(stat.AttemptsPerGameAvg, 'E', -1, 64),
			fmt.Sprintf("%v", stat.Yards), // Use fmt as the parsed is not populated
			strconv.FormatFloat(stat.AvgYardsPerAttempt, 'E', -1, 64),
			strconv.Itoa(stat.Touchdowns),
			fmt.Sprintf("%v", stat.Longest), // Using fmt as the parsed fields are not populated
			strconv.Itoa(stat.FirstDown),
			strconv.FormatFloat(stat.PerFirstDown, 'E', -1, 64),
			strconv.Itoa(stat.Yards20),
			strconv.Itoa(stat.Yards40),
			strconv.Itoa(stat.Fumbles),
		}
		if err = csvWriter.Write(rec); err != nil {
			return err
		}
	}
	csvWriter.Flush()
	if err = csvWriter.Error(); err != nil {
		return err
	}
	return nil
}

// setSelectedOptions takes string filters from request and populates
// type specific options, that are needed for processing the request
func setSelectedOptions(params *Request) {
	switch params.SortFilter {
	case "Yards":
		params.IsYardsSortSelected = true
	case "Longest":
		params.IsLongestSortSelected = true
	case "Touchdowns":
		params.IsTouchdownsSortSelected = true
	default:
		params.IsDefaultSortSelected = true
	}
	// Default to ascending order
	if params.OrderFilter == "asc" || params.OrderFilter == "" {
		params.IsOrderAsc = true
	} else if params.OrderFilter == "desc" {
		params.IsOrderDesc = true
	}
}

// validateOptions checks if the user input is valid
// before we can proceed with processing the request
func validateOptions(params Request) error {
	if params.ResponseFilter != "Search" && params.ResponseFilter != "Download" {
		return fmt.Errorf("Invalid response filter: %s", params.ResponseFilter)
	}
	return nil
}
