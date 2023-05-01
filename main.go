package main

import (
	"tibia-stats/dynamo"
	"tibia-stats/scraper"
	"tibia-stats/tibia"
	"tibia-stats/tracker"
	"tibia-stats/utils/logger"
	"time"
)

func main() {
	start := time.Now()
	dr, err := dynamo.InitializeDeathRepository()
	if err != nil {
		logger.Error.Fatal(err)
	}
	guildExpRepository, err := dynamo.InitializeGuildExpRepository()
	if err != nil {
		logger.Error.Fatal(err)
	}

	guildsRepository, err := dynamo.InitializeGuildRepository()
	if err != nil {
		logger.Error.Fatal(err)
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

	logger.Info.Printf("Initialized in %v", time.Since(start))

	select {}
}
