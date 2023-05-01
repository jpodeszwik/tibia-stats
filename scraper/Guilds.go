package scraper

import (
	"tibia-stats/tibia"
	"tibia-stats/utils/logger"
	"time"
)

type Guilds struct {
	api     *tibia.ApiClient
	handler Handler[[]string]
	worlds  *Worlds
}

func (g *Guilds) fetchGuilds() error {
	start := time.Now()

	var allGuilds []string
	worlds := g.worlds.getWorlds()
	for _, world := range worlds {
		guilds, err := retry(func() ([]tibia.OverviewGuild, error) {
			return g.api.FetchGuilds(world)
		}, 5)
		if err != nil {
			return err
		}

		for _, guild := range guilds {
			allGuilds = append(allGuilds, guild.Name)
		}
	}

	g.handler.Handle(allGuilds)

	logger.Info.Printf("Finished fetching %v worlds in %v", len(worlds), time.Since(start))
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

func NewGuilds(api *tibia.ApiClient, worlds *Worlds, handler Handler[[]string]) *Guilds {
	return &Guilds{
		api:     api,
		handler: handler,
		worlds:  worlds,
	}
}
