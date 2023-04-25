package main

import (
	"log"
	"tibia-stats/dynamo"
	"tibia-stats/scraper"
	"tibia-stats/tibia"
	"tibia-stats/tracker"
)

func main() {
	dr, err := dynamo.InitializeDeathRepository()
	if err != nil {
		log.Fatal(err)
	}

	apiClient := tibia.NewApiClient()
	ot := scraper.NewOnlineScraper(apiClient)
	ot.Start()

	dh := tracker.NewDeathTracker(dr)
	dt := scraper.NewCharacterProfilesScraper(apiClient, ot, dh)
	dt.Start()

	select {}
}
