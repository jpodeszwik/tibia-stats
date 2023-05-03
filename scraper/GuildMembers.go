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
	worlds      *Worlds
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
	worlds := gm.worlds.getWorlds()

	playerGuild := make(map[string]string)
	guilds := 0
	for _, world := range worlds {
		members, err := gm.fetchWorldGuildMembers(world)
		if err != nil {
			logger.Error.Printf("Failed to fetch guilds for world %v", world)
			return err
		}
		logger.Debug.Printf("Fetched guild members for %v", world)

		guilds += len(members)

		for guildName, guildMembers := range members {
			for _, member := range guildMembers {
				playerGuild[member] = guildName
			}
		}
	}

	gm.m.Lock()
	gm.playerGuild = playerGuild
	gm.m.Unlock()
	gm.handler(playerGuild)

	logger.Info.Printf("Finished fetching %v guilds for members, %v memberships found in %v", guilds, len(playerGuild), time.Since(start))
	return nil
}

func (gm *GuildMembers) fetchWorldGuildMembers(world string) (map[string][]string, error) {
	guilds := gm.guilds.getGuilds(world)

	guildMembers := make(map[string][]string)
	for _, guildName := range guilds {
		guild, err := retry(func() (*tibia.GuildResponse, error) {
			return gm.api.FetchGuild(guildName)
		}, 5)
		if err != nil {
			return nil, err
		}
		for _, member := range guild.Members {
			guildMembers[guildName] = append(guildMembers[guildName], member.Name)
		}
	}

	return guildMembers, nil
}

func (gm *GuildMembers) GetPlayerGuild() map[string]string {
	gm.m.RLock()
	defer gm.m.RUnlock()
	return gm.playerGuild
}

func NewGuildMembers(client *tibia.ApiClient, worlds *Worlds, guilds *Guilds, handler Handler[map[string]string]) *GuildMembers {
	return &GuildMembers{
		api:         client,
		worlds:      worlds,
		guilds:      guilds,
		playerGuild: make(map[string]string),
		handler:     handler,
	}
}
