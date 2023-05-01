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

	worlds, err := retry(func() ([]tibia.OverviewWorld, error) {
		return ot.api.FetchWorlds()
	}, 3)
	if err != nil {
		logger.Error.Printf("Failed to fetch worlds %v", err)
		return
	}

	for _, world := range worlds {
		for i := 0; i < 3; i++ {
			players, err := retry(func() ([]tibia.OnlinePlayers, error) {
				return ot.api.FetchOnlinePlayers(world.Name)
			}, 3)
			if err != nil {
				logger.Error.Printf("Failed to fetch players %v", err)
			} else {
				for _, player := range players {
					newLastSeen[player.Name] = time.Now()
				}
				break
			}
		}
	}

	logger.Info.Printf("Finished fetching online players in %v, onlineCount %v", time.Since(start), len(newLastSeen))
	ot.m.Lock()
	defer ot.m.Unlock()
	ot.lastSeen = newLastSeen
}

func (ot *OnlineScraper) GetLastSeen() map[string]time.Time {
	ot.m.RLock()
	defer ot.m.RUnlock()
	return ot.lastSeen
}
