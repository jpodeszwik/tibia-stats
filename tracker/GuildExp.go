package tracker

import (
	"sync"
	"tibia-stats/domain"
	"tibia-stats/dynamo"
	"tibia-stats/utils/formats"
	"time"
)

type GuildExp struct {
	lastExp      map[string]*domain.GuildExp
	repository   *dynamo.GuildExpRepository
	guildMembers map[string]string
	m            sync.RWMutex
}

func (ge *GuildExp) HandleGuildMembers(guildMembers map[string]string) {
	ge.m.Lock()
	ge.guildMembers = guildMembers
	ge.m.Unlock()
}

func (ge *GuildExp) HandleWorldExperience(exp domain.WorldExperience) {
	ge.m.RLock()
	playerGuild := ge.guildMembers
	ge.m.RUnlock()

	guildExperience := make(map[string]int64)
	for playerName, experience := range exp.Experience {
		guild, exists := playerGuild[playerName]
		if !exists {
			continue
		}

		guildExperience[guild] += experience
	}

	for guildName, experience := range guildExperience {
		ge.handleGuildExp(domain.GuildExp{
			GuildName: guildName,
			Exp:       experience,
			Date:      time.Now(),
		})
	}
}

func (ge *GuildExp) handleGuildExp(exp domain.GuildExp) {
	last, exists := ge.lastExp[exp.GuildName]

	if !exists || last.Exp != exp.Exp || last.Date.Format(formats.IsoDate) != exp.Date.Format(formats.IsoDate) {
		ge.repository.StoreGuildExp(exp)
		ge.lastExp[exp.GuildName] = &exp
	}
}

func NewGuildExp(repository *dynamo.GuildExpRepository) *GuildExp {
	return &GuildExp{
		lastExp:      make(map[string]*domain.GuildExp),
		repository:   repository,
		guildMembers: make(map[string]string),
	}
}
