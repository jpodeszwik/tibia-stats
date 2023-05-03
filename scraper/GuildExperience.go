package scraper

import (
	"tibia-stats/domain"
	"tibia-stats/tibia"
	"tibia-stats/utils/logger"
	"time"
)

const guildExperienceRefreshInterval = 2 * time.Hour

type GuildExperience struct {
	api          *tibia.ApiClient
	handler      Handler[domain.GuildExp]
	worlds       *Worlds
	guildMembers *GuildMembers
}

func (ge *GuildExperience) Start() {
	go func() {
		ticker := time.NewTicker(guildExperienceRefreshInterval)
		ge.fetchGuildsExperience()
		for range ticker.C {
			ge.fetchGuildsExperience()
		}
	}()
}

func (ge *GuildExperience) fetchGuildsExperience() error {
	start := time.Now()
	worlds := ge.worlds.getWorlds()
	ret := make(map[string]int64)
	for _, world := range worlds {
		guildsExp, err := ge.fetchWorldGuildsExperience(world)
		if err != nil {
			logger.Error.Printf("Failed to fetch experience for world %v", world)
			continue
		}
		for guildName, exp := range guildsExp {
			ret[guildName] = exp
		}
	}

	logger.Info.Printf("Finished fetching %v worlds %v guilds experiences in %v", len(worlds), len(ret), time.Since(start))
	return nil
}

func (ge *GuildExperience) fetchWorldGuildsExperience(world string) (map[string]int64, error) {
	playerGuild := ge.guildMembers.getPlayerGuild()

	guildExp := make(map[string]int64)
	for _, profession := range tibia.AllProfessions {
		for page := 0; page < 20; page++ {
			highScorePage, err := retry(func() ([]tibia.HighscoreResponse, error) {
				return ge.api.FetchHighscore(world, profession, tibia.Exp, page)
			}, 3)

			if err != nil {
				return nil, err
			}
			for _, highScoreEntry := range highScorePage {
				guild, exists := playerGuild[highScoreEntry.Name]
				if !exists {
					continue
				}

				guildExp[guild] += highScoreEntry.Value
			}
		}
	}

	for guildName, exp := range guildExp {
		ge.handler.Handle(domain.GuildExp{
			GuildName: guildName,
			Exp:       exp,
			Date:      time.Now(),
		})
	}

	return guildExp, nil
}

func NewGuildExperience(client *tibia.ApiClient, worlds *Worlds, guildMembers *GuildMembers, handler Handler[domain.GuildExp]) *GuildExperience {
	return &GuildExperience{
		api:          client,
		worlds:       worlds,
		guildMembers: guildMembers,
		handler:      handler,
	}
}
