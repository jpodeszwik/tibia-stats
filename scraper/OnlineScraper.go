package scraper

import (
	"sync"
	"tibia-stats/tibia"
	"tibia-stats/utils/logger"
	"time"
)

const fetchPlayersInterval = 40 * time.Second
const maxOfflineDuration = 4 * time.Hour

type OnlineScraper struct {
	api      *tibia.ApiClient
	worlds   *Worlds
	lastSeen map[string]time.Time
	m        sync.RWMutex
}

func NewOnlineScraper(api *tibia.ApiClient, worlds *Worlds) *OnlineScraper {
	return &OnlineScraper{
		api:      api,
		worlds:   worlds,
		lastSeen: make(map[string]time.Time),
		m:        sync.RWMutex{},
	}
}

func (ot *OnlineScraper) Start() {
	logger.Info.Printf("Starting")

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

	worlds := ot.worlds.getWorlds()
	for _, world := range worlds {
		onlinePlayers, err := retry(func() ([]tibia.OnlinePlayers, error) {
			return ot.api.FetchOnlinePlayers(world)
		}, 3)
		if err != nil {
			logger.Error.Printf("Failed to fetch players %v", err)
		} else {
			for _, player := range onlinePlayers {
				newLastSeen[player.Name] = time.Now()
			}
			break
		}
	}

	logger.Info.Printf("Finished fetching online players in %v, onlineCount %v", time.Since(start), len(newLastSeen))
	ot.m.Lock()
	ot.lastSeen = newLastSeen
	ot.m.Unlock()
}

func (ot *OnlineScraper) GetLastSeen() map[string]time.Time {
	ot.m.RLock()
	defer ot.m.RUnlock()
	return ot.lastSeen
}
