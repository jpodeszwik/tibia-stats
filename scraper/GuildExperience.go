package scraper

import (
	"fmt"
	"tibia-stats/domain"
	"tibia-stats/tibia"
	"tibia-stats/utils/logger"
	"time"
)

const guildExperienceRefreshInterval = 2 * time.Hour

type GuildExperience struct {
	api     *tibia.ApiClient
	handler Handler[domain.GuildExp]
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
	worlds, err := retry(func() ([]tibia.OverviewWorld, error) {
		return ge.api.FetchWorlds()
	}, 3)
	if err != nil {
		return err
	}
	ret := make(map[string]int64)
	for _, world := range worlds {
		guildsExp, err := ge.fetchWorldGuildsExperience(world.Name)
		if err != nil {
			logger.Error.Printf("Failed to fetch experience for world %v", world.Name)
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
	guilds, err := retry(func() ([]tibia.OverviewGuild, error) {
		return ge.api.FetchGuilds(world)
	}, 3)
	if err != nil {
		return nil, err
	}

	playerGuild := make(map[string]string)
	for _, overviewGuild := range guilds {
		guild, err := retry(func() (*tibia.GuildResponse, error) {
			return ge.api.FetchGuild(overviewGuild.Name)
		}, 3)
		if err != nil {
			return nil, err
		}
		for _, member := range guild.Members {
			playerGuild[member.Name] = overviewGuild.Name
		}
	}

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

func NewGuildExperience(client *tibia.ApiClient, handler Handler[domain.GuildExp]) *GuildExperience {
	return &GuildExperience{api: client, handler: handler}
}

func retry[T any](f func() (T, error), times int) (T, error) {
	var errs []error
	for i := 0; i < times; i++ {
		val, err := f()
		if err == nil {
			return val, err
		}
		errs = append(errs, err)
	}

	var zero T
	return zero, fmt.Errorf("failed after %v tries, %v", times, errs)
}
