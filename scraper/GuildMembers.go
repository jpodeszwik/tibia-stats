package scraper

import (
	"sync"
	"tibia-stats/tibia"
	"tibia-stats/utils/logger"
	"time"
)

const guildMembersRefreshInterval = 2 * time.Hour

type GuildMembers struct {
	api         *tibia.ApiClient
	guilds      *Guilds
	playerGuild map[string]string
	handler     Handler[map[string]string]
	m           sync.RWMutex
}

func (gm *GuildMembers) Start() {
	ticker := time.NewTicker(guildMembersRefreshInterval)
	err := gm.fetchGuildMembers()
	if err != nil {
		logger.Error.Fatal(err)
	}

	go func() {
		for range ticker.C {
			gm.fetchGuildMembers()
		}
	}()
}

func (gm *GuildMembers) fetchGuildMembers() error {
	start := time.Now()
	guilds := gm.guilds.getGuilds()

	playerGuild := make(map[string]string)
	for _, guildName := range guilds {
		guild, err := retry(func() (*tibia.GuildResponse, error) {
			return gm.api.FetchGuild(guildName)
		}, 5)
		if err != nil {
			return err
		}

		for _, member := range guild.Members {
			playerGuild[member.Name] = guildName
		}
	}

	gm.m.Lock()
	gm.playerGuild = playerGuild
	gm.m.Unlock()
	gm.handler(playerGuild)

	logger.Info.Printf("Finished fetching %v guilds for members, %v memberships found in %v", len(guilds), len(playerGuild), time.Since(start))
	return nil
}

func (gm *GuildMembers) GetPlayerGuild() map[string]string {
	gm.m.RLock()
	defer gm.m.RUnlock()
	return gm.playerGuild
}

func NewGuildMembers(client *tibia.ApiClient, guilds *Guilds, handler Handler[map[string]string]) *GuildMembers {
	return &GuildMembers{
		api:         client,
		guilds:      guilds,
		playerGuild: make(map[string]string),
		handler:     handler,
	}
}
