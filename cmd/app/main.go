package main

import (
	"log"
	"net/http"

	nflrushing "github.com/geekgunda/nfl-rushing"
)

func main() {
	var shouldImport bool
	var err error
	if err = nflrushing.InitDBConn(); err != nil {
		log.Fatalf("Failed connecting to DB: %v", err)
	}
	// Comment the line below, to avoid duplicate import upon restart
	shouldImport = true
	if shouldImport {
		if err = nflrushing.ImportStats(); err != nil {
			log.Fatalf("Failed to import stats: %v", err)
		}
	}

	h := nflrushing.NewHandler()
	http.Handle("/rushingstats", &h)
	log.Println("Starting HTTP Server")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("Error while listenAndServe: ", err)
	}
}
