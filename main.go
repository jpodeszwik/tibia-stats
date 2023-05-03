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

	worldsScraper := scraper.NewWorlds(apiClient)
	worldsScraper.Start()

	guildScraper := scraper.NewGuilds(apiClient, worldsScraper, tracker.NewGuilds(guildsRepository))
	guildScraper.Start()

	onlineScraper := scraper.NewOnlineScraper(apiClient, worldsScraper)
	onlineScraper.Start()

	guildMembersScraper := scraper.NewGuildMembers(apiClient, worldsScraper, guildScraper)
	guildMembersScraper.Start()

	characterProfilesScraper := scraper.NewCharacterProfilesScraper(apiClient, onlineScraper, tracker.NewDeathTracker(dr))
	characterProfilesScraper.Start()

	guildExperienceScraper := scraper.NewGuildExperience(apiClient, worldsScraper, guildMembersScraper, tracker.NewGuildExp(guildExpRepository))
	guildExperienceScraper.Start()

	logger.Info.Printf("Initialized in %v", time.Since(start))

	select {}
}
