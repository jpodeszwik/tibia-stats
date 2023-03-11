package actions

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

	chunks := 8
	expChunks := slices.SplitSlice(expData, chunks)
	res := make(chan error, chunks)
	defer close(res)

	for _, slice := range expChunks {
		go func(eds []repository.ExpData) {
			res <- expRepository.StoreExperiences(eds)
		}(slice)
	}

	for i := 0; i < chunks; i++ {
		err = <-res
		if nil != err {
			log.Printf("Error when storing chunk of data %v", err)
		}
	}

	log.Printf("Done")
	return nil
}
