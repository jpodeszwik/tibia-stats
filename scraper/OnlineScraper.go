package scraper

import (
	"log"
	"sync"
	"tibia-stats/tibia"
	"time"
)

const fetchPlayersInterval = 40 * time.Second
const maxOfflineDuration = 4 * time.Hour

type OnlineScraper struct {
	api      *tibia.ApiClient
	lastSeen map[string]time.Time
	m        sync.RWMutex
}

func NewOnlineScraper(api *tibia.ApiClient) *OnlineScraper {
	return &OnlineScraper{
		api:      api,
		lastSeen: make(map[string]time.Time),
		m:        sync.RWMutex{},
	}
}

func (ot *OnlineScraper) Start() {
	log.Printf("Starting")

	fetchPlayersTicker := time.NewTicker(fetchPlayersInterval)
	ot.fetchOnlinePlayers()

	go func() {
		for range fetchPlayersTicker.C {
			ot.fetchOnlinePlayers()
		}
	}()
}

func (ot *OnlineScraper) fetchOnlinePlayers() {
	start := time.Now()
	players := ot.GetLastSeen()
	newLastSeen := make(map[string]time.Time)
	lastSeenLimit := time.Now().Add(-maxOfflineDuration)

	if players != nil {
		for player, lastSeen := range players {
			if lastSeen.After(lastSeenLimit) {
				newLastSeen[player] = lastSeen
			}
		}
	}

	worlds, err := ot.api.FetchWorlds()
	if err != nil {
		log.Printf("Failed to fetch worlds %v", err)
		return
	}

	for _, world := range worlds {
		players, err := ot.api.FetchOnlinePlayers(world.Name)
		if err != nil {
			log.Printf("Failed to fetch players %v", err)
			continue
		}

		for _, player := range players {
			newLastSeen[player.Name] = time.Now()
		}
	}

	log.Printf("Finished fetching online players in %v, onlineCount %v", time.Since(start), len(newLastSeen))
	ot.m.Lock()
	defer ot.m.Unlock()
	ot.lastSeen = newLastSeen
}

func (ot *OnlineScraper) GetLastSeen() map[string]time.Time {
	ot.m.RLock()
	defer ot.m.RUnlock()
	return ot.lastSeen
}
