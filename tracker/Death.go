package tracker

import (
	"log"
	"sync"
	"tibia-stats/domain"
	"tibia-stats/dynamo"
	"tibia-stats/tibia"
	"time"
)

const isoDateTime = "2006-01-02T15:04:05Z"

type Death struct {
	lastDeath map[string]time.Time
	repo      *dynamo.DeathRepository
	m         sync.RWMutex
}

func NewDeathTracker(repo *dynamo.DeathRepository) *Death {
	return &Death{
		lastDeath: make(map[string]time.Time),
		repo:      repo,
		m:         sync.RWMutex{},
	}
}

func (dh *Death) HandleProfileRefresh(character *tibia.Characters) error {
	characterName := character.Character.Name
	if len(character.Deaths) == 0 {
		return nil
	}

	dh.m.RLock()
	lastDeathTime, exists := dh.lastDeath[characterName]
	dh.m.RUnlock()

	var minDeathTime time.Time
	if exists {
		minDeathTime = lastDeathTime
	} else {
		death, err := dh.repo.GetLastDeath(characterName)
		if err != nil {
			log.Printf("Failed to get last death")
			return err
		}

		if death == nil {
			minDeathTime = time.Time{}
		} else {
			minDeathTime = death.Time
		}
	}

	var deaths []domain.Death
	var maxDeathTime = minDeathTime
	for i := len(character.Deaths) - 1; i >= 0; i-- {
		death := character.Deaths[i]

		parsedTime, err := time.Parse(isoDateTime, death.Time)
		if err != nil {
			log.Printf("Failed to parse time %v", err)
			continue
		}

		if parsedTime.After(minDeathTime) {
			deaths = append(deaths, domain.Death{
				CharacterName: characterName,
				Guild:         character.Character.Guild.Name,
				Time:          parsedTime,
				Reason:        death.Reason,
			})
		}
		if parsedTime.After(maxDeathTime) {
			maxDeathTime = parsedTime
		}
	}

	if len(deaths) > 0 {
		for {
			err := dh.repo.StoreDeaths(deaths)
			if err == nil {
				break
			} else {
				log.Printf("Error %v storing deaths, retrying", err)
			}
		}
	}

	if !exists || maxDeathTime.After(lastDeathTime) {
		dh.m.Lock()
		dh.lastDeath[characterName] = maxDeathTime
		dh.m.Unlock()
	}

	return nil
}
