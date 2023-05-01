package scraper

import (
	"sync"
	"tibia-stats/tibia"
	"tibia-stats/utils/logger"
	"time"
)

type Worlds struct {
	api    *tibia.ApiClient
	worlds []string
	m      sync.RWMutex
}

func (w *Worlds) Start() {
	err := w.fetchWorlds()
	if err != nil {
		logger.Error.Fatal(err)
	}

	go func() {
		ticker := time.NewTicker(4 * time.Hour)
		for range ticker.C {
			w.fetchWorlds()
		}
	}()
}

func (w *Worlds) fetchWorlds() error {
	worlds, err := retry(func() ([]tibia.OverviewWorld, error) {
		return w.api.FetchWorlds()
	}, 5)
	if err != nil {
		logger.Error.Printf("Failed to fetch worlds %v", err)
		return err
	}

	var newWorlds []string
	for _, world := range worlds {
		newWorlds = append(newWorlds, world.Name)
	}

	w.m.Lock()
	w.worlds = newWorlds
	w.m.Unlock()

	return nil
}

func (w *Worlds) getWorlds() []string {
	w.m.RLock()
	defer w.m.RUnlock()
	return w.worlds
}

func NewWorlds(api *tibia.ApiClient) *Worlds {
	return &Worlds{
		api: api,
		m:   sync.RWMutex{},
	}
}
