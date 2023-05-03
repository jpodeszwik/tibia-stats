package tracker

import (
	"tibia-stats/domain"
	"tibia-stats/dynamo"
	"tibia-stats/utils/logger"
)

type HighScore struct {
	repo *dynamo.HighScoreRepository
}

func (pe *HighScore) HandleHighScore(experience domain.WorldExperience) {
	for {
		err := pe.repo.StoreHighScore(experience)
		if err != nil {
			logger.Error.Printf("Failed to store HighScore %v", err)
		} else {
			return
		}
	}
}

func NewHighScore(repo *dynamo.HighScoreRepository) *HighScore {
	return &HighScore{repo: repo}
}
