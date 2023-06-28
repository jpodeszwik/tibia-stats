package scraper

import (
	"tibia-stats/domain"
	"tibia-stats/tibia"
	"tibia-stats/utils/logger"
	"time"
)

const highScoreRefreshInterval = 2 * time.Hour

type HighScore struct {
	api     *tibia.ApiClient
	handler Handler[domain.WorldExperience]
	worlds  *Worlds
}

func (ge *HighScore) Start() {
	logger.Info.Printf("Starting")
	go func() {
		ticker := time.NewTicker(highScoreRefreshInterval)
		ge.fetchHighScores()
		for range ticker.C {
			ge.fetchHighScores()
		}
	}()
}

func (ge *HighScore) fetchHighScores() error {
	start := time.Now()
	worlds := ge.worlds.getWorlds()
	ret := make(map[string]int64)
	for _, world := range worlds {
		err := ge.fetchHighscore(world)
		if err != nil {
			logger.Error.Printf("Failed to fetch highscore for world %v", world)
			continue
		}
		logger.Debug.Printf("Finished fetching %v highscore", world)
	}

	logger.Info.Printf("Finished fetching %v worlds %v highScores in %v", len(worlds), len(ret), time.Since(start))
	return nil
}

func (ge *HighScore) fetchHighscore(world string) error {
	worldExperience := make(map[string]int64)
	for _, profession := range tibia.AllProfessions {
		for page := 0; page < 20; page++ {
			highScorePage, err := retry(func() ([]tibia.HighscoreResponse, error) {
				return ge.api.FetchHighscore(world, profession, tibia.Exp, page)
			}, 3)

			if err != nil {
				return err
			}
			for _, highScoreEntry := range highScorePage {
				worldExperience[highScoreEntry.Name] = highScoreEntry.Value
			}
		}
	}

	ge.handler(domain.WorldExperience{
		World:      world,
		Experience: worldExperience,
	})
	return nil
}

func NewHighScore(client *tibia.ApiClient, worlds *Worlds, handler Handler[domain.WorldExperience]) *HighScore {
	return &HighScore{
		api:     client,
		worlds:  worlds,
		handler: handler,
	}
}
