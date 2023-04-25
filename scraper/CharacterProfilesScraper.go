package scraper

import (
	"fmt"
	"log"
	"sync"
	"tibia-stats/tibia"
	"tibia-stats/tracker"
	"time"
)

const refreshProfilesInterval = 30 * time.Minute
const workers = 4

type CharacterProfilesScraper struct {
	ot  *OnlineScraper
	dh  *tracker.Death
	api *tibia.ApiClient
}

func NewCharacterProfilesScraper(api *tibia.ApiClient, ot *OnlineScraper, dh *tracker.Death) *CharacterProfilesScraper {
	return &CharacterProfilesScraper{
		ot:  ot,
		api: api,
		dh:  dh,
	}
}

func (dt *CharacterProfilesScraper) Start() {
	log.Printf("Starting")

	refreshProfilesTicker := time.NewTicker(refreshProfilesInterval)
	dt.refreshProfiles()

	go func() {
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
			log.Printf("Fetched %v out of %v profiles", profilesFetched, len(charactersToRefresh))
		}
	}
	close(fetchProfileWork)
	wg.Wait()
	log.Printf("Finished fetching %v character profiles in %v with %v workers", profilesFetched, time.Since(start), workers)
}

func (dt *CharacterProfilesScraper) work(fetchProfilesWork <-chan string) {
	for characterName := range fetchProfilesWork {
		character, err := dt.fetchCharacter(characterName)
		if err != nil {
			log.Printf("Failed to fetch character of %v %v", characterName, err)
			continue
		}

		err = dt.dh.HandleProfileRefresh(character)
		if err != nil {
			log.Printf("Failed to handleProfileRefresh for character %v %v", characterName, err)
		}
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