package scraper

import (
	"fmt"
	"sync"
	"tibia-stats/tibia"
	"tibia-stats/utils/logger"
	"time"
)

const refreshProfilesInterval = 30 * time.Minute
const workers = 4

type CharacterProfilesScraper struct {
	ot      *OnlineScraper
	handler Handler[*tibia.Characters]
	api     *tibia.ApiClient
}

func NewCharacterProfilesScraper(api *tibia.ApiClient, ot *OnlineScraper, handler Handler[*tibia.Characters]) *CharacterProfilesScraper {
	return &CharacterProfilesScraper{
		ot:      ot,
		api:     api,
		handler: handler,
	}
}

func (dt *CharacterProfilesScraper) Start() {
	logger.Info.Printf("Starting")
	go func() {
		refreshProfilesTicker := time.NewTicker(refreshProfilesInterval)
		dt.refreshProfiles()
		for range refreshProfilesTicker.C {
			dt.refreshProfiles()
		}
	}()
}

func (dt *CharacterProfilesScraper) refreshProfiles() {
	charactersToRefresh := dt.ot.GetLastSeen()
	start := time.Now()

	fetchProfileWork := make(chan string, 100)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			dt.work(fetchProfileWork)
			wg.Done()
		}()
	}

	profilesFetched := 0
	for characterName := range charactersToRefresh {
		fetchProfileWork <- characterName
		profilesFetched++
		if profilesFetched%1000 == 0 {
			logger.Info.Printf("Fetched %v out of %v profiles", profilesFetched, len(charactersToRefresh))
		}
	}
	close(fetchProfileWork)
	wg.Wait()
	logger.Info.Printf("Finished fetching %v character profiles in %v with %v workers", profilesFetched, time.Since(start), workers)
}

func (dt *CharacterProfilesScraper) work(fetchProfilesWork <-chan string) {
	for characterName := range fetchProfilesWork {
		character, err := dt.fetchCharacter(characterName)
		if err != nil {
			logger.Error.Printf("Failed to fetch character of %v %v", characterName, err)
			continue
		}

		dt.handler(character)
	}
}

func (dt *CharacterProfilesScraper) fetchCharacter(characterName string) (*tibia.Characters, error) {
	var errs []error
	for i := 0; i < 3; i++ {
		character, err := dt.api.FetchCharacter(characterName)
		if err == nil {
			return character, nil
		}
	}
	return nil, fmt.Errorf("failed to fetch character %v %v", characterName, errs)
}
