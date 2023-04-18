package domain

import (
	"fmt"
	"log"
	"tibia-stats/repository"
	"tibia-stats/slices"
	"tibia-stats/tibia"
	"time"
)

const retries = 3

type fetchedHighscoresResult struct {
	highscores []tibia.HighscoreResponse
	err        error
}

type fetchExperienceJob struct {
	world      string
	profession tibia.Profession
	page       int
}

func FetchExperience(ac *tibia.ApiClient, expRepository repository.ExpRepository) error {
	worlds, err := ac.FetchWorlds()
	if err != nil {
		return err
	}

	log.Printf("Found %v worlds", len(worlds))
	jobsCount := len(worlds) * 4 * 20

	log.Printf("%v jobs", jobsCount)

	jobs := make(chan fetchExperienceJob, jobsCount)
	for _, world := range worlds {
		for _, profession := range tibia.AllProfessions {
			for page := 1; page <= 20; page++ {
				jobs <- fetchExperienceJob{world: world.Name, profession: profession, page: page}
			}
		}
	}
	close(jobs)

	workers := 8

	fetchResults := make(chan fetchedHighscoresResult, jobsCount)
	defer close(fetchResults)

	log.Printf("Fetching highscores with %v workers", workers)
	for i := 0; i < workers; i++ {
		go func() {
			for job := range jobs {
				err := fetchHighscore(ac, job, fetchResults)
				if err != nil {
					log.Printf("Error fetching highscore %v", err)
				}
			}
		}()
	}

	allHighscores := make([]tibia.HighscoreResponse, 0)
	for i := 0; i < jobsCount; i++ {
		if i%100 == 0 && i != 0 {
			log.Printf("%v done", i)
		}
		res := <-fetchResults
		if res.err == nil {
			allHighscores = append(allHighscores, res.highscores...)
		}

	}

	expData := slices.MapSlice(allHighscores, func(hr tibia.HighscoreResponse) repository.ExpData {
		return repository.ExpData{
			Name: hr.Name,
			Exp:  hr.Value,
			Date: time.Now(),
		}
	})

	log.Printf("Storing Experiences")
	err = expRepository.StoreExperiences(expData)
	if err != nil {
		log.Printf("Error storing exp data")
	}

	log.Printf("done")
	return err
}

func fetchHighscore(ac *tibia.ApiClient, job fetchExperienceJob, fetchResults chan fetchedHighscoresResult) error {
	for i := 1; i <= retries; i++ {
		response, err := ac.FetchHighscore(job.world, job.profession, tibia.Exp, job.page)
		if err != nil {
			log.Printf("Error fetching highscore %v %v %v %v tries left", job.world, job.profession, job.page, retries-i)
			continue
		}
		fetchResults <- fetchedHighscoresResult{
			highscores: response,
			err:        err,
		}
		return nil
	}
	return fmt.Errorf("failed ot fetch highscore after %v tries", retries)
}
