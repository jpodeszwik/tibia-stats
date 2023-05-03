package tracker

import (
	"sync"
	"tibia-stats/domain"
	"tibia-stats/dynamo"
	"tibia-stats/utils/formats"
	"time"
)

type GuildExp struct {
	lastExp             map[string]*domain.GuildExp
	guildExpRepository  *dynamo.GuildExpRepository
	highScoreRepository *dynamo.HighScoreRepository
	guildMembers        map[string]string
	m                   sync.RWMutex
}

func (ge *GuildExp) HandleGuildMembers(guildMembers map[string]string) {
	ge.m.Lock()
	ge.guildMembers = guildMembers
	ge.m.Unlock()
}

func (ge *GuildExp) HandleWorldExperience(exp domain.WorldExperience) {
	now := time.Now()
	previousDayHighScore, err := ge.highScoreRepository.GetHighScore(exp.World, now.Add(-24*time.Hour))
	if err != nil {
		return
	}

	ge.m.RLock()
	playerGuild := ge.guildMembers
	ge.m.RUnlock()

	highScoreExp := make(map[string]int64)
	gainedExp := make(map[string]int64)
	for playerName, experience := range exp.Experience {
		guild, exists := playerGuild[playerName]
		if !exists {
			continue
		}

		highScoreExp[guild] += experience
		if previousDayHighScore != nil {
			previousDayExperience, exists := previousDayHighScore.Experience[playerName]
			if exists {
				gainedExp[guild] += experience - previousDayExperience
			}
		}
	}

	for guildName, experience := range highScoreExp {
		ge.handleGuildExp(domain.GuildExp{
			GuildName:    guildName,
			HighScoreExp: experience,
			GainedExp:    gainedExp[guildName],
			Date:         now,
		})
	}
}

func (ge *GuildExp) handleGuildExp(exp domain.GuildExp) {
	last, exists := ge.lastExp[exp.GuildName]

	if !exists || last.HighScoreExp != exp.HighScoreExp || last.Date.Format(formats.IsoDate) != exp.Date.Format(formats.IsoDate) {
		ge.guildExpRepository.StoreGuildExp(exp)
		ge.lastExp[exp.GuildName] = &exp
	}
}

func NewGuildExp(guildExpRepository *dynamo.GuildExpRepository, highScoreRepository *dynamo.HighScoreRepository) *GuildExp {
	return &GuildExp{
		lastExp:             make(map[string]*domain.GuildExp),
		guildExpRepository:  guildExpRepository,
		highScoreRepository: highScoreRepository,
		guildMembers:        make(map[string]string),
	}
}
