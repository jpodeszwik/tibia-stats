package tracker

import (
	"tibia-stats/domain"
	"tibia-stats/dynamo"
	"tibia-stats/utils/logger"
)

type PlayerExp struct {
	repo *dynamo.HighScoreRepository
}

func (pe *PlayerExp) HandleHighScore(experience domain.WorldExperience) {
	for {
		err := pe.repo.StoreHighScore(experience)
		if err != nil {
			logger.Error.Printf("Failed to store HighScore %v", err)
		} else {
			return
		}
	}
}

func NewPlayerExp(repo *dynamo.HighScoreRepository) *PlayerExp {
	return &PlayerExp{repo: repo}
}
