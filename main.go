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
	guildExpRepository, err := dynamo.InitializeGuildExpRepository()
	if err != nil {
		log.Fatal(err)
	}

	guildsRepository, err := dynamo.InitializeGuildRepository()
	if err != nil {
		log.Fatal(err)
	}

	apiClient := tibia.NewApiClient()
	ot := scraper.NewOnlineScraper(apiClient)
	ot.Start()

	dt := scraper.NewCharacterProfilesScraper(apiClient, ot, tracker.NewDeathTracker(dr))
	dt.Start()

	guildExperience := scraper.NewGuildExperience(apiClient, tracker.NewGuildExp(guildExpRepository))
	guildExperience.Start()

	guildScraper := scraper.NewGuilds(apiClient, tracker.NewGuilds(guildsRepository))
	guildScraper.Start()

	select {}
}
