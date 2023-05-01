package scraper

import (
	"log"
	"tibia-stats/tibia"
	"time"
)

type Guilds struct {
	api     *tibia.ApiClient
	handler Handler[[]string]
}

func (g *Guilds) fetchGuilds() error {
	start := time.Now()
	worlds, err := retry(func() ([]tibia.OverviewWorld, error) {
		return g.api.FetchWorlds()
	}, 3)
	if err != nil {
		return err
	}

	var allGuilds []string
	for _, world := range worlds {
		guilds, err := retry(func() ([]tibia.OverviewGuild, error) {
			return g.api.FetchGuilds(world.Name)
		}, 5)
		if err != nil {
			return err
		}

		for _, guild := range guilds {
			allGuilds = append(allGuilds, guild.Name)
		}
	}

	g.handler.Handle(allGuilds)

	log.Printf("Finished fetching %v worlds in %v", len(worlds), time.Since(start))
	return nil
}

func (g *Guilds) Start() {
	g.fetchGuilds()
	go func() {
		ticker := time.NewTicker(2 * time.Hour)
		for range ticker.C {
			g.fetchGuilds()
		}
	}()
}

func NewGuilds(api *tibia.ApiClient, handler Handler[[]string]) *Guilds {
	return &Guilds{
		api:     api,
		handler: handler,
	}
}
