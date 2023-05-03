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
	logger.Info.Printf("Starting")

	deathRepository, err := dynamo.InitializeDeathRepository()
	if err != nil {
		logger.Error.Fatal(err)
	}
	guildExpRepository, err := dynamo.InitializeGuildExpRepository()
	if err != nil {
		logger.Error.Fatal(err)
	}

	guildRepository, err := dynamo.InitializeGuildRepository()
	if err != nil {
		logger.Error.Fatal(err)
	}

	highScoreRepository, err := dynamo.InitializeHighScoreRepository()
	if err != nil {
		logger.Error.Fatal(err)
	}

	logger.Info.Printf("Initializing repositories")

	logger.Info.Printf("Initializing trackers")
	guildsTracker := tracker.NewGuilds(guildRepository)
	deathTracker := tracker.NewDeathTracker(deathRepository)
	guildExpTracker := tracker.NewGuildExp(guildExpRepository)
	playerExpTracker := tracker.NewPlayerExp(highScoreRepository)

	apiClient := tibia.NewApiClient()

	logger.Info.Printf("Initializing scrapers")
	worldsScraper := scraper.NewWorlds(apiClient)
	onlineScraper := scraper.NewOnlineScraper(apiClient, worldsScraper)
	guildScraper := scraper.NewGuilds(apiClient, worldsScraper, guildsTracker.Handle)
	guildMembersScraper := scraper.NewGuildMembers(apiClient, guildScraper, guildExpTracker.HandleGuildMembers)
	characterProfilesScraper := scraper.NewCharacterProfilesScraper(apiClient, onlineScraper, deathTracker.Handle)
	highScoreScraper := scraper.NewHighScore(apiClient, worldsScraper, combineTrackers(guildExpTracker.HandleWorldExperience, playerExpTracker.HandleHighScore))

	logger.Info.Printf("Starting scrapers")
	worldsScraper.Start()
	onlineScraper.Start()
	guildScraper.Start()
	guildMembersScraper.Start()
	characterProfilesScraper.Start()
	highScoreScraper.Start()

	logger.Info.Printf("Initialized in %v", time.Since(start))

	select {}
}

func combineTrackers[T any](funcs ...func(T)) func(T) {
	return func(t T) {
		for _, f := range funcs {
			f(t)
		}
	}
}
