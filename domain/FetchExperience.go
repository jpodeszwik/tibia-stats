package domain

import (
	"log"
	"tibia-exp-tracker/repository"
	"tibia-exp-tracker/slices"
	"tibia-exp-tracker/tibia"
	"time"
)

func FetchExperience(ac *tibia.ApiClient, expRepository repository.ExpRepository, world string) error {
	log.Printf("Fetching exp for %v", world)
	response, err := ac.FetchAllHighscores(world, tibia.Exp)
	if err != nil {
		return err
	}
	log.Printf("Fetched %d highscores", len(response))

	expData := slices.MapSlice(response, func(hr tibia.HighscoreResponse) repository.ExpData {
		return repository.ExpData{
			Name: hr.Name,
			Exp:  hr.Value,
			Date: time.Now(),
		}
	})

	err = expRepository.StoreExperiences(expData)

	return nil
}
